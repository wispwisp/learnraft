#!/bin/bash

curl -X POST 127.0.0.1:8090/add \
     -H 'Content-Type: application/json' \
     -d '{"key":"me","val":"{\"inner_json\":42, \"other_field\":\"some text\"}"}'
