package middleware

import (
	"net/http"

	"github.com/jtmelton/appsensor-reverse-proxy/Godeps/_workspace/src/github.com/golang/glog"
)

func Trend(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		glog.Info("Considering trend")

//		go func() {
//			// simulate long running task
//			glog.Info("about to sleep")
//			time.Sleep(2 * time.Second)
//			glog.Info("done sleeping")
//		}()
		
		next.ServeHTTP(w, r)

	})
}
