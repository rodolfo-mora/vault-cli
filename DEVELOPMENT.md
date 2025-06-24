# Development Guide

This document provides instructions for developing and testing the Vault CLI application.

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose (optional)
- Make (optional)

## Development Setup

### Using Makefile

The project includes a comprehensive Makefile with the following targets:

```bash
# Build for current platform
make build-local

# Build for all platforms
make build-all

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Run linter
make lint

# Clean build artifacts
make clean

# Install dependencies
make deps

# Show all available targets
make help
```

### Using Docker for Development

1. Start the development environment with Vault:
```bash
docker-compose up -d vault
```

2. Create a transit key in Vault:
```bash
docker-compose exec vault vault secrets enable transit
docker-compose exec vault vault write -f transit/keys/test-key
```

3. Test the CLI:
```bash
docker-compose run --rm vault-cli encrypt --transit-key test-key --text "Hello, World!"
```

### Building for Different Platforms

```bash
# Build for Linux
make build-linux

# Build for macOS
make build-darwin

# Build for Windows
make build-windows

# Build for all platforms
make build-all
```

## Docker Usage

### Building the Image

```bash
# Build the image
docker build -t vault-cli .

# Build with specific tag
docker build -t vault-cli:v1.0.0 .
```

### Running the Container

```bash
# Show help
docker run --rm vault-cli --help

# Encrypt a string
docker run --rm vault-cli encrypt \
  --vault-addr http://127.0.0.1:8200 \
  --vault-token your-token \
  --transit-key my-key \
  --text "Hello, World!"

# Encrypt a file (with volume mount)
docker run --rm -v $(pwd):/app/data vault-cli encrypt \
  --vault-addr http://127.0.0.1:8200 \
  --vault-token your-token \
  --transit-key my-key \
  --file /app/data/input.txt \
  --output /app/data/encrypted/
```

### Using Docker Compose

The `docker-compose.yml` file provides a complete development environment:

```bash
# Start Vault server
docker-compose up -d vault

# Wait for Vault to be ready, then run commands
docker-compose run --rm vault-cli encrypt --transit-key test-key --text "Hello, World!"
docker-compose run --rm vault-cli decrypt --transit-key test-key --text "vault:v1:encrypted-data"
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test -v ./vault
```

### Integration Testing

1. Start Vault server:
```bash
docker-compose up -d vault
```

2. Set up transit key:
```bash
docker-compose exec vault vault secrets enable transit
docker-compose exec vault vault write -f transit/keys/test-key
```

3. Run integration tests:
```bash
# Test string encryption
docker-compose run --rm vault-cli encrypt --transit-key test-key --text "test-data"

# Test file encryption
echo "test-data" > test.txt
docker-compose run --rm -v $(pwd):/app/data vault-cli encrypt \
  --transit-key test-key --file /app/data/test.txt --output /app/data/
```

## Code Quality

### Formatting

```bash
# Format code
make fmt

# Or manually
go fmt ./...
```

### Linting

```bash
# Run linter (requires golangci-lint)
make lint

# Install golangci-lint if not available
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Building

```bash
# Build for current platform
make build-local

# Build for all platforms
make build-all

# Create release packages
make release
```

## Release Process

1. Update version in code (if using version variables)
2. Create a git tag:
```bash
git tag v1.0.0
git push origin v1.0.0
```

3. Build release packages:
```bash
make release
```

4. The release packages will be created in `build/release/`

## Troubleshooting

### Docker Issues

1. **Permission denied**: Ensure proper volume mounts and user permissions
2. **Network connectivity**: Use `--network host` or proper Docker networking
3. **File not found**: Check volume mount paths and file locations

### Build Issues

1. **Dependencies**: Run `make deps` to ensure all dependencies are downloaded
2. **Platform-specific builds**: Use the appropriate make target for your target platform
3. **Clean build**: Run `make clean` to remove any cached build artifacts

### Vault Issues

1. **Connection refused**: Ensure Vault server is running and accessible
2. **Authentication failed**: Check token validity and permissions
3. **Transit key not found**: Ensure the transit engine is enabled and the key exists 