.PHONY: run
run:
	go run ./app

.PHONY: test
test:
	go test ./app/...	

.PHONY: build
build:
	@echo 'Building ...'
	go build -ldflags="-s -w" -o=./bin/vt-manager ./app
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o=./bin/linux_amd64/vt-manager-amd64 ./app