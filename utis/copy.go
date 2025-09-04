package utis

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func CopyHTTPRequest(req *http.Request) (*http.Request, error) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = req.Body.Close()
		if err != nil {
			logrus.Errorf("could not close request bodyBytes on copy. err: %v", err)
		}
	}()
	copyReq := req.Clone(req.Context())
	if bodyBytes != nil {
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		copyReq.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	return copyReq, nil
}

func CopyHTTPResponse(resp *http.Response) (*http.Response, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			logrus.Errorf("could not close response body on copy. err: %v", err)
		}
	}()
	newBody := io.NopCloser(bytes.NewReader(bodyBytes))
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	newHeader := make(http.Header)
	for k, v := range resp.Header {
		newHeader[k] = v
	}

	return &http.Response{
		Status:           resp.Status,
		StatusCode:       resp.StatusCode,
		Proto:            resp.Proto,
		ProtoMajor:       resp.ProtoMajor,
		ProtoMinor:       resp.ProtoMinor,
		Header:           newHeader,
		Body:             newBody,
		ContentLength:    resp.ContentLength,
		TransferEncoding: resp.TransferEncoding,
		Close:            resp.Close,
		Uncompressed:     resp.Uncompressed,
		Trailer:          resp.Trailer,
		Request:          resp.Request, // Это ссылка, при необходимости нужно клонировать и Request
		TLS:              resp.TLS,
	}, nil
}
