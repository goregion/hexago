# Dockerfile for binance-tick-consumer service
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binance-tick-consumer service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o binance-tick-consumer ./cmd/binance-tick-consumer

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy the binary
COPY --from=builder /app/binance-tick-consumer .

# Change to non-root user
USER appuser

# This service typically doesn't expose ports as it's a consumer
# EXPOSE 8080

CMD ["./binance-tick-consumer"]