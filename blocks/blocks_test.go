package blocks

import (
	"testing"
	"time"
)

func TestApplies(t *testing.T) {

	//current time in unix - ms since epoch
	now := time.Now()
	ip := "1.2.3.4"
	resource := "/a/b/c"

	oneHour, _ := time.ParseDuration("1h")

	afterNow := Block{
		Endtime:   time.Now().Add(oneHour).Unix(),
		Ipaddress: "1.2.3.4",
		Resource:  "/a/b/c",
	}

	//fmt.Println("tester: ", tester)
	//fmt.Println("afterNow: ", afterNow)

	if !afterNow.Applies(ip, resource, now) {
		t.Error("expected", "does not apply", "got", "applies")
	}
}
