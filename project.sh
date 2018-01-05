#!/bin/bash
#set -ex
COMMAND=$1

#TARGET_OS=$2
#
#if [ -z "$TARGET_OS" ]; then
#	TARGET_OS="darwin"
#fi
#
#case "$TARGET_OS" in
#    darwin)
#        export GOARCH="amd64"
#        export GOOS="darwin"
#        ;;
#    linux)
#        export GOARCH="amd64"
#        export GOOS="linux"
#        ;;
#    *)
#        echo "Unknown Target OS: $TARGET_OS"
#        exit 1
#        ;;
#esac
#
#export CGO_ENABLED=1
#
#docker run --rm --name nxt-rtd-build \
#	-v $GOPATH:/go \
#	-w /go/src/github.com/kimmyfek/next_rtd \
#	-e GOOS=$GOOS \
#	-e GOARCH=$GOARCH \
#	#-e CGO_ENABLED=$CGO_ENABLED \
#	golang:1.8 \
#	go build -v -o nxt-$GOOS-$GOARCH
function rundb() {
	echo ""
	echo "-----------------------"
	echo "Running DB"
	echo "-----------------------"

	docker run --rm -it --name nxt-mysql \
		-e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
		-e MYSQL_DATABASE=rtd \
		-p 3306:3306 \
		mysql:5.7
}

case "$COMMAND" in
	build)
		build
		echo "Done"
		;;
	rundb)
		rundb
		echo "Done"
		;;
esac
