#!/bin/bash

# Build go microservice
cd ./micro
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/main .

docker-compose up