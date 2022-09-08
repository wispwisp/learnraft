#!/bin/bash

# 8090

cd ../src/ && go build

./learnraft -port 8090 -log ~/logs/node1.log &
NODE1PID=$!

./learnraft -port 8091 -log ~/logs/node2.log &
NODE2PID=$!

./learnraft -port 8092 -log ~/logs/node3.log &
NODE3PID=$!

cd ../scripts

echo "$NODE1PID $NODE2PID $NODE3PID"

sleep 3

curl -X POST 127.0.0.1:8090/addnode \
     -H 'Content-Type: application/json' \
     -d '{"uri":"127.0.0.1:8091"}'

curl -X POST 127.0.0.1:8090/addnode \
     -H 'Content-Type: application/json' \
     -d '{"uri":"127.0.0.1:8092"}'

# 8091

curl -X POST 127.0.0.1:8091/addnode \
     -H 'Content-Type: application/json' \
     -d '{"uri":"127.0.0.1:8090"}'

curl -X POST 127.0.0.1:8091/addnode \
     -H 'Content-Type: application/json' \
     -d '{"uri":"127.0.0.1:8092"}'

# 8092

curl -X POST 127.0.0.1:8092/addnode \
     -H 'Content-Type: application/json' \
     -d '{"uri":"127.0.0.1:8090"}'

curl -X POST 127.0.0.1:8092/addnode \
     -H 'Content-Type: application/json' \
     -d '{"uri":"127.0.0.1:8091"}'

ps uax | grep learnraft
