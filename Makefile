# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt

# Binary names
SERVER_BINARY=chat-server
CLI_BINARY=chat-cli

CLI_INSTALL_DIR := $(HOME)/.local/bin

# Default target
all: build

# Build both binaries
build: build-server build-cli

# Build chat-server
build-server:
	cd chat-server && \
	$(GOBUILD) -o bin/$(SERVER_BINARY) .

# Build chat-cli
build-cli:
	cd chat-cli && \
	$(GOBUILD) -o bin/$(CLI_BINARY) .

# Install the CLI binary to ~/.local/bin
install-cli: build-cli
	@mkdir -p $(CLI_INSTALL_DIR)
	cp chat-cli/bin/$(CLI_BINARY) $(CLI_INSTALL_DIR)/$(CLI_BINARY)
	chmod +x $(CLI_INSTALL_DIR)/$(CLI_BINARY)
	@echo "Installed $(CLI_BINARY) to $(CLI_INSTALL_DIR)"

# Format code
fmt:
	$(GOFMT) -s -w .

# Install dependencies
deps:
	cd chat-server && \
	$(GOMOD) download && \
	$(GOMOD) tidy 

	cd chat-cli && \
	$(GOMOD) download && \
	$(GOMOD) tidy


# Run the server
run-server: build-server
	./chat-server/bin/$(SERVER_BINARY)

# Help
help:
	@echo "Available targets:"
	@echo "  build        - Build both chat-server and chat-cli"
	@echo "  build-server - Build chat-server only"
	@echo "  build-cli    - Build chat-cli only"
	@echo "  install-cli  - Install chat-cli to ~/.local/bin"
	@echo "  fmt          - Format Go code"
	@echo "  deps         - Install and tidy dependencies"
	@echo "  run-server   - Build and run chat-server"
	@echo "  run-cli      - Build and run chat-cli"
	@echo "  help         - Show this help"

.PHONY: all build build-server build-cli install-cli fmt deps clean run-server run-cli test help
