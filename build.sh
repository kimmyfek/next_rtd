#!/bin/bash

docker run --rm -v "$PWD":/go/src/github.com/kimmyfek/next_rtd -w /go/src/github.com/kimmyfek/next_rtd golang:1.8 go get && export GOOS=darwin && go build -v -o nxt-darwin
