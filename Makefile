.PHONY: test
test: build-bin-monotask
	MONOTASK_BINARY=$(shell pwd)/bin/monotask go test -race -v -shuffle=on -failfast -fullpath ./... 

.PHONY: build-bin-monotask
build-bin-monotask:
	go build -o ./bin/monotask ./cmd/monotask

.PHONY: clean
clean: 
	rm -rf ./bin
