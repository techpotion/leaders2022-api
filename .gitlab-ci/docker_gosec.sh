#!/bin/bash

docker run --rm \
    -v ${CI_PROJECT_DIR}:/app \
    -w /app \
    securego/gosec:2.10.0 \
    /app/...
