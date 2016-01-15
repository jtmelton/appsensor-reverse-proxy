package blocks

import (
	"encoding/json"
	"time"

	"github.com/golang/glog"
)

func ExpireBlocks() {
	glog.Info("Evaluating blocks for expiration.")

	var removableBlocks []string

	for _, element := range StoredBlocks.Flatten() {

		var block Block

		//glog.Info(element)
		if err := json.Unmarshal([]byte(element.(string)), &block); err != nil {
			panic(err)
		}

		t, err := unixToTime(block.Endtime)
		if err != nil {
			glog.Fatal(err)
		}

		if time.Now().After(t) {
			glog.Info("Expiring block: ", block)
			removableBlocks = append(removableBlocks, element.(string))
		} else {
			glog.Infof("Not expiring block:")
			glog.Infof("\tnow: ", time.Now())
			glog.Infof("\tt: ", t)
			glog.Infof("\tnow unix: ", time.Now().Unix())
			glog.Infof("\tt unix: ", t.Unix())
		}

	}

	if len(removableBlocks) > 0 {
		for _, removableBlock := range removableBlocks {
			StoredBlocks.Remove(removableBlock)
		}
	}

}
