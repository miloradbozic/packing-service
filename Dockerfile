# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o packing-service main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary and config
COPY --from=builder /app/packing-service .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/templates ./templates

EXPOSE 8080

CMD ["./packing-service"]