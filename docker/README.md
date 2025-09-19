# Docker Files

This directory contains individual Dockerfile configurations for each service in the hexago project.

## Files

- `all-in-one.Dockerfile` - Builds a container with all services running together
- `backoffice-api.Dockerfile` - Builds the gRPC API service
- `binance-tick-consumer.Dockerfile` - Builds the Binance WebSocket consumer service  
- `ohlc-generator.Dockerfile` - Builds the OHLC data generator service

## Usage

All Dockerfiles should be built from the project root directory using the `-f` flag:

```bash
# Build from project root
docker build -f docker/service-name.Dockerfile -t image-name .
```

## Architecture

Each Dockerfile follows the same multi-stage build pattern:

1. **Builder stage**: Uses `golang:1.24-alpine` to compile the Go binary
2. **Runtime stage**: Uses minimal `alpine:latest` with security best practices
   - Non-root user (appuser:1001)
   - CA certificates for HTTPS
   - Minimal attack surface

## Security

- All containers run as non-root user (UID 1001)
- CA certificates included for secure external connections
- Minimal Alpine Linux base image
- No unnecessary packages or tools included