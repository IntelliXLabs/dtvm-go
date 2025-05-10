.PHONY: test build clean

# Default target
all: test

# Build WASM before running tests
test:
	@echo "Building WASM files..."
	bash build_wasm.sh
	@echo "Running tests..."
	go test ./... -v

# Add a build target if needed
build:
	go build ./...

# Clean target to remove generated files
clean:
	rm -f example/*.wasm
	go clean
