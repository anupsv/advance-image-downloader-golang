GO := go
GOFLAGS :=
BUILD_DIR := build
VENDOR_DIR := vendor

.PHONY: build
build-linux-amd64:
	@echo "Building for Linux (AMD64)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/your-app-name-linux-amd64 .

.PHONY: build-linux-arm
build-linux-arm:
	@echo "Building for Linux (ARM)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=arm $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/your-app-name-linux-arm .

.PHONY: vendor
vendor:
	@echo "Vendoring libraries..."
	@$(GO) mod vendor

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build-linux-amd64   : Build for Linux (AMD64)"
	@echo "  build-linux-arm     : Build for Linux (ARM)"
	@echo "  vendor              : Vendor the libraries"
	@echo "  clean               : Clean the build artifacts"
	@echo "  help                : Show this help message"

.DEFAULT_GOAL := build-linux-amd64
