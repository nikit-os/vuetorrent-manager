VERSION = $(shell git describe --tags --always --dirty)

.PHONY: run
run:
	go run ./app

.PHONY: test
test:
	go test -v ./app/...	

.PHONY: build
build:
	@echo 'Building ...'
	GOOS=darwin GOARCH=arm64 go build -v -ldflags="-s -w -X n1kit0s/vt-manager/app/cmd.version=${VERSION}" -o=./bin/vt-manager-arm64-${VERSION} ./app
	GOOS=linux GOARCH=amd64 go build -v -ldflags="-s -w -X n1kit0s/vt-manager/app/cmd.version=${VERSION}" -o=./bin/vt-manager-amd64-${VERSION} ./app