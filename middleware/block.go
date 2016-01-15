package middleware

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/jtmelton/appsensor-reverse-proxy/blocks"
	"github.com/jtmelton/appsensor-reverse-proxy/connections"
)

func Block(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := connections.FindIp(r)
		resource := r.URL.Path

		shouldBlock := false

		for _, element := range blocks.StoredBlocks.Flatten() {

			var block blocks.Block

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
