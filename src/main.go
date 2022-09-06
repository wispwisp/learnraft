package main

import (
	"encoding/json"
	"flag"
	"io"
	"net/http"

	log "github.com/wispwisp/learnraft/logger"
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

	fileName := "./conf.json"
	example := SomeJSONStruct{From: fileName}

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

		if encodeErr := json.NewEncoder(w).Encode(example); encodeErr != nil {
			log.Error("Encode to json failed, err: ", encodeErr)
			http.Error(w, "Encode to json failed", http.StatusBadRequest)
			// http.NotFound(w, req)
			return
		}
	})

	log.Info("Server started on", *args.Port, "port")
	http.ListenAndServe(":"+*args.Port, nil)
}
