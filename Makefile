GO := go
GOFLAGS :=
BUILD_DIR := build
VENDOR_DIR := vendor

.PHONY: build
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/your-app-name .

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
	@echo "  build           : Build the project"
	@echo "  vendor          : Vendor the libraries"
	@echo "  clean           : Clean the build artifacts"
	@echo "  help            : Show this help message"

.DEFAULT_GOAL := help
