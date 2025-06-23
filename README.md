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

## Troubleshooting

### Common Issues

1. **"vault address is required"**: Set `VAULT_ADDR` environment variable or use `--vault-addr`
2. **"vault token is required"**: Set `VAULT_TOKEN` environment variable or use `--vault-token`
3. **"transit key is required"**: Always specify `--transit-key`
4. **"cannot specify both --file and --text"**: Choose either file or text input, not both

### Performance Tips

1. **Increase batch size** for faster processing of large files
2. **Use local Vault instances** for development to reduce network latency
3. **Monitor Vault server metrics** during bulk operations

## License

This project is licensed under the MIT License. 