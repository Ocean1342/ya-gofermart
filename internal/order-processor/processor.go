package orderprocessor

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	accrualsystem "gofermart/internal/accrual-system"
	"gofermart/internal/storage"
	"io"
	"time"
)

type Process interface {
	Process(ctx context.Context)
}

type OrderQueueItem struct {
	OrderID    int
	UserID     int
	RetryTimes int
}

type OrderProcessor struct {
	AccrualService accrualsystem.AccrualService
	Storage        storage.Storage
	Queue          chan OrderQueueItem
	tick           time.Duration
}

func New(service accrualsystem.AccrualService, storage storage.Storage, tick time.Duration) *OrderProcessor {
	ch := make(chan OrderQueueItem, 1_000_000)
	return &OrderProcessor{
		AccrualService: service,
		Storage:        storage,
		Queue:          ch,
		tick:           tick,
	}
}
func (op *OrderProcessor) Process(ctx context.Context) {
	logger := logrus.WithFields(map[string]interface{}{
		"HANDLER": "OrderProcessor",
	})
	queueTicker := time.NewTicker(op.tick)
	storageInterrogateTicker := time.NewTicker(op.tick + 1)
	for {
		select {
		case <-ctx.Done():
			logger.Info("Order processor stop by ctx")
			return
		case <-storageInterrogateTicker.C:
			go op.interrogateProcess(ctx)
		case <-queueTicker.C:
			item, ok := <-op.Queue
			if !ok {
				logger.Info("queue channel closed")
				return
			}
			if item.OrderID != 0 && item.UserID != 0 {
				go op.process(ctx, *logger, item)
			}
		case item := <-op.Queue:
			go op.process(ctx, *logger, item)
		}
	}
}

func (op *OrderProcessor) process(ctx context.Context, logger logrus.Entry, item OrderQueueItem) {
	resp := op.AccrualService.OrderProcess(item.OrderID)
	if resp == nil {
		logger.Errorf("empty response from accrualService")
		item.RetryTimes++
		op.Queue <- item
		return
	}
	switch resp.StatusCode {
	case 200:
		var acResp AccrualResponse
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Errorf("could not read resp body. err: %v", err)
			item.RetryTimes++
			op.Queue <- item
			return
		}
		err = json.Unmarshal(bytes, &acResp)
		if err != nil {
			logger.Errorf("could not unmarshal resp body. err: %v", err)
			item.RetryTimes++
			op.Queue <- item
			return
		}
		orderID, err := acResp.GetOrderID()
		if err != nil {
			logger.Errorf("could not convert Order id")
			item.RetryTimes++
			op.Queue <- item
			return
		}
		logger.Infof("user id %d get raw accrual: %f counted accrual: %d", item.UserID, acResp.Accrual, acResp.GetAccrual())
		err = op.Storage.UpdateOrderAndBalance(ctx, orderID, item.UserID, acResp.GetAccrual(), acResp.GetStatus())
		if err != nil {
			logger.Errorf("UpdateOrderAndBalance err: %v", err)
			item.RetryTimes++
			op.Queue <- item
			return
		}
	case 204:
		bytes, err := io.ReadAll(resp.Body)
		if err != nil && !errors.Is(err, io.EOF) {
			logger.Errorf("recived 204. orderID %d read body err:%s", item.OrderID, err)
		}
		logger.Errorf("recived 204. orderID %d response body: %s", item.OrderID, string(bytes))
		err = op.Storage.UpdateOrderAndBalance(ctx, item.OrderID, item.UserID, 0, "INVALID")
		if err != nil {
			logger.Errorf("could not change status for order id:%d", item.OrderID)
		}
	case 429:
		logger.Errorf("too many requests.sleeps for 60 sec")
		time.Sleep(60 * time.Second)
		op.Queue <- item
	case 500:
		logger.Errorf("Accrual 500")
		op.Queue <- item
	default:
		logger.Errorf("undefined response from Accrual system item: %d", item)
		item.RetryTimes++
		op.Queue <- item
	}
}
