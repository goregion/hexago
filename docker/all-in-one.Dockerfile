# Dockerfile for all-in-one service
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the all-in-one service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o all-in-one ./cmd/all-in-one

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy the binary
COPY --from=builder /app/all-in-one .

# Change to non-root user
USER appuser

# Expose ports if needed (you may need to adjust based on your services)
EXPOSE 8080 9090

CMD ["./all-in-one"]