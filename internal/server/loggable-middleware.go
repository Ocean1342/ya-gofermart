package server

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gofermart/pkg/common"
	"gofermart/utis"
	"io"
	"net/http"
)

type logResponseWriter struct {
	innerWriter http.ResponseWriter
	reqID       string
}

func (l *logResponseWriter) Header() http.Header {
	return l.innerWriter.Header()
}
func (l *logResponseWriter) Write(b []byte) (int, error) {
	l.Header().Set(common.XRequestId, l.reqID)
	logrus.Debugf("RESPONSE: REQUEST-ID: `%s`, body: `%s`, headers: `%s`", l.reqID, string(b), l.Header())
	return l.innerWriter.Write(b)
}

func (l *logResponseWriter) WriteHeader(statusCode int) {
	l.innerWriter.WriteHeader(statusCode)
}

func loggable(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		reqCopy, err := utis.CopyHTTPRequest(r)
		if err != nil {
			logrus.Errorf("could not copy request")
		}
		reqID, ok := r.Context().Value(middleware.RequestIDKey).(string)
		if !ok || reqID == "" {
			reqIDFromHeader := r.Header.Get(common.XRequestId)
			if reqIDFromHeader == "" {
				reqID = uuid.New().String()
				logrus.Errorf("empty request id. set manualy")
			} else {
				reqID = reqIDFromHeader
				logrus.Info("setting request ID from header")
			}
		}
		bodyBytes, err := io.ReadAll(reqCopy.Body)
		logrus.Debugf("REQUEST: REQUEST-ID:`%s`.url:`%s`. method:`%s` body: `%s`. headers: `%s`",
			reqID, reqCopy.RequestURI, reqCopy.Method, string(bodyBytes), reqCopy.Header)
		logRW := &logResponseWriter{innerWriter: w, reqID: reqID}
		next.ServeHTTP(logRW, r)
	}
	return http.HandlerFunc(fn)
}
