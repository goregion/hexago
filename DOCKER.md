# Docker Instructions for Hexago

This document provides instructions for building and running the Hexago application using Docker.

## Prerequisites

- Docker
- Docker Compose

## Project Structure

The project includes the following Docker files:

```
Dockerfile                           # Generic multi-stage Dockerfile
docker/all-in-one.Dockerfile        # All-in-one service
docker/backoffice-api.Dockerfile    # Backoffice API service
docker/binance-tick-consumer.Dockerfile  # Binance tick consumer
docker/ohlc-generator.Dockerfile    # OHLC generator
docker-compose.yml                   # Docker Compose configuration
.dockerignore                        # Docker ignore file
```

## Building Individual Services

### Build all-in-one service
```bash
docker build -f docker/all-in-one.Dockerfile -t hexago-all-in-one .
```

### Build backoffice-api service
```bash
docker build -f docker/backoffice-api.Dockerfile -t hexago-backoffice-api .
```

### Build binance-tick-consumer service
```bash
docker build -f docker/binance-tick-consumer.Dockerfile -t hexago-binance-tick-consumer .
```

### Build ohlc-generator service
```bash
docker build -f docker/ohlc-generator.Dockerfile -t hexago-ohlc-generator .
```

### Build using generic Dockerfile with target
```bash
# Build any service using the generic Dockerfile
docker build --build-arg TARGET_CMD=all-in-one -t hexago-all-in-one .
docker build --build-arg TARGET_CMD=backoffice-api -t hexago-backoffice-api .
docker build --build-arg TARGET_CMD=binance-tick-consumer -t hexago-binance-tick-consumer .
docker build --build-arg TARGET_CMD=ohlc-generator -t hexago-ohlc-generator .
```

## Running with Docker Compose

### Option 1: All-in-one mode (default)
This runs all services in a single container along with Redis and MySQL:

```bash
docker-compose up -d
```

### Option 2: Microservices mode
This runs each service in its own container:

```bash
docker-compose --profile microservices up -d
```

### Stop services
```bash
docker-compose down
```

### Stop and remove volumes
```bash
docker-compose down -v
```

## Environment Variables

The following environment variables can be configured:

- `BINANCE_API_KEY`: Your Binance API key (optional for consumer service)
- `BINANCE_SECRET_KEY`: Your Binance secret key (optional for consumer service)
- `REDIS_URL`: Redis connection URL (default: redis:6379)
- `MYSQL_DSN`: MySQL Data Source Name (format: user:password@tcp(host:port)/database?options)
- `SYMBOLS`: Comma-separated list of trading symbols (e.g., BTCUSDT,ETHUSDT,ADAUSDT)

### Create a .env file for environment variables:
```bash
echo "BINANCE_API_KEY=your_api_key_here" > .env
echo "BINANCE_SECRET_KEY=your_secret_key_here" >> .env
echo "SYMBOLS=BTCUSDT,ETHUSDT,ADAUSDT" >> .env
```

## Service Ports

- **All-in-one**: 8080, 9090
- **Backoffice API**: 9091
- **Redis**: 6379
- **MySQL**: 3306

## Health Checks

The docker-compose.yml includes health checks for:
- Redis: `redis-cli ping`
- MySQL: `mysqladmin ping`

## Logs

View logs for all services:
```bash
docker-compose logs -f
```

View logs for a specific service:
```bash
docker-compose logs -f all-in-one
docker-compose logs -f backoffice-api
docker-compose logs -f binance-tick-consumer
docker-compose logs -f ohlc-generator
```

## Troubleshooting

1. **Build issues**: Make sure all dependencies are properly defined in go.mod
2. **Connection issues**: Ensure all services are healthy before dependent services start
3. **Port conflicts**: Check if ports 3306, 6379, 8080, 9090, 9091 are available
4. **Permission issues**: The containers run as non-root user (appuser) for security

## Development

For development, you can mount your source code as a volume:

```yaml
volumes:
  - .:/app
```

Add this to any service in docker-compose.yml for live code reloading during development.