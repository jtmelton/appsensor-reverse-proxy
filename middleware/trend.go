package middleware

import (
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/jtmelton/appsensor-reverse-proxy/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/jtmelton/appsensor-reverse-proxy/connections"
)

const (
	DELIMETER              = ":"
	DURATION_TIMESTAMP_KEY = "APPSENSOR_DURATION_TIMESTAMPS"
)

var (
	lastSerialized        = time.Now()
	serializationDuration = time.Duration(1) * time.Minute
	userResourceCounts    = make(map[string]int)
	mutex                 = &sync.Mutex{}
)

func Trend(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		go storeRequest(r)

		next.ServeHTTP(w, r)

	})
}

func storeRequest(r *http.Request) {
	ip := connections.FindIp(r)
	resource := r.URL.Path

	key := lastSerialized.String() + DELIMETER + ip + DELIMETER + resource

	mutex.Lock()
	// increment count of times this ip/resource pair has been seen
	userResourceCounts[key] = userResourceCounts[key] + 1
	mutex.Unlock()
	runtime.Gosched()

	// if we hit the serialization time ..
	if time.Now().After(lastSerialized.Add(serializationDuration)) {
		glog.Infof("Purging now - comparing %s + ms (%s) <--> %s",
			lastSerialized.String(), serializationDuration.String(),
			lastSerialized.Add(serializationDuration).String(), time.Now())

		serializedTimestampString := lastSerialized.String()

		// update last serialized for next batch
		lastSerialized = time.Now()

		flushToRedis(serializedTimestampString)
	}
}

func flushToRedis(serializedTimestampString string) {
	c := connections.RedisPool.Get()
	defer c.Close()

	// lock on main map
	mutex.Lock()

	// copy the main map to backup
	var tempData = make(map[string]int)
	for k, v := range userResourceCounts {
		tempData[k] = v
	}

	// "remake" main map
	userResourceCounts = make(map[string]int)

	// unlock main map
	mutex.Unlock()

	// send data from backup map to redis
	for k, v := range tempData {
		c.Send("SET", k, v)
	}

	c.Send("SADD", DURATION_TIMESTAMP_KEY, serializedTimestampString)

	c.Flush()

}
