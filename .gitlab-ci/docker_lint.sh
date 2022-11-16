#!/bin/bash

docker run --rm \
  -v ${CI_PROJECT_DIR}:/app \
  -w /app \
  golangci/golangci-lint:v1.44-alpine \
  golangci-lint run -v --timeout 5m
