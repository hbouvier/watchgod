GOCC=go
# To Compile the linux version using docker simply invoke the makefile like this:
#
# make GOCC="docker run --rm -t -v ${GOPATH}:/go hbouvier/go-lang:1.5"

all: get-deps build

build: fmt test
	${GOCC} install github.com/hbouvier/watchgod

fmt:
	${GOCC} fmt github.com/hbouvier/watchgod
	${GOCC} fmt github.com/hbouvier/watchgod/libwatchgod

test:
	${GOCC} test -v -cpu 4 -count 1 ./libwatchgod/...

get-deps:
	${GOCC} get github.com/hashicorp/logutils
