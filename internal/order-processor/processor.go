package orderprocessor

import (
	"context"
	"github.com/sirupsen/logrus"
	accrualsystem "gofermart/internal/accrual-system"
	"gofermart/internal/storage"
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
				logger.Infof("ticker recieved item orderID:%d, userID: %d retrytimes: %d", item.OrderID, item.UserID, item.RetryTimes)
				go op.process(ctx, *logger, item)
			}
		case item := <-op.Queue:
			logger.Infof("recieved item orderID:%d, userID: %d retrytimes: %d", item.OrderID, item.UserID, item.RetryTimes)
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
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		err := op.StatusOkHandle(ctx, logger, resp, item)
		if err != nil {
			logger.Errorf("could not handle 200 response. err: %v", err)
			item.RetryTimes++
			op.Queue <- item
		}
	case 204:
		err := op.StatusNoContentHandle(ctx, logger, resp, item)
		if err != nil {
			logger.Errorf("could not handle 204 response. err: %v", err)
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
