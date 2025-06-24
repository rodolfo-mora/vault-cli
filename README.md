# Vault CLI

A comprehensive CLI tool for encrypting and decrypting data using HashiCorp Vault's transit engine. This tool supports both string and file operations with bulk processing capabilities, making it suitable for handling large files and directories with hundreds of files.

## Features

- **String Encryption/Decryption**: Encrypt or decrypt individual strings
- **File Encryption/Decryption**: Process entire files with bulk operations
- **Bulk Processing**: Efficiently handle large files by processing in configurable batches
- **Flexible Output**: Choose to override original files or specify output directories
- **Environment Variable Support**: Use standard Vault environment variables
- **Base64 Encoding**: Automatic base64 encoding/decoding for proper data handling

## Installation

### Option 1: Build from source

1. Clone the repository:
```bash
git clone <repository-url>
cd vault-cli
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the application:
```bash
go build -o vault-cli
```

### Option 2: Using Makefile

1. Clone the repository:
```bash
git clone <repository-url>
cd vault-cli
```

2. Build using Makefile:
```bash
make build
```

3. Install to system (optional):
```bash
make install
```

### Option 3: Using Docker

1. Build the Docker image:
```bash
docker build -t vault-cli .
```

2. Run the container:
```bash
docker run --rm vault-cli --help
```

### Option 4: Using Docker Compose (with Vault server)

1. Start the development environment:
```bash
docker-compose up -d vault
```

2. Wait for Vault to be ready, then run the CLI:
```bash
docker-compose run --rm vault-cli encrypt --transit-key test-key --text "Hello, World!"
```

## Configuration

The tool supports both command-line arguments and environment variables for Vault configuration:

### Environment Variables
- `VAULT_ADDR`: Vault server address (e.g., `http://127.0.0.1:8200`)
- `VAULT_TOKEN`: Vault authentication token

### Required Parameters
- `--transit-key`: The transit key name to use for encryption/decryption

## Usage

### String Operations

#### Encrypt a string:
```bash
# Using command-line arguments
./vault-cli encrypt \
  --vault-addr http://127.0.0.1:8200 \
  --vault-token your-token \
  --transit-key my-key \
  --text "Hello, World!"

# Using environment variables
export VAULT_ADDR=http://127.0.0.1:8200
export VAULT_TOKEN=your-token
./vault-cli encrypt --transit-key my-key --text "Hello, World!"

# Using Docker
docker run --rm vault-cli encrypt \
  --vault-addr http://127.0.0.1:8200 \
  --vault-token your-token \
  --transit-key my-key \
  --text "Hello, World!"
```

#### Decrypt a string:
```bash
./vault-cli decrypt \
  --vault-addr http://127.0.0.1:8200 \
  --vault-token your-token \
  --transit-key my-key \
  --text "vault:v1:encrypted-data-here"
```

### File Operations

#### Encrypt a file:
```bash
# Encrypt and create new file with .encrypted extension
./vault-cli encrypt \
  --vault-addr http://127.0.0.1:8200 \
  --vault-token your-token \
  --transit-key my-key \
  --file input.txt

# Encrypt and save to specific output folder
./vault-cli encrypt \
  --vault-addr http://127.0.0.1:8200 \
  --vault-token your-token \
  --transit-key my-key \
  --file input.txt \
  --output ./encrypted/

# Override original file
./vault-cli encrypt \
  --vault-addr http://127.0.0.1:8200 \
  --vault-token your-token \
  --transit-key my-key \
  --file input.txt \
  --override

# Using Docker with volume mount
docker run --rm -v $(pwd):/app/data vault-cli encrypt \
  --vault-addr http://127.0.0.1:8200 \
  --vault-token your-token \
  --transit-key my-key \
  --file /app/data/input.txt \
  --output /app/data/encrypted/
```

