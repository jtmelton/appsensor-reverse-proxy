package middleware

import (
	"net/http"

	"github.com/golang/glog"
)

func Recovery(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				glog.Errorf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
