bin/monotask:
	go build -o ./bin/monotask ./cmd/monotask

.PHONY: test
test: clean
	@${MAKE} bin/monotask

	MONOTASK_BINARY=$(shell pwd)/bin/monotask go test -fullpath -shuffle=on ./... 

.PHONY: test-coverage
test-coverage: clean
	@mkdir -p "$(shell pwd)/unit.coverprofile"
	@mkdir -p "$(shell pwd)/binary.coverprofile"
	@GOFLAGS="-cover" ${MAKE} bin/monotask

	@# GOCOVERDIR is for the binary built with -cover, see https://go.dev/doc/build-cover#running
	MONOTASK_BINARY="$(shell pwd)/bin/monotask" BINARY_GOCOVERDIR="$(shell pwd)/binary.coverprofile" \
		go test ./... -test.shuffle=on -test.fullpath -cover -test.gocoverdir="$(shell pwd)/unit.coverprofile" 

	@mkdir -p "$(shell pwd)/merged.coverprofile"
	@go tool covdata merge -i=unit.coverprofile,binary.coverprofile -o merged.coverprofile
	@go tool covdata textfmt -i=merged.coverprofile -o coverage.out

	@# Replace the module path with the absolute repository path
	@# So I can jump to files.
	@go tool cover -func=./coverage.out \
		| sort \
		| sed "s|github.com/IlyasYOY/monotask|$(shell pwd)|g"
	@go tool cover -html=./coverage.out -o coverage.html

.PHONY: clean
clean: 
	rm -fr binary.coverprofile/ merged.coverprofile/ unit.coverprofile/ coverage.out coverage.html bin/monotask
