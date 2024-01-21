.PHONY: run
run:
	go run ./app

.PHONY: test
test:
	go test -v ./app/...	

.PHONY: build
build:
	@echo 'Building ...'
	go build -v -ldflags="-s -w" -o=./bin/vt-manager ./app
	GOOS=linux GOARCH=amd64 go build -v -ldflags="-s -w" -o=./bin/linux_amd64/vt-manager-amd64 ./app