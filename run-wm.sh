#!/bin/bash

#CGO_CFLAGS_ALLOW=-Xpreprocessor \


ENV="dev" \
SITE="testapi.com" \
go run cmd/rainier/main.go
