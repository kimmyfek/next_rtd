#!/bin/bash
#set -ex

COMMAND=$1
TARGET_OS=$2

function build() {

	if [ -z "$TARGET_OS" ]; then
		TARGET_OS="darwin"
	fi

	case "$TARGET_OS" in
		darwin)
			export GOARCH="amd64"
			export GOOS="darwin"
			export CGO_ENABLED=0
			;;
		linux)
			export GOARCH="amd64"
			export GOOS="linux"
			export CGO_ENABLED=1
			;;
		*)
			echo "Unknown Target OS: $TARGET_OS"
			exit 1
			;;
	esac


	docker run --rm --name nxt-rtd-build \
		-v $(pwd):/go/src/github.com/kimmyfek/next_rtd \
		-w /go/src/github.com/kimmyfek/next_rtd \
		-e GOOS=$GOOS \
		-e GOARCH=$GOARCH \
		-e CGO_ENABLED=$CGO_ENABLED \
		golang:1.8 \
		/bin/bash -c "go get && go build -v -o nxt-$GOOS-$GOARCH"
}

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

function debug() {
	# Should do a docker delete nxt-mysql-debug and not care if it fails
	docker run --rm -d --name nxt-mysql-debug \
		-e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
		-e MYSQL_DATABASE=rtd \
		-p 3306:3306 \
		mysql:5.7
	echo "Taking a nap while mysql turns on"
	sleep 10
	./nxt-darwin-amd64 --level=debug
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
	debug)
		build
		debug
esac
