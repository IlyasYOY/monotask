.PHONY: test
test: bulid-bin-monotask
	MONOTASK_BINARY=./bin/monotask go test -race -v -shuffle=on -failfast -fullpath ./... 

.PHONY: bulid-bin-monotask:
bulid-bin-monotask:
	go build -o ./bin/monotask ./cmd/monotask

.PHONY: clean
clean: 
	rm -rf ./bin
