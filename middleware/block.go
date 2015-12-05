package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/jtmelton/appsensor-reverse-proxy/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/jtmelton/appsensor-reverse-proxy/blocks"
)

func Block(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := findIp(r)
		resource := r.URL.Path

		shouldBlock := false

		for _, element := range blocks.StoredBlocks.Flatten() {

			var block blocks.Block

			//glog.Info(element)
			if err := json.Unmarshal([]byte(element.(string)), &block); err != nil {
				panic(err)
			}

			if block.Applies(ip, resource, time.Now()) {
				shouldBlock = true
				glog.Info("Found a matching block - denying request: ", block)
				break
			}

		}

		if shouldBlock {

			// deny access
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Access Denied"))

		} else {
			next.ServeHTTP(w, r)
		}

	})
}

func findIp(r *http.Request) string {

	if r.Header.Get("X-Forwarded-For") != "" {
		ip := r.Header.Get("X-Forwarded-For")
		return ip
	}

	remoteip, _, _ := net.SplitHostPort(r.RemoteAddr)

	return remoteip
}
