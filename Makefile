PROJECT_NAME := graphql-federation
WASM_FILE := $(PROJECT_NAME).wasm
MAIN_FILE := cmd/graphql/graphql-federation.go

# Go variables
GO_VERSION := 1.24
GO := go

# Build directories
BUILD_DIR := build
DIST_DIR := dist


.PHONY: all build clean

# Default target
all: build


# Build the WASM file using Go wasip1 mode
build:
	@echo "Building WASM file with Go wasip1 mode..."
	@mkdir -p $(BUILD_DIR)
	env GOOS=wasip1 GOARCH=wasm $(GO) build -buildmode=c-shared -o $(BUILD_DIR)/$(WASM_FILE) $(MAIN_FILE)
	@echo "WASM file built: $(BUILD_DIR)/$(WASM_FILE)"
	@ls -lh $(BUILD_DIR)/$(WASM_FILE)

# Create distribution package
dist: build
	@echo "Creating distribution package..."
	@mkdir -p $(DIST_DIR)
	@cp $(BUILD_DIR)/$(WASM_FILE) $(DIST_DIR)/
	@cp README.md $(DIST_DIR)/ 2>/dev/null || true
	@echo "Distribution package created in $(DIST_DIR)/"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	@echo "Clean completed"
