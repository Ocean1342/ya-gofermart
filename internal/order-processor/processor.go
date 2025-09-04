package order_processor

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	accrual_system "gofermart/internal/accrual-system"
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
	AccrualService accrual_system.AccrualService
	Storage        storage.Storage
	Queue          chan OrderQueueItem
}

func New(service accrual_system.AccrualService, storage storage.Storage) *OrderProcessor {
	ch := make(chan OrderQueueItem, 1_000_000)
	return &OrderProcessor{
		AccrualService: service,
		Storage:        storage,
		Queue:          ch,
	}
}
func (op *OrderProcessor) Process(ctx context.Context) {
	logger := logrus.WithFields(map[string]interface{}{
		"HANDLER": "OrderProcessor",
	})
	for {
		select {
		case <-ctx.Done():
			logger.Info("order processor stop by ctx")
			return
		case item := <-op.Queue:
			op.process(ctx, *logger, item)
		}
	}
}

func (op *OrderProcessor) process(ctx context.Context, logger logrus.Entry, item OrderQueueItem) {
	resp := op.AccrualService.OrderProcess(item.OrderID)
	switch resp.StatusCode {
	case 200:
		//маршалить структуру
		var acResp AccrualResponse
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Errorf("could not read resp body. err: %s", err)
			item.RetryTimes++
			op.Queue <- item
			return
		}
		err = json.Unmarshal(bytes, &acResp)
		if err != nil {
			logger.Errorf("could not unmarshal resp body. err: %s", err)
			item.RetryTimes++
			op.Queue <- item
			return
		}
		orderID, err := acResp.GetOrderID()
		if err != nil {
			logger.Errorf("could not convert order id")
			item.RetryTimes++
			op.Queue <- item
			return
		}
		err = op.Storage.UpdateOrderAndBalance(ctx, orderID, item.UserID, acResp.GetAccrual(), acResp.GetStatus())
		if err != nil {
			item.RetryTimes++
			op.Queue <- item
			return
		}
	case 204:
		bytes, err := io.ReadAll(resp.Body)
		if err != nil && !errors.Is(err, io.EOF) {
			logger.Errorf("recived 204. item %s read body err:%s", item, err)
		}
		logger.Errorf("recived 204. item %s response body: %s", item, string(bytes))
		item.RetryTimes++
		op.Queue <- item
	case 429:
		logger.Errorf("too many requests.sleeps for 60 sec")
		time.Sleep(60 * time.Second)
		op.Queue <- item
	case 500:
		logger.Errorf("accrual 500")
		op.Queue <- item
	default:
		logger.Errorf("undefined response from accrual system item: %d", item)
	}
}
