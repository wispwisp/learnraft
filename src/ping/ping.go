package ping

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/wispwisp/learnraft/logger"
	"github.com/wispwisp/learnraft/node"
)

func StartPingEx() (statuses chan int) {
	statuses = make(chan int)

	go func() {
		defer close(statuses)

		counter := 0
		ticker := time.NewTicker(1 * time.Second)
		for {
			<-ticker.C
			statuses <- counter
			counter++

			log.Info("Status send:", counter)
		}
	}()

	return
}

func RecievePingEx(statuses chan int) {
	go func() {
		for status := range statuses {
			log.Info("Status recieved:", status)

		}
	}()
}

func SendVoteToOtherNodes(logger *log.FileLogger, nodesInfo *node.NodesInfo, v *node.Vote) {
	ni := nodesInfo.Get()
	for _, nodeInfo := range ni {
		go func(uri string) {
			logger.Info("Send vote to", uri)

			jsonData, err := json.Marshal(v)
			if err != nil {
				logger.Error(err)
				return
			}

			resp, err := http.Post(uri, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				logger.Error(err)
				return
			}

			// Log result
			var res map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&res)
			logger.Info("Vote response:", res)

		}("http://" + nodeInfo.Uri + "/vote")
	}
}

func Elections(logger *log.FileLogger, nodeState *node.NodeState, nodesInfo *node.NodesInfo) {
	go func() {
		counter := 0
		ticker := time.NewTicker(5 * time.Second) // TODO: 150-300 ms randomized
		for {
			<-ticker.C
			counter++

			v := &node.Vote{NodeName: nodeState.GetUri()}
			SendVoteToOtherNodes(logger, nodesInfo, v)
		}
	}()
}
