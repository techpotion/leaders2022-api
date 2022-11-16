#!/bin/bash

docker run --rm -v ${CI_PROJECT_DIR}:/app -w /app \
    -e CGO_ENABLED=0 \
    -e GOOS=linux \
    -e GOARCH=amd64 \
    golang:1.17-alpine3.15 \
    go test -v ./...
