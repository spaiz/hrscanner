#!/usr/bin/env bash
# Use this helper to build static compiled binary using docker
set -ev
REPOSITORY=/go/src/github.com/spaiz/hrscanner
docker build -t hrscanner-build -f Dockerfile.build .
docker run -it --rm -w /go/src/github.com/spaiz/hrscanner -v ${PWD}:${REPOSITORY} hrscanner-build:latest ${REPOSITORY}/docker/build.sh