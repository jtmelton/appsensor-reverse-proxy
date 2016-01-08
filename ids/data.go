package ids

import (
	"log"
	"os"
	"strings"

	"github.com/jtmelton/appsensor-reverse-proxy/Godeps/_workspace/src/github.com/golang/glog"
)

type Event struct {
	User            User            `json:"user"`
	DetectionPoint  DetectionPoint  `json:"detectionPoint"`
	Timestamp       string          `json:"timestamp"`
	DetectionSystem DetectionSystem `json:"detectionSystem"`
}

type User struct {
	Username  string    `json:"username"`
	IpAddress IpAddress `json:"ipAddress"`
}

type IpAddress struct {
	Address string `json:"address"`
}

type DetectionPoint struct {
	Category string `json:"category"`
	Label    string `json:"label"`
}

type DetectionSystem struct {
	DetectionSystemId string    `json:"detectionSystemId"`
	IpAddress         IpAddress `json:"ipAddress"`
}

type Attack Event

var (
	RestUrl         string
	urlExists       bool
	RestHeaderName  string
	nameExists      bool
	RestHeaderValue string
	valueExists     bool
	ClientIp        string
	ipExists        bool
)

// read environment variables for data for http client construction
func InitializeAppSensorRestClient() {

	RestUrl, urlExists = os.LookupEnv("APPSENSOR_REST_ENGINE_URL")
	RestHeaderName, nameExists = os.LookupEnv("APPSENSOR_CLIENT_APPLICATION_ID_HEADER_NAME")
	RestHeaderValue, valueExists = os.LookupEnv("APPSENSOR_CLIENT_APPLICATION_ID_HEADER_VALUE")
	ClientIp, ipExists = os.LookupEnv("APPSENSOR_CLIENT_APPLICATION_IP_ADDRESS")

	missingFields := make([]string, 0)

	if !urlExists {
		missingFields = append(missingFields, "APPSENSOR_REST_ENGINE_URL")
	}

	if !nameExists {
		glog.Info("The APPSENSOR_CLIENT_APPLICATION_ID_HEADER_NAME env var not set, using default value")
		RestHeaderName = "X-Appsensor-Client-Application-Name"
	}

	if !valueExists {
		missingFields = append(missingFields, "APPSENSOR_CLIENT_APPLICATION_ID_HEADER_VALUE")
	}

	if !ipExists {
		missingFields = append(missingFields, "APPSENSOR_CLIENT_APPLICATION_ID_HEADER_VALUE")
	}

	if len(missingFields) > 0 {
		log.Fatal("The following environment variable(s) must be populated: [ " + strings.Join(missingFields, " , ") + " ]")
	}

	glog.Info("Rest client information configured properly.")

}
