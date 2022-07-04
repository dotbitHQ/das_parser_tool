# Use -a flag to prevent code cache problems.
GO_BUILD=go build -ldflags -s -v
BINARY_NAME=tx_parser

update:
	export GOPRIVATE="github.com/dotbitHQ"
	go mod tidy

tx_parser_linux:
	export GOOS=linux
	export GOARCH=amd64
	$(GO_BUILD) -o $(BINARY_NAME) main.go
	mkdir -p bin/linux
	mv $(BINARY_NAME) bin/linux/
	@echo "build $(BINARY_NAME) successfully."

tx_parser_mac:
	export GOOS=darwin
	export GOARCH=amd64
	$(GO_BUILD) -o $(BINARY_NAME) main.go
	mkdir -p bin/mac
	mv $(BINARY_NAME) bin/mac/
	@echo "build $(BINARY_NAME) successfully."

tx_parser_win: BINARY_NAME=tx_parser.exe
tx_parser_win:
	export GOOS=windows
	export GOARCH=amd64
	$(GO_BUILD) -o $(BINARY_NAME) main.go
	mkdir -p bin/win
	mv $(BINARY_NAME) bin/win/
	@echo "build $(BINARY_NAME) successfully."

# default
default: tx_parser_linux
