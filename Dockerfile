# Dockerfile
FROM golang:1.9.4-alpine3.6
MAINTAINER Alexander R. <spaizadv@gmail.com>

ARG REPOSITORY_PATH=${GOPATH}/src/github.com/spaiz/hrscanner
COPY ./data ${REPOSITORY_PATH}/data
COPY ./artifacts/hrscanner ${GOPATH}/bin
WORKDIR ${REPOSITORY_PATH}
CMD ["hrscanner"]