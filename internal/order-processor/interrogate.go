package orderprocessor

import (
	"context"
	"github.com/sirupsen/logrus"
)

func (op *OrderProcessor) interrogateProcess(ctx context.Context) {
	unprocessedOrders, err := op.Storage.GetUnprocessedOrders(ctx, 50)
	if err != nil {
		logrus.Errorf("could not get orders. err: %v", err)
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
