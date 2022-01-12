# build file
GOCMD=go
# Use -a flag to prevent code cache problems.
GOBUILD=$(GOCMD) build -ldflags -s -v

BIN_BINARY_NAME=das_parser_tool
parser_tool:
	$(GOBUILD) -o $(BIN_BINARY_NAME) cmd/main.go
	@echo "Build $(BIN_BINARY_NAME) successfully. You can run ./$(BIN_BINARY_NAME) now.If you can't see it soon,wait some seconds"

update:
	export GOPRIVATE="github.com/DeAccountSystems"
	go mod tidy
	go mod vendor
