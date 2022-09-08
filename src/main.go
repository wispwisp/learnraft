package main

import (
	"encoding/json"
	"flag"
	"io"
	"net/http"

	log "github.com/wispwisp/learnraft/logger"
	"github.com/wispwisp/learnraft/node"
	"github.com/wispwisp/learnraft/ping"
)

type Args struct {
	Port *string
	Init *bool
}

func registerArgs() (args Args) {
	args.Port = flag.String("port", "8090", "server port")
	args.Init = flag.Bool("init", false, "make initial actions")
	flag.Parse()
	return
}

type SomeJSONStruct struct {
	From   string `json:"from"`
	Amount int    `json:"amount"`
}

func main() {
	args := registerArgs()

	nodesFileName := "./nodes.json"
	var nodesInfo node.NodesInfo
	if err := nodesInfo.LoadFromFile(nodesFileName); err != nil {
		log.Info("Fail to load from", nodesFileName, "error:", err)
	}

	if false {
		statuses := ping.StartPing()
		ping.RecievePing(statuses)
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		log.Info("'/ping' HTTP handler")

		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Error("error parsing request body:", err)
			http.Error(w, "error parsing request", http.StatusBadRequest)
			return
		}

		var jsonRes map[string]interface{}
		err = json.Unmarshal(body, &jsonRes)
		if err != nil {
			log.Error("Unmarshal err:", err)
			return
		}

		log.Info(jsonRes)

		nodes := nodesInfo.Get()
		if encodeErr := json.NewEncoder(w).Encode(nodes); encodeErr != nil {
			log.Error("Encode to json failed, err: ", encodeErr)
			http.Error(w, "Encode to json failed", http.StatusBadRequest)
			// http.NotFound(w, req)
			return
		}
	})

	http.HandleFunc("/addnode", func(w http.ResponseWriter, req *http.Request) {
		log.Info("'/addnode' HTTP handler")

		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Error("error parsing request body:", err)
			http.Error(w, "error parsing request", http.StatusBadRequest)
			return
		}

		var nodeInfo node.NodeInfo
		err = json.Unmarshal(body, &nodeInfo)
		if err != nil {
			log.Error("error parsing node info:", err)
			http.Error(w, "error parsing node info", http.StatusBadRequest)
			return
		}

		nodesInfo.Add(&nodeInfo)
	})

	log.Info("Server started on", *args.Port, "port")
	http.ListenAndServe(":"+*args.Port, nil)
}