#### Decrypt a file:
```bash
# Decrypt and create new file (removes .encrypted extension)
./vault-cli decrypt \
  --vault-addr http://127.0.0.1:8200 \
  --vault-token your-token \
  --transit-key my-key \
  --file input.txt.encrypted

# Decrypt and save to specific output folder
./vault-cli decrypt \
  --vault-addr http://127.0.0.1:8200 \
  --vault-token your-token \
  --transit-key my-key \
  --file input.txt.encrypted \
  --output ./decrypted/

# Override original file
./vault-cli decrypt \
  --vault-addr http://127.0.0.1:8200 \
  --vault-token your-token \
  --transit-key my-key \
  --file input.txt.encrypted \
  --override
```

## Advanced Options

### Batch Processing
Control the number of items processed per batch for large files:

```bash
./vault-cli encrypt \
  --vault-addr http://127.0.0.1:8200 \
  --vault-token your-token \
  --transit-key my-key \
  --file large-file.txt \
  --max-items 500
```

### Custom Transit Path
If your transit engine is mounted at a different path:

```bash
./vault-cli encrypt \
  --vault-addr http://127.0.0.1:8200 \
  --vault-token your-token \
  --transit-path my-transit \
  --transit-key my-key \
  --text "Hello, World!"
```

## Development

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

## Best Practices

### For Large Files
1. **Adjust batch size**: Use `--max-items` to optimize performance based on your Vault server capacity
2. **Monitor memory usage**: Large files are processed line by line to minimize memory consumption
3. **Use output directories**: Specify `--output` to keep original files intact

### For Multiple Files
1. **Process in parallel**: Run multiple instances for different files
2. **Use consistent naming**: The tool automatically handles `.encrypted` extensions
3. **Backup originals**: Always keep backups before using `--override`

### Security Considerations
1. **Token management**: Use appropriate token policies with minimal required permissions
2. **Network security**: Ensure Vault communication is over HTTPS in production
3. **File permissions**: Be mindful of file permissions on encrypted output
4. **Container security**: When using Docker, avoid mounting sensitive directories

## Error Handling

The tool provides detailed error messages for common issues:
- Invalid Vault configuration
- Authentication failures
- Missing or invalid transit keys
- File access issues
- Network connectivity problems

## Examples

### Encrypting a configuration file:
```bash
./vault-cli encrypt \
  --vault-addr https://vault.company.com \
  --transit-key config-key \
  --file config.yaml \
  --output ./secure-configs/
```

### Decrypting multiple files:
```bash
for file in *.encrypted; do
  ./vault-cli decrypt \
    --vault-addr https://vault.company.com \
    --transit-key config-key \
    --file "$file" \
    --output ./decrypted/
done
```

### Processing sensitive data:
```bash
# Encrypt database credentials
./vault-cli encrypt \
  --vault-addr https://vault.company.com \
  --transit-key db-key \
  --text "postgresql://user:password@localhost:5432/db" \
  > encrypted_creds.txt
```

### Using Docker in CI/CD:
```bash
# Build and run in CI pipeline
docker build -t vault-cli .
docker run --rm vault-cli encrypt \
  --vault-addr $VAULT_ADDR \
  --vault-token $VAULT_TOKEN \
  --transit-key $TRANSIT_KEY \
  --file config.yaml \
  --output ./secure/
```

## Troubleshooting

### Common Issues

1. **"vault address is required"**: Set `VAULT_ADDR` environment variable or use `--vault-addr`
2. **"vault token is required"**: Set `VAULT_TOKEN` environment variable or use `--vault-token`
3. **"transit key is required"**: Always specify `--transit-key`
4. **"cannot specify both --file and --text"**: Choose either file or text input, not both

### Docker Issues

1. **Permission denied**: Ensure proper volume mounts and user permissions
2. **Network connectivity**: Use `--network host` or proper Docker networking
3. **File not found**: Check volume mount paths and file locations

### Performance Tips

1. **Increase batch size** for faster processing of large files
2. **Use local Vault instances** for development to reduce network latency
3. **Monitor Vault server metrics** during bulk operations
4. **Use multi-stage Docker builds** for smaller production images

## License

This project is licensed under the MIT License. 