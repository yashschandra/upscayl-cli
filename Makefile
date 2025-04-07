# Project settings
APP_NAME := upscayl
OUTPUT_DIR := release/build
MAIN_PKG := .

# Default target
all: clean linux mac-amd mac-arm

# Linux AMD64 build
linux:
	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_PKG)
	@tar -czf $(OUTPUT_DIR)/$(APP_NAME)-linux-amd64.tar.gz -C $(OUTPUT_DIR) $(APP_NAME)-linux-amd64
	@rm $(OUTPUT_DIR)/$(APP_NAME)-linux-amd64

# macOS AMD64 build
mac-amd:
	@echo "Building for macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(APP_NAME)-darwin-amd64 $(MAIN_PKG)
	@tar -czf $(OUTPUT_DIR)/$(APP_NAME)-darwin-amd64.tar.gz -C $(OUTPUT_DIR) $(APP_NAME)-darwin-amd64
	@rm $(OUTPUT_DIR)/$(APP_NAME)-darwin-amd64

# macOS ARM64 build
mac-arm:
	@echo "Building for macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(APP_NAME)-darwin-arm64 $(MAIN_PKG)
	@tar -czf $(OUTPUT_DIR)/$(APP_NAME)-darwin-arm64.tar.gz -C $(OUTPUT_DIR) $(APP_NAME)-darwin-arm64
	@rm $(OUTPUT_DIR)/$(APP_NAME)-darwin-arm64

# Clean
clean:
	@echo "Cleaning..."
	@rm -rf $(OUTPUT_DIR)
	@mkdir -p $(OUTPUT_DIR)

.PHONY: all linux mac-amd mac-arm clean
