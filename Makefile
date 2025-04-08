# Project settings
APP_NAME := upscayl
OUTPUT_DIR := release/build
MAIN_PKG := .

# Extract version from version.go
VERSION := $(shell grep 'const Version' version.go | sed 's/[^"]*"\([^"]*\)".*/\1/')

# Default target
all: clean linux mac-amd mac-arm

linux:
	@echo "Building for Linux (amd64)... Version: $(VERSION)"
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(APP_NAME) $(MAIN_PKG)
	@tar -czf $(OUTPUT_DIR)/$(APP_NAME)-v$(VERSION)-linux-amd64.tar.gz -C $(OUTPUT_DIR) $(APP_NAME)
	@rm $(OUTPUT_DIR)/$(APP_NAME)

mac-amd:
	@echo "Building for macOS (amd64)... Version: $(VERSION)"
	GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(APP_NAME) $(MAIN_PKG)
	@tar -czf $(OUTPUT_DIR)/$(APP_NAME)-v$(VERSION)-darwin-amd64.tar.gz -C $(OUTPUT_DIR) $(APP_NAME)
	@rm $(OUTPUT_DIR)/$(APP_NAME)

mac-arm:
	@echo "Building for macOS (arm64)... Version: $(VERSION)"
	GOOS=darwin GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(APP_NAME) $(MAIN_PKG)
	@tar -czf $(OUTPUT_DIR)/$(APP_NAME)-v$(VERSION)-darwin-arm64.tar.gz -C $(OUTPUT_DIR) $(APP_NAME)
	@rm $(OUTPUT_DIR)/$(APP_NAME)

local:
	@echo "Building for local..."
	go build -o $(OUTPUT_DIR)/local/$(APP_NAME) $(MAIN_PKG)

# Clean
clean:
	@echo "Cleaning..."
	@rm -rf $(OUTPUT_DIR)
	@mkdir -p $(OUTPUT_DIR)

.PHONY: all linux mac-amd mac-arm clean
