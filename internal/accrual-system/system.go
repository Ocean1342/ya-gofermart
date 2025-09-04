package accrual_system

import (
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

//TODO: туть будет реальная имплементация сервиса

type System struct {
	client *http.Client
	url    string
	path   string
}

func New(url, path string) *System {
	return &System{
		client: initClient(),
		url:    url,
		path:   path,
	}
}

func (s *System) OrderProcess(orderID int) *http.Response {
	l := logrus.WithField("HANDLER", "AccrualSystem")
	l.Infof("start request orderID:%d", orderID)
	reqUrl := fmt.Sprintf("%s/%s/%d", s.url, s.path, orderID)
	resp, err := s.client.Get(reqUrl)
	if err != nil {
		l.Errorf("could not send request to accrual system. err:%s", err)
		return nil
	}
	return resp
}

func initClient() *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 5 * time.Second
	return retryClient.StandardClient()
}
