package middleware

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/jtmelton/appsensor-reverse-proxy/ids"
)

// this functionality covers the a completely invalid HTTP verb is used (ie. GOTO)
// whitelist is: [HEAD, GET, POST, PUT, DELETE, TRACE, OPTIONS, CONNECT]
// https://www.owasp.org/index.php/AppSensor_DetectionPoints#RE2:_Attempt_to_Invoke_Unsupported_HTTP_Method

// TRACE Verb might be valid , but it is not a very good idea to have it enabled.
func InvalidVerbs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		go evaluateInvalidVerbs(r)

		next.ServeHTTP(w, r)

	})
}

func evaluateInvalidVerbs(r *http.Request) {

	if !isValidVerb(r.Method) {
		glog.Info("Invalid HTTP verb seen, creating event.")
		go ids.AddEvent("Request", "RE2", r)
	}

}

func isValidVerb(verb string) bool {
	switch verb {
	case
		"OPTIONS",
		"GET",
		"HEAD",
		"POST",
		"PUT",
		"DELETE",
		"TRACE",
		"CONNECT":
		return true
	}
	return false
}
