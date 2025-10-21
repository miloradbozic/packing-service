.PHONY: build test run docker-build docker-run clean

build:
	go build -o bin/packing-service main.go

test:
	go test ./... -v

run:
	go run main.go

docker-build:
	docker build -t packing-service .

docker-run:
	docker-compose up

docker-down:
	docker-compose down

clean:
	rm -rf bin/

lint:
	go fmt ./...
	go vet ./...

coverage:
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out