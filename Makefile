GOCC=go
VERSION=v$(shell cat VERSION.txt)
INSTALL_FLAGS=-ldflags "-X main.version=${VERSION}"
# To Compile the linux version using docker simply invoke the makefile like this:
#
# make GOCC="docker run --rm -t -v ${GOPATH}:/go hbouvier/go-lang:1.5"

USERNAME=hbouvier
PROJECTNAME=watchgod

# all: get-deps fmt darwin linux arm windows build coverage
all: get-deps fmt darwin linux arm build coverage

clean:
	rm -rf coverage.out \
	       ${GOPATH}/pkg/{linux_amd64,darwin_amd64,linux_arm}/github.com/${USERNAME}/${PROJECTNAME} \
	       ${GOPATH}/bin/{linux_amd64,darwin_amd64,linux_arm}/${PROJECTNAME} \
		   watchgod.upx \
		   watchgod \
	       release

build: gen fmt test
	${GOCC} install ${INSTALL_FLAGS} github.com/${USERNAME}/${PROJECTNAME}

gen:
	${GOCC} generate github.com/${USERNAME}/${PROJECTNAME}/...
	
fmt:
	${GOCC} fmt github.com/${USERNAME}/${PROJECTNAME}

test:
	# ${GOCC} test -v -cpu 4 -count 1 -coverprofile=coverage.out github.com/${USERNAME}/${PROJECTNAME}
	${GOCC} test -v -cpu 4 -count 1 -coverprofile=coverage.out github.com/${USERNAME}/${PROJECTNAME}/process/...

coverage:
	${GOCC} tool cover -html=coverage.out

get-deps:
	${GOCC} install github.com/kulshekhar/fungen
	${GOCC} get github.com/hashicorp/logutils \
	            github.com/kulshekhar/fungen \
				github.com/opencontainers/runc/libcontainer/user \
				github.com/opencontainers/runc/libcontainer/system

linux:
	# GOOS=linux GOARCH=amd64 CGO_ENABLED=0 ${GOCC} install github.com/kulshekhar/fungen
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 ${GOCC} install ${INSTALL_FLAGS} github.com/${USERNAME}/${PROJECTNAME}
	@if [[ $(shell uname | tr '[:upper:]' '[:lower:]') == $@ ]] ; then mkdir -p ${GOPATH}/bin/$@_amd64 && mv ${GOPATH}/bin/${PROJECTNAME} ${GOPATH}/bin/$@_amd64/ ; fi
	# Pass from 9.4M to 5.9M when using `-ldflags="-s -w"` and then to 1.7M when also using `upx -f --brute`
	# go build -ldflags="-s -w" watchgod.go
	# upx -f --brute watchgod
	docker run -ti -v ${GOPATH}:/go golang /bin/sh -c "apt-get -y update && apt-get install -y upx && cd src/github.com/${USERNAME}/${PROJECTNAME} && go build -ldflags='-s -w -X main.version=${VERSION}' ${PROJECTNAME}.go && upx -f --brute ${PROJECTNAME}"
	mv ${PROJECTNAME} ${GOPATH}/bin/$@_amd64/

darwin:
	# GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 ${GOCC} install github.com/kulshekhar/fungen
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 ${GOCC} install ${INSTALL_FLAGS} github.com/${USERNAME}/${PROJECTNAME}
	@if [[ $(shell uname | tr '[:upper:]' '[:lower:]') == $@ ]] ; then mkdir -p ${GOPATH}/bin/$@_amd64 && mv ${GOPATH}/bin/${PROJECTNAME} ${GOPATH}/bin/$@_amd64/ ; fi

arm:
	# GOOS=linux GOARCH=arm CGO_ENABLED=0 ${GOCC} install github.com/kulshekhar/fungen
	GOOS=linux GOARCH=arm CGO_ENABLED=0 ${GOCC} install ${INSTALL_FLAGS} github.com/${USERNAME}/${PROJECTNAME}
	@if [[ $(shell uname | tr '[:upper:]' '[:lower:]') == $@ ]] ; then mkdir -p ${GOPATH}/bin/$@_amd64 && mv ${GOPATH}/bin/${PROJECTNAME} ${GOPATH}/bin/$@_amd64/ ; fi

windows:
	# GOOS=windows GOARCH=amd64 CGO_ENABLED=0 ${GOCC} install github.com/kulshekhar/fungen
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 ${GOCC} install ${INSTALL_FLAGS} github.com/${USERNAME}/${PROJECTNAME}
	@if [[ $(shell uname | tr '[:upper:]' '[:lower:]') == $@ ]] ; then mkdir -p ${GOPATH}/bin/$@_amd64 && mv ${GOPATH}/bin/${PROJECTNAME}.exe ${GOPATH}/bin/$@_amd64/ ; fi

# release: linux darwin arm windows
# 	@mkdir -p release/bin/{linux_amd64,darwin_amd64,linux_arm,windows_amd64}
# 	for i in linux_amd64 darwin_amd64 linux_arm; do cp ${GOPATH}/bin/$${i}/${PROJECTNAME} release/bin/$${i}/ ; done
# 	cp ${GOPATH}/bin/windows_amd64/${PROJECTNAME}.exe release/bin/windows_amd64/
# 	COPYFILE_DISABLE=1 tar cvzf release/${PROJECTNAME}.v`cat VERSION`.tgz release/bin
# 	zip -r release/${PROJECTNAME}.v`cat VERSION`.zip release/bin

release: linux darwin arm
	@mkdir -p release/bin/{linux_amd64,darwin_amd64,linux_arm}
	for i in linux_amd64 darwin_amd64 linux_arm; do cp ${GOPATH}/bin/$${i}/${PROJECTNAME} release/bin/$${i}/; done
	cd release && COPYFILE_DISABLE=1 tar cvzf ${PROJECTNAME}.${VERSION}.tgz bin
	cd release && zip -r ${PROJECTNAME}.${VERSION}.zip bin

info:
	@echo "version ${VERSION}"

tag-release:
	# git tag ${VERSION}
	# git push origin v${VERSION}
	./upload-release.sh
