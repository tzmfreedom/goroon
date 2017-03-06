.PHONY: test install uninstall

test: build
	@go test -v ./...

install: build
	@go install ./goroon

uninstall:

clean:
	@go clean

build: clean
	@goimports -w .
	@gofmt -w .
	@go build -o cmd goroon
