GOCC=go
VERSION=v1.0.3
# To Compile the linux version using docker simply invoke the makefile like this:
#
# make GOCC="docker run --rm -t -v ${GOPATH}:/go hbouvier/go-lang:1.5"

all: get-deps build

clean:
	rm -f coverage.out

build: fmt test
	${GOCC} install github.com/hbouvier/watchgod

fmt:
	${GOCC} fmt github.com/hbouvier/watchgod
	${GOCC} fmt github.com/hbouvier/watchgod/libwatchgod

test:
	${GOCC} test -v -cpu 4 -count 1 -coverprofile=coverage.out github.com/hbouvier/watchgod/libwatchgod/...
	${GOCC} tool cover -html=coverage.out

get-deps:
	${GOCC} get github.com/hashicorp/logutils

release_darwin_amd64:
	rm -rf "release/bin/darwin_amd64" "${GOPATH}/bin/watchgod"
	mkdir -p release/bin/darwin_amd64
	go install github.com/hbouvier/watchgod
	cp ${GOPATH}/bin/watchgod release/bin/darwin_amd64/

release_linux_adm64:
	rm -rf "release/bin/linux_adm64" "${GOPATH}/bin/watchgod"
	mkdir -p release/bin/linux_adm64
	docker run --rm -t -v ${GOPATH}:/go hbouvier/go-lang:1.5 install github.com/hbouvier/watchgod
	cp ${GOPATH}/bin/watchgod release/bin/linux_adm64/

archives: release_darwin_amd64 release_linux_adm64
	cd release && zip -r watchgod_${VERSION}.zip bin/
	cd release && COPYFILE_DISABLE=1 tar cvzf watchgod_${VERSION}.tgz  bin/