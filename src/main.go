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
	Port    *string
	Init    *bool
	LogFile *string
}

func registerArgs() (args Args) {
	args.Port = flag.String("port", "8090", "server port")
	args.Init = flag.Bool("init", false, "make initial actions")
	args.LogFile = flag.String("log", "./node.log", "file for logs")
	flag.Parse()
	return
}

func main() {
	args := registerArgs()

	logger, err := log.NewFileLogger(*args.LogFile)
	if err != nil {
		panic("Fail to initialize logger")
	}

	defer logger.Close()

	nodesFileName := "./nodes.json"
	var nodesInfo node.NodesInfo
	if err := nodesInfo.LoadFromFile(nodesFileName); err != nil {
		logger.Info("Fail to load from", nodesFileName, "error:", err)
	}

	if false {
		statuses := ping.StartPing()
		ping.RecievePing(statuses)
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		logger.Info("'/ping' HTTP handler")

		body, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Error("error parsing request body:", err)
			http.Error(w, "error parsing request", http.StatusBadRequest)
			return
		}

		var jsonRes map[string]interface{}
		err = json.Unmarshal(body, &jsonRes)
		if err != nil {
			logger.Error("Unmarshal err:", err)
			return
		}

		logger.Info(jsonRes)

		nodes := nodesInfo.Get()
		if encodeErr := json.NewEncoder(w).Encode(nodes); encodeErr != nil {
			logger.Error("Encode to json failed, err: ", encodeErr)
			http.Error(w, "Encode to json failed", http.StatusBadRequest)
			// http.NotFound(w, req)
			return
		}
	})

	http.HandleFunc("/addnode", func(w http.ResponseWriter, req *http.Request) {
		logger.Info("'/addnode' HTTP handler")

		body, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Error("error parsing request body:", err)
			http.Error(w, "error parsing request", http.StatusBadRequest)
			return
		}

		var nodeInfo node.NodeInfo
		err = json.Unmarshal(body, &nodeInfo)
		if err != nil {
			logger.Error("error parsing node info:", err)
			http.Error(w, "error parsing node info", http.StatusBadRequest)
			return
		}

		nodesInfo.Add(&nodeInfo)
	})

	http.HandleFunc("/nodes", func(w http.ResponseWriter, req *http.Request) {
		if encodeErr := json.NewEncoder(w).Encode(nodesInfo.Get()); encodeErr != nil {
			logger.Error("Encode to json failed, err: ", encodeErr)
			http.Error(w, "Encode to json failed", http.StatusBadRequest)
			// http.NotFound(w, req)
			return
		}
	})

	logger.Info("Server started on", *args.Port, "port")
	err = http.ListenAndServe(":"+*args.Port, nil)
	if err != nil {
		logger.Error("Server start error", err)
	}
}
