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

// this functionality covers the situation where a user has requested a resource
// that is not specified (ie. allowed) in the resources yml file
// see https://www.owasp.org/index.php/AppSensor_DetectionPoints#ACE3:_Force_Browsing_Attempt
func UnexpectedResources(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		go evaluateUnexpectedResources(r)

		next.ServeHTTP(w, r)

	})
}

func evaluateUnexpectedResources(r *http.Request) {

	if staticPaths.Exists(r.URL.Path) {
		// found static route, we're good
		return
	}

	for _, re := range regexPaths {
		if re.MatchString(r.URL.Path) {
			// found match with regex, we're good
			return
		}
	}

	//create event - didn't find a match
	glog.Info("Did not find matching resource - creating event.")
	go ids.AddEvent("Access Control", "ACE1", r)

}

var staticPaths = set.New()
var regexPaths = make([]*regexp.Regexp, 0)

//var expectedResources = set.New()

type ResourcesConfig struct {
	Resources []string
}

func PopulateExpectedResources(resourcesYamlFile *string) {
	var config ResourcesConfig

	source, err := ioutil.ReadFile(*resourcesYamlFile)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}

	for _, path := range config.Resources {

		if strings.HasPrefix(path, "REGEX|") {
			trimmed := path[6:len(path)]

			// regex
			glog.Info("Mapped regex route = ", trimmed)

			r, _ := regexp.Compile(trimmed)

			regexPaths = append(regexPaths, r)
		} else {
			// static route
			glog.Info("Mapped static route = ", path)

			staticPaths.Add(path)
			//= append(staticPaths, path)
		}

	}

	glog.Info("All regex routes: ", regexPaths)
	glog.Info("All static routes: ", staticPaths)

}
