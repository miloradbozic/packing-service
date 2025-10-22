# Packing Service

A Go-based HTTP API service that calculates optimal packing solutions for customer orders.

## Problem Statement

The service calculates the optimal number of packs needed to fulfill customer orders with the following rules:
1. Only whole packs can be sent
2. Send out the least amount of items to fulfill the order
3. Send out as few packs as possible

## Features

- Calculates minimum items needed to fulfill orders
- Minimizes number of packs used
- Web UI for easy interaction
- RESTful API for programmatic access
- PostgreSQL database storage for pack sizes
- CRUD API for managing pack sizes
- Database migrations for schema management
- Docker support with PostgreSQL

## Prerequisites

- **Go 1.21 or later** - [Download here](https://golang.org/dl/)
- **Docker and Docker Compose** - [Download here](https://docs.docker.com/get-docker/)

## Quick Start

### 🚀 One-Command Setup (Recommended)
```bash
# Clone the repository
git clone https://github.com/miloradbozic/packing-service.git
# Alternatively, clone the repository with ssh
git clone git@github.com:miloradbozic/packing-service.git

# Move to the service repo
cd packing-service

# Start everything (database + service)
make dev
```

Visit http://localhost:8080 for the web UI

> **Note:** If you get "address already in use" error, make sure port 8080 is free or stop any services using it.

### Alternative Setup Options

#### Run locally (without database)
```bash
go run main.go
```

#### Run with Docker (includes PostgreSQL)
```bash
docker-compose up
```

#### Run locally with PostgreSQL
```bash
# Start PostgreSQL
docker-compose up -d postgres

# Run migrations
make migrate

# Start the service
go run main.go
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

## Pack Size Management API

### List All Pack Sizes

**Endpoint:** `GET /api/v1/pack-sizes`

**Response:**
```json
{
  "pack_sizes": [
    {
      "id": 1,
      "size": 250,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### Create Pack Size

**Endpoint:** `POST /api/v1/pack-sizes`

**Request:**
```json
{
  "size": 750,
}
```

### Update Pack Size

**Endpoint:** `PUT /api/v1/pack-sizes/{id}`

**Request:**
```json
{
  "size": 750,
}
```

### Delete Pack Size

**Endpoint:** `DELETE /api/v1/pack-sizes/{id}`

**Response:** `204 No Content`

## Configuration

### Database Configuration

The service now uses PostgreSQL for storing pack sizes. Database settings can be configured in `config.yaml`:

```yaml
database:
  host: "localhost"
  port: 5432
  user: "packing_user"
  password: "packing_password"
  dbname: "packing_service"
  sslmode: "disable"
  max_open_conns: 25
  max_idle_conns: 25
  conn_max_lifetime: "5m"
```

### Legacy Configuration

Pack sizes in `config.yaml` are now used only for initial migration:
```yaml
packs:
  sizes:
    - 250
    - 500
    - 1000
    - 2000
    - 5000
```

## Database Management

### Run Migrations
```bash
make migrate
```

### Development Setup
```bash
make dev  # Starts PostgreSQL, runs migrations, and starts the service
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