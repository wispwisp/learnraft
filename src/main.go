package main

import (
	"encoding/json"
	"flag"
	"io"
	"net/http"

	"github.com/wispwisp/learnraft/mylogger"
	"github.com/wispwisp/learnraft/node"
	"github.com/wispwisp/learnraft/ping"
	"github.com/wispwisp/learnraft/storage"
)

type Args struct {
	Addr    *string
	Port    *string
	Init    *bool
	LogFile *string
}

func registerArgs() (args Args) {
	args.Addr = flag.String("addr", "127.0.0.1", "server addr")
	args.Port = flag.String("port", "8090", "server port")
	args.Init = flag.Bool("init", false, "make initial actions")
	args.LogFile = flag.String("log", "./node.log", "file for logs")
	flag.Parse()
	return
}

func main() {
	args := registerArgs()

	var logger mylogger.Logger
	if true { // TODO: logger type from args
		filelogger, err := mylogger.NewFileLogger(*args.LogFile)
		if err != nil {
			panic("Fail to initialize logger")
		}
		defer filelogger.Close()
		logger = filelogger
	}

	nodesFileName := "./nodes.json"
	var nodesInfo node.NodesInfo
	if err := nodesInfo.LoadFromFile(nodesFileName); err != nil {
		logger.Info("Fail to load from", nodesFileName, "error:", err)
	}

	var stor storage.Storage
	if true { // TODO: storage type from args
		stor = storage.NewFileStorage("./node_"+*args.Port, logger)
	}

	nodeState := node.NewNodeState(*args.Addr, *args.Port)

	ping.Elections(logger, nodeState, &nodesInfo)

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

	http.HandleFunc("/vote", func(w http.ResponseWriter, req *http.Request) {
		logger.Info("'/vote' HTTP handler")

		body, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Error("error parsing request body:", err)
			http.Error(w, "error parsing request", http.StatusBadRequest)
			return
		}

		var vote node.Vote
		err = json.Unmarshal(body, &vote)
		if err != nil {
			logger.Error("error parsing vote:", err)
			http.Error(w, "error parsing vote", http.StatusBadRequest)
			return
		}

		logger.Info("Vote recieved:", vote)

		// Accept new leader
		voteResponse := node.VoteResponse{NewLeader: vote.NodeName}
		if encodeErr := json.NewEncoder(w).Encode(voteResponse); encodeErr != nil {
			logger.Error("Encode response to json failed, err: ", encodeErr)
			http.Error(w, "Encode response to json failed", http.StatusBadRequest)
			return
		}
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, req *http.Request) {
		logger.Info("'/add' HTTP handler")

		// TODO (Log Replication):
		// If not Leader - proxy to leader
		// Add to storage (Uncommited)
		// Send to others
		// Wait for response from majority (N/2)
		// Commit on leader. Response to client here.
		// Send to other that change commited
		// Other set that commited.

		body, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Error("error parsing request body:", err)
			http.Error(w, "error parsing request", http.StatusBadRequest)
			return
		}

		var msg storage.Message
		err = json.Unmarshal(body, &msg)
		if err != nil {
			logger.Error("error parsing message:", err)
			http.Error(w, "error parsing message", http.StatusBadRequest)
			return
		}

		logger.Info("Message:", msg)

		// TODO: json as a value - embedded commands
		var jsonRes map[string]interface{}
		err = json.Unmarshal([]byte(msg.Value), &jsonRes)
		if err != nil {
			logger.Error("Unmarshal err:", err)
			return
		} else {
			logger.Info("Messages json:", jsonRes)
		}

		success := stor.Add(msg.Key, msg.Value)
		if !success {
			logger.Info("add message failed")
		}
	})

	// TODO: recieve leader ping. Drop timeout for candidate to become a leader
	// If no leader ping recieved, follower become candidate and vote for himself
	// --- USE TERM, if node not voted in this term, its vote for candidate

	logger.Info("Server started on", *args.Port, "port")
	err := http.ListenAndServe(":"+*args.Port, nil)
	if err != nil {
		logger.Error("Server start error:", err)
	}
}
