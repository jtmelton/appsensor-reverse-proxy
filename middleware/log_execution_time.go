package middleware

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/jtmelton/appsensor-reverse-proxy/Godeps/_workspace/src/github.com/golang/glog"
)

// from https://justinas.org/writing-http-middleware-in-go/
func LogExecutionTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		rec := httptest.NewRecorder()
		// passing a ResponseRecorder instead of the original RW
		next.ServeHTTP(rec, r)

		// we copy the original headers first
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}

		w.WriteHeader(rec.Code)

		// then write out the original body
		w.Write(rec.Body.Bytes())

		glog.Infof("%v %s on %s %s in %v", rec.Code, http.StatusText(rec.Code), r.Method, r.URL.Path, time.Since(start))
		
	})
}
