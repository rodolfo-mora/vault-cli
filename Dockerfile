# Multi-stage build for vault-cli
# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o vault-cli .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S vaultcli && \
    adduser -u 1001 -S vaultcli -G vaultcli

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/vault-cli .

# Change ownership to non-root user
RUN chown -R vaultcli:vaultcli /app

# Switch to non-root user
USER vaultcli

# Set environment variables
ENV VAULT_ADDR=""
ENV VAULT_TOKEN=""

# Expose any necessary ports (if needed for future features)
# EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["./vault-cli"]

# Default command (can be overridden)
CMD ["--help"] 