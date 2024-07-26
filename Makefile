GO := go
BIN := ddoh

$(BIN): $(wildcard *.go) $(wildcard */*.go)
	CGO_ENABLED=0 $(GO) build -ldflags='-s -w -buildid=' .

clean:
	rm -rf $(BIN)
	go clean
