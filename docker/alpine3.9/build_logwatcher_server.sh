# from https://github.com/restic/restic/tree/master/docker
#!/bin/sh

set -e

echo "Build binary using golang docker image"
docker run --rm -ti \
  -v $(pwd):/go/src/github.com/xtforgame/log_mge_utils \
  -w /go/src/github.com/xtforgame/log_mge_utils \
  -e CGO_ENABLED=1 \
  -e GOOS=linux \
  -e GO111MODULE=on \
  golang:1.12-alpine3.9 go build -mod=vendor -o ./build/alpine3.9/logwatcher_server main/logwatcher_server.go

echo "Build docker image xtforgame/logwatcher_server:latest"
docker build --rm -t xtforgame/logwatcher_server:latest -f docker/alpine3.9/logwatcher_server/Dockerfile .

# docker run --rm -ti \
#   -p 8080:8080 \
#   -v $(pwd)/tmp:/usr/log_mge_utils \
#   -w /usr/log_mge_utils \
#   xtforgame/logwatcher_server:latest log_mge_utils ./forweb ./pgbackrest-backup ./output
