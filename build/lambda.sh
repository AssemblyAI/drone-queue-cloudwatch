#!/usr/bin/env bash

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main ./

zip drone-queue-cloudwatch.zip main