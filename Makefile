.PHONY: run swag seed test build migrate clean

# Load environment variables from .env file before running the command
load-env:
	@echo "Loading environment variables..."
	@set -a && source .env && set +a

# Run server, it will generate docs
run: swag
	go run cmd/server/main.go

# Generate docs
swag:
	swag init -g cmd/server/main.go -o ./docs/

# Migrate database
migrate:
	go run cmd/migrate/main.go

# Migrate database with clean option (drop all tables first)
clean:
	go run cmd/migrate/main.go --clean

# Seed database
seed:
	go run cmd/seed/main.go

# Test
test:
	go test ./...

# Build
build:
	go build -o bin/meet-book-api cmd/server/main.go

docker-build:
	docker build -t meet-book-api-dev .

docker-run:
	docker run -p 8080:8080 --rm meet-book-api-dev