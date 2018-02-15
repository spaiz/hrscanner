#!/usr/bin/env sh
# Used by instance from the Dockerfile.build image to create statically compiled binary
set -ev
REPOSITORY_PATH=${GOPATH}/src/github.com/spaiz/hrscanner
CGO_ENABLED=0 GOOS=linux go build -v -o ${REPOSITORY_PATH}/artifacts/hrscanner -a -tags netgo -ldflags '-w' .