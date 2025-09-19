# Dockerfile for backoffice-api service
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the backoffice-api service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o backoffice-api ./cmd/backoffice-api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy the binary
COPY --from=builder /app/backoffice-api .

# Change to non-root user
USER appuser

# Expose gRPC port (adjust as needed)
EXPOSE 9090

CMD ["./backoffice-api"]