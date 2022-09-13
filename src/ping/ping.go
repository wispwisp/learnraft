package ping

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/wispwisp/learnraft/mylogger"
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
		}
	}()

	return
}

// func RecievePingEx(statuses chan int) {
// 	go func() {
// 		for status := range statuses {
// 			log.Info("Status recieved:", status)
// 		}
// 	}()
// }

func SendVoteToOtherNodes(logger mylogger.Logger, nodesInfo *node.NodesInfo, v *node.Vote) []node.VoteResponse {
	ni := nodesInfo.Get()

	l := len(ni)

	var wg sync.WaitGroup

	votes := make([]node.VoteResponse, l)

	for i := 0; i < l; i++ {
		wg.Add(1)

		go func(iteration int, uri string) {
			defer wg.Done()

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
			// var res map[string]interface{}
			// json.NewDecoder(resp.Body).Decode(&res)
			// logger.Info("Vote response:", res)

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Error("error response request body:", err)
				return
			}

			var voteResp node.VoteResponse
			err = json.Unmarshal(body, &voteResp)
			if err != nil {
				logger.Error("error parsing vote response:", err)
				return
			}

			votes[iteration] = voteResp

		}(i, "http://"+ni[i].Uri+"/vote")
	}

	wg.Wait()
	return votes
}

func Elections(logger mylogger.Logger, nodeState *node.NodeState, nodesInfo *node.NodesInfo) {
	go func() {
		// r := rand.New(rand.NewSource(time.Now().UnixNano()))
		// v := 150 + r.Intn(150) // 150-300 ms randomized

		counter := 0
		ticker := time.NewTicker(1 * time.Second)
		for {
			<-ticker.C // TODO: only when not a candidate

			if nodeState.GetState() != node.CANDIDATE {
				<-ticker.C
				continue
			}

			// TODO: 150-300 ms randomized timeout

			counter++

			// Vote for self
			v := &node.Vote{NodeName: nodeState.GetUri()}
			votes := SendVoteToOtherNodes(logger, nodesInfo, v)
			logger.Info("votes recieved:", votes)

			// TODO: if consensus (N/2) accepts (all returns this node URI) - became a leader.
		}
	}()
}
