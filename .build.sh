#!/bin/sh

set -e
set -x

# compile the main binary
GOOS=linux GOARCH=amd64 CGO_ENABLED=0         go build -ldflags "-X main.build=${DRONE_BUILD_NUMBER}" -a -tags netgo -o release/linux/amd64/image-promote github.com/drone-plugins/image-promote/cmd/promote
