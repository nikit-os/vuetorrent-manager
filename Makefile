VERSION = $(shell git describe --tags --always --dirty)
TARGETOS = darwin
TARGETARCH = amd64

.PHONY: run
run:
	go run ./app

.PHONY: test
test:
	go test -v ./app/...	

.PHONY: build
build:
	@echo 'Building ...'
	GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -ldflags="-s -w -X n1kit0s/vt-manager/app/cmd.version=${VERSION}" -o=./bin/vt-manager-${TARGETOS}-${TARGETARCH}_${VERSION} ./app