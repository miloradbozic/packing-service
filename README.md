# Packing Service

A Go-based HTTP API service that calculates optimal packing solutions for customer orders.

## Problem Statement

The service calculates the optimal number of packs needed to fulfill customer orders with the following rules:
1. Only whole packs can be sent
2. Send out the least amount of items to fulfill the order
3. Send out as few packs as possible

## Features

- ✅ Calculates minimum items needed to fulfill orders
- ✅ Minimizes number of packs used
- ✅ Web UI for easy interaction
- ✅ RESTful API for programmatic access
- ✅ Configurable pack sizes via YAML
- ✅ Docker support

## Quick Start

### Run locally
```bash
go run main.go
```

Visit http://localhost:8080 for the web UI

### Run with Docker
```bash
docker-compose up
```

### Run tests
```bash
make test
```

## API Documentation

### Calculate Packing

**Endpoint:** `POST /api/v1/calculate`

**Request:**
```json
{
  "items": 501
}
```

**Response:**
```json
{
  "items_ordered": 501,
  "total_items_shipped": 750,
  "total_packs": 2,
  "packs": [
    {
      "size": 500,
      "quantity": 1
    },
    {
      "size": 250,
      "quantity": 1
    }
  ],
  "excess_items": 249
}
```

### Get Configuration

**Endpoint:** `GET /api/v1/config`

**Response:**
```json
{
  "pack_sizes": [250, 500, 1000, 2000, 5000]
}
```

## Configuration

Pack sizes can be modified in `config.yaml`:
```yaml
packs:
  sizes:
    - 250
    - 500
    - 1000
    - 2000
    - 5000
```

## Test Examples

| Items Ordered | Packs Sent | Total Items |
|--------------|------------|-------------|
| 1 | 1×250 | 250 |
| 250 | 1×250 | 250 |
| 251 | 1×500 | 500 |
| 501 | 1×500, 1×250 | 750 |
| 12001 | 2×5000, 1×2000, 1×250 | 12250 |

## Deployment

The service can be deployed to any platform that supports Docker or Go binaries (Heroku, Railway, Google Cloud Run, etc.)