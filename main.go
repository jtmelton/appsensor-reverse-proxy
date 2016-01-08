// appsensor-reverse-proxy project main.go
package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/jtmelton/appsensor-reverse-proxy/blocks"
	"github.com/jtmelton/appsensor-reverse-proxy/connections"
	"github.com/jtmelton/appsensor-reverse-proxy/ids"
	"github.com/jtmelton/appsensor-reverse-proxy/middleware"

	"github.com/jtmelton/appsensor-reverse-proxy/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/jtmelton/appsensor-reverse-proxy/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/jtmelton/appsensor-reverse-proxy/Godeps/_workspace/src/github.com/justinas/alice"
)

type Prox struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func New(target string) *Prox {
	url, _ := url.Parse(target)

	return &Prox{target: url, proxy: httputil.NewSingleHostReverseProxy(url)}
}

func (p *Prox) handle(w http.ResponseWriter, r *http.Request) {
	glog.Info("saw a request for ", r.URL)

	p.proxy.ServeHTTP(w, r)
}

func doEvery(d time.Duration, f func()) {
	for {
		time.Sleep(d)
		f()
	}
}

func parseFlags() (*bool, *string, *string, *string, http.Handler) {

	const (
		defaultPort                     = ":8080"
		defaultPortUsage                = "default server port, ':80', ':8080'..."
		defaultTarget                   = "http://127.0.0.1:8090"
		defaultTargetUsage              = "default redirect url, 'http://127.0.0.1:8080'"
		defaultBooleanUsage             = "true or false"
		defaultRefreshUsage             = "number of seconds between block refreshes"
		defaultExpireUsage              = "number of seconds between expiration sweeps"
		defaultBlocksRefreshUrl         = "block refresh url, 'http://localhost:8090/api/v1.0/blocks'"
		defaultCertFile                 = "cert file (e.g. cert.pem)"
		defaultKeyFile                  = "key file (e.g. key.pem)"
		defaultResourceVerbsMappingFile = "resource verbs mapping file (e.g. resource-verbs-mapping.yml)"
		defaultResourcesFile            = "resources file (e.g. resources.yml)"
	)

	// flags
	port := flag.String("port", defaultPort, defaultPortUsage)
	url := flag.String("url", defaultTarget, defaultTargetUsage)

	// for running in TLS mode
	enableTls := flag.Bool("enable-tls", false, defaultBooleanUsage)
	certFile := flag.String("cert-file", "cert.pem", defaultCertFile)
	keyFile := flag.String("key-file", "key.pem", defaultKeyFile)

	// related to trend detection
	enableTrends := flag.Bool("enable-trend-tracking", false, defaultBooleanUsage)

	// related to blocking
	enableBlocking := flag.Bool("enable-blocking", false, defaultBooleanUsage)
	refreshSeconds := flag.Int("blocking-refresh-rate-seconds", 30, defaultRefreshUsage)
	expireSeconds := flag.Int("blocking-expire-rate-seconds", 30, defaultExpireUsage)
	blocksRefreshUrl := flag.String("blocking-blocks-refresh-url", "http://localhost", defaultBlocksRefreshUrl)

	// for redis, if used
	redisAddress := flag.String("redis-address", ":6379", "Address to the Redis server")
	maxRedisConnections := flag.Int("max-redis-connections", 10, "Max connections to Redis")

	// related to invalid http verbs
	// (instead of GET, POST, HEAD ... you receive DOG, SHOE, FOREST, etc.)
	enableInvalidVerbs := flag.Bool("enable-invalid-verb-checking", false, defaultBooleanUsage)

	// related to unexpected http verbs
	// requires config file
	// config file states this file reachable over GET or HEAD, and you get a POST
	enableUnexpectedVerbs := flag.Bool("enable-unexpected-verb-checking", false, defaultBooleanUsage)
	resourceVerbsMappingFile := flag.String("resource-verbs-mapping-file", "resource-verbs-mapping.yml", defaultResourceVerbsMappingFile)

	// related to unexpected resources
	// requires config file
	// config file states all resources expected in the app (think site-map)
	// this analysis should help detect scanners probing for non-existent resources
	enableUnexpectedResources := flag.Bool("enable-unexpected-resource-checking", false, defaultBooleanUsage)
	resourcesFile := flag.String("resources-file", "resources.yml", defaultResourcesFile)

	flag.Parse()

	glog.Info("------------------------------------------------")
	glog.Infof("Settings:")
	glog.Infof("\tServer port %s", *port)
	glog.Infof("\tProxy url: %s", *url)
	glog.Infof("\tEnable TLS: %t", *enableTls)
	glog.Infof("\t\tCert File: %s", *certFile)
	glog.Infof("\t\tKey File: %s", *keyFile)
	glog.Infof("\tEnable trend tracking: %t", *enableTrends)
	glog.Infof("\tEnable blocking: %t", *enableBlocking)
	glog.Infof("\t\tRefresh rate (seconds): %d", *refreshSeconds)
	glog.Infof("\t\tExpire rate (seconds): %d", *expireSeconds)
	glog.Infof("\t\tBlock refresh url: %s", *blocksRefreshUrl)
	glog.Infof("\tEnable invalid verbs: %t", *enableInvalidVerbs)
	glog.Infof("\tEnable unexpected verbs: %t", *enableUnexpectedVerbs)
	glog.Infof("\t\tResource verbs mapping file: %s", *resourceVerbsMappingFile)
	glog.Infof("\tEnable unexpected resources: %t", *enableUnexpectedResources)
	glog.Infof("\t\tResources file: %s", *resourcesFile)
	glog.Info("------------------------------------------------")

	// proxy
	proxy := New(*url)

	proxyHandler := http.HandlerFunc(proxy.handle)

	// chain default handlers (clear context, recovery, and log execution time)
	chain := alice.New(
		context.ClearHandler,
		middleware.Recovery,
		middleware.LogExecutionTime)

	// if blocking is enabled, add it to alice chain
	if *enableBlocking {
		chain = chain.Append(middleware.Block)
		go doEvery((time.Duration(*refreshSeconds) * time.Second), func() { blocks.RefreshBlocks(blocksRefreshUrl) })
		go doEvery((time.Duration(*expireSeconds) * time.Second), blocks.ExpireBlocks)
	}

	// if trending is enabled, add it to alice chain
	if *enableTrends {
		chain = chain.Append(middleware.Trend)
	}

	// if invalid verb checking is enabled, add it to alice chain
	if *enableInvalidVerbs {
		chain = chain.Append(middleware.InvalidVerbs)
	}

	// if unexpected verb checking is enabled, add it to alice chain
	if *enableUnexpectedVerbs {
		middleware.PopulateExpectedVerbs(resourceVerbsMappingFile)
		chain = chain.Append(middleware.UnexpectedVerbs)
	}

	// if unexpected resource checking is enabled, add it to alice chain
	if *enableUnexpectedResources {
		middleware.PopulateExpectedResources(resourcesFile)
		chain = chain.Append(middleware.UnexpectedResources)
	}

	// if any of these settings are enabled, connect to redis
	if *enableTrends {
		connections.ConnectRedis(redisAddress, maxRedisConnections)
	}

	// if any of these settings are enabled, setup rest client to appsensor backend
	if *enableInvalidVerbs || *enableTrends ||
		*enableUnexpectedVerbs || *enableUnexpectedResources {
		ids.InitializeAppSensorRestClient()
	}

	// always use reverse proxy as chain target (final)
	chainHandler := chain.Then(proxyHandler)

	return enableTls, port, certFile, keyFile, chainHandler
}

func main() {

	enableTls, port, certFile, keyFile, chainHandler := parseFlags()

	if connections.RedisPool != nil {
		defer connections.RedisPool.Close()
	}

	if *enableTls {

		// serve over https
		// if you want to test this locally, generate a cert and key file using
		// go by running a command similar to :
		// "go run /usr/local/go/src/crypto/tls/generate_cert.go --host localhost"
		// and change the location to wherever go is installed on your system
		err := http.ListenAndServeTLS(*port, *certFile, *keyFile, chainHandler)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// serve over http (non-secure)
		glog.Infof("WARNING: Running over plaintext http (insecure). " +
			"To enable https, use '-enable-tls'")
		http.ListenAndServe(*port, chainHandler)
	}

}
