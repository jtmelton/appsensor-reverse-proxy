package ids

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jtmelton/appsensor-reverse-proxy/connections"
	"github.com/parnurzeal/gorequest"
)

func AddEvent(category string, label string, r *http.Request) {

	// grab ip from request and use for username and ip address
	ip := connections.FindIp(r)

	event := &Event{
		User: User{
			Username: ip,
			IpAddress: IpAddress{
				Address: ip,
			},
		},
		DetectionPoint: DetectionPoint{
			Category: category,
			Label:    label,
		},
		Timestamp: time.Now().Format(time.RFC3339),
		DetectionSystem: DetectionSystem{
			DetectionSystemId: RestHeaderValue,
			IpAddress: IpAddress{
				Address: ClientIp,
			},
		},
	}

	json, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err)
		return
	}

	request := gorequest.New()
	// resp, body, errs :=
	request.Post(RestUrl+"/events").
		Set(RestHeaderName, RestHeaderValue).
		Send(string(json)).
		End()

}

func AddAttack(category string, label string, r *http.Request) {

	// grab ip from request and use for username and ip address
	ip := connections.FindIp(r)

	attack := &Attack{
		User: User{
			Username: ip,
			IpAddress: IpAddress{
				Address: ip,
			},
		},
		DetectionPoint: DetectionPoint{
			Category: category,
			Label:    label,
		},
		Timestamp: time.Now().Format(time.RFC3339),
		DetectionSystem: DetectionSystem{
			DetectionSystemId: RestHeaderValue,
			IpAddress: IpAddress{
				Address: ClientIp,
			},
		},
	}

	json, err := json.Marshal(attack)
	if err != nil {
		fmt.Println(err)
		return
	}

	request := gorequest.New()
	// resp, body, errs :=
	request.Post(RestUrl+"/attacks").
		Set(RestHeaderName, RestHeaderValue).
		Send(string(json)).
		End()

}
