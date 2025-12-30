.PHONY: test
test: bin/monotask
	MONOTASK_BINARY=$(abspath $<) go test -race -v -shuffle=on -failfast -fullpath ./... 

bin/monotask:
	go build -o $@ ./cmd/monotask

