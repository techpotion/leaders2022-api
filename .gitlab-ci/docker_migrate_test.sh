#!/bin/bash

docker run --rm \
    -v ${CI_PROJECT_DIR}/migrations:/migrations \
    --network host migrate/migrate:4 \
	-path=/migrations/ \
	-database ${DB_POSTGRESQL_URI_TEST}?sslmode=disable \
	up
