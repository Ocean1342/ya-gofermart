package orderprocessor

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

func (op *OrderProcessor) interrogateProcess(ctx context.Context) {
	unprocessedOrders, err := op.Storage.GetUnprocessedOrders(ctx, 50)
	if err != nil {
		logrus.Errorf("interrogateProcess could not get orders. err: %v", err)
		time.Sleep(10 * time.Second)
		return
	}
	for _, order := range unprocessedOrders {
		op.Queue <- OrderQueueItem{
			OrderID:    order.OrderID,
			UserID:     order.UserID,
			RetryTimes: 0,
		}
	}
}
