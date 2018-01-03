#!/bin/bash
set -ex

TARGET_OS=$2

if [ -z "$TARGET_OS" ]; then
	TARGET_OS="darwin"
fi

case "$TARGET_OS" in
    darwin)
        export GOARCH="amd64"
        export GOOS="darwin"
        ;;
    linux)
        export GOARCH="amd64"
        export GOOS="linux"
        ;;
    *)
        echo "Unknown Target OS: $TARGET_OS"
        exit 1
        ;;
esac

export CGO_ENABLED=1

docker run --rm --name nxt-rtd-build \
	-v $GOPATH:/go \
	-w /go/src/github.com/kimmyfek/next_rtd \
	-e GOOS=$GOOS \
	-e GOARCH=$GOARCH \
	#-e CGO_ENABLED=$CGO_ENABLED \
	golang:1.8 \
	go build -v -o nxt-$GOOS-$GOARCH
