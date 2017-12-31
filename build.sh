#!/bin/bash

case "$OSTYPE" in
    darwin*)
        export HOST_OS="darwin"
        ;;
    linux*)
        export HOST_OS="linux"
        ;;
    *)
        echo "unknown OS Type: $OSTYPE"
        exit 1
        ;;
esac

docker run --rm -v "$PWD":/go/src/github.com/kimmyfek/next_rtd -w /go/src/github.com/kimmyfek/next_rtd golang:1.8 /bin/bash -c "go get && export GOOS=$HOST_OS && go build -v -o nxt-darwin"
