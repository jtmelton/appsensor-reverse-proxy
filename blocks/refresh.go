package blocks

import (
	"encoding/json"
	"time"

	"github.com/golang/glog"
	"github.com/parnurzeal/gorequest"
)

func RefreshBlocks(blockRefreshUrl *string) {
	//"http://localhost:8090/api/v1.0/blocks"
	_, body, _ := gorequest.New().Get(*blockRefreshUrl).End()

	var blocks Blocks

	if err := json.Unmarshal([]byte(body), &blocks); err != nil {
		panic(err)
	}

	for _, element := range blocks {

		// temporary hack for dealing w/ java block store
		// it keeps timestamp in microseconds for some reason
		element.Endtime = element.Endtime / 1000

		t, _ := unixToTime(element.Endtime)

		// only add if this is still a valid time block
		if time.Now().Before(t) {
			//back to json so it's hashable
			jsonStr, _ := json.Marshal(element)
			StoredBlocks.Add(string(jsonStr))
		}

	}

	//	glog.Info(StoredBlocks.Len())
	glog.Infof("Retrieved %d blocks, total stored: %d", len(blocks), StoredBlocks.Len())
	//	glog.Info("STORED blocks: ", StoredBlocks.Flatten())
	//	glog.Info("Printing response: ", resp)
	//	glog.Info("Printing body: ", body)
}
