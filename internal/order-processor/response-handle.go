package orderprocessor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

func (op *OrderProcessor) StatusOkHandle(ctx context.Context, logger logrus.Entry, resp *http.Response, item OrderQueueItem) error {
	var acResp AccrualResponse
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read resp body. err: %v", err)
	}
	err = json.Unmarshal(bytes, &acResp)
	if err != nil {
		return fmt.Errorf("could not unmarshal resp body. err: %v", err)
	}
	orderID, err := acResp.GetOrderID()
	if err != nil {
		return fmt.Errorf("could not convert Order id")
	}
	logger.Infof("user id %d get raw accrual: %f counted accrual: %d", item.UserID, acResp.Accrual, acResp.GetAccrual())
	ctxProcess, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	err = op.Storage.UpdateOrderAndBalance(ctxProcess, orderID, item.UserID, acResp.GetAccrual(), acResp.GetStatus())
	if err != nil {
		return fmt.Errorf("UpdateOrderAndBalance err: %v", err)
	}
	return nil
}

func (op *OrderProcessor) StatusNoContentHandle(ctx context.Context, logger logrus.Entry, resp *http.Response, item OrderQueueItem) error {
	bytes, err := io.ReadAll(resp.Body)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("recived 204. orderID %d read body err:%s", item.OrderID, err)
	}
	logger.Errorf("recived 204. orderID %d response body: %s", item.OrderID, string(bytes))
	ctxProcess, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	err = op.Storage.UpdateOrderAndBalance(ctxProcess, item.OrderID, item.UserID, 0, "INVALID")
	if err != nil {
		return fmt.Errorf("could not change status for order id:%d err:%v", item.OrderID, err)
	}
	return nil
}
