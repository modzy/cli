.PHONY: all
all: fmt lint vet test

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	staticcheck ./...

.PHONY: test
test:
	go test -covermode=count -coverprofile=combined.coverprofile ./...

.PHONY: vet
vet:
	go vet ./modzy/...

.PHONY: build
build:
	go build ./modzy/...

.PHONY: clean
clean:
	find . -name \*.coverprofile -delete
