package blocks

import (
	//"fmt"
	"time"

	"github.com/jtmelton/appsensor-reverse-proxy/Godeps/_workspace/src/github.com/Workiva/go-datastructures/set"
)

type Block struct {
	Ipaddress string `json:"ipAddress"`
	Resource  string `json:"resource"`
	Endtime   int64  `json:"endTime"`
}

func (b Block) Applies(ip string, res string, t time.Time) bool {

	end, _ := unixToTime(b.Endtime)
	if end.Before(t) {
		return false
	}
	
	// both have values
	if ip != "" && res != "" {
		//fmt.Println("both")
		return ip == b.Ipaddress && res == b.Resource
	}
	
	// only ip has value
	if ip != "" {
		//fmt.Println("ip only")
		return ip == b.Ipaddress
	}

	// only res has value
	if res != "" {
		//fmt.Println("res only")
		return res == b.Resource
	}

	return false
}

type Blocks []Block

var StoredBlocks = set.New()

func unixToTime(s int64) (time.Time, error) {
	//return time.Unix(0, ms*int64(time.Millisecond)), nil
	return time.Unix(s, 0), nil
}
