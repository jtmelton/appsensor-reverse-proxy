package connections

import (
	"net"
	"net/http"

	"github.com/jtmelton/appsensor-reverse-proxy/Godeps/_workspace/src/github.com/garyburd/redigo/redis"
	"github.com/jtmelton/appsensor-reverse-proxy/Godeps/_workspace/src/github.com/golang/glog"
)

var (
	RedisPool *redis.Pool
)

func ConnectRedis(redisAddress *string, maxRedisConnections *int) {

	glog.Infof("Attempting redis connection at %s with %d connections",
		*redisAddress, *maxRedisConnections)

	RedisPool = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("dtcp", *redisAddress)

		if err != nil {
			glog.Fatal("redis connection could not be made")
			return nil, err
		}

		return c, err
	}, *maxRedisConnections)
}

func FindIp(r *http.Request) string {

	if r.Header.Get("X-Forwarded-For") != "" {
		ip := r.Header.Get("X-Forwarded-For")
		return ip
	}

	remoteip, _, _ := net.SplitHostPort(r.RemoteAddr)

	return remoteip
}
