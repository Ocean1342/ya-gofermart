package server

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

func decompressable(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Content-Encoding")
		if h == "gzip" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				r.Body = io.NopCloser(bytes.NewBuffer(body))
				next.ServeHTTP(w, r)
				return
			}
			dec, err := gzip.NewReader(bytes.NewBuffer(body))
			if dec == nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer dec.Close()

			r.Body = io.NopCloser(bytes.NewBuffer(body))
			//defer r.Body.Close()
			next.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}
