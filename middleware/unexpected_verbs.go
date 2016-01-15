package middleware

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/Workiva/go-datastructures/set"
	"github.com/golang/glog"
	"github.com/jtmelton/appsensor-reverse-proxy/ids"
)

// this functionality covers the situation where a valid HTTP verb is used in an
// unexpected place (ie. GET when expecting POST)
// see https://www.owasp.org/index.php/AppSensor_DetectionPoints#RE1:_Unexpected_HTTP_Command
func UnexpectedVerbs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		go evaluateUnexpectedVerbs(r)

		next.ServeHTTP(w, r)

	})
}

func evaluateUnexpectedVerbs(r *http.Request) {

	if config.EnableGlobalPreflightRequests && r.Method == "OPTIONS" {
		// using OPTIONS always allowed in preflight mode
		return
	}

	if staticPaths.Exists(r.URL.Path) {
		for _, verb := range config.Resources[r.URL.Path] {
			if verb == r.Method {
				// found matching verb .. bail
				return
			}
		}

		glog.Info("Invalid verb was found for static path, creating event.")
		go ids.AddEvent("Request", "RE1", r)
		return
	}

	for _, re := range regexPaths {
		if re.MatchString(r.URL.Path) {

			for _, verb := range config.Resources["REGEX|"+r.URL.Path] {
				if verb == r.Method {
					// found matching verb .. bail
					return
				}
			}

			glog.Info("Invalid verb was found for regex path, creating event.")
			go ids.AddEvent("Request", "RE1", r)
			return
		}
	}

	if config.EvaluateUnlistedResources {
		glog.Info("No resource listing matched, creating event.")
		go ids.AddEvent("Request", "RE1", r)
	}
}

var config VerbsConfig

type VerbsConfig struct {
	EnableGlobalPreflightRequests bool
	EvaluateUnlistedResources     bool
	Resources                     map[string][]string
}

var staticResourcePaths = set.New()
var regexResourcePaths = make([]*regexp.Regexp, 0)

func PopulateExpectedVerbs(verbsYamlFile *string) {
	//var config VerbsConfig

	source, err := ioutil.ReadFile(*verbsYamlFile)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}

	for key, _ := range config.Resources {

		if strings.HasPrefix(key, "REGEX|") {
			trimmed := key[6:len(key)]

			// regex
			glog.Info("Mapped regex route = ", trimmed)

			r, _ := regexp.Compile(trimmed)

			regexResourcePaths = append(regexResourcePaths, r)
		} else {
			// static route
			glog.Info("Mapped static route = ", key)

			staticResourcePaths.Add(key)
		}

	}

	glog.Info("All regex resources: ", regexResourcePaths)
	glog.Info("All static resources: ", staticResourcePaths)

	glog.Info("parsed out config ", config)
}
