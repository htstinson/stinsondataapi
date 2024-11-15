.PHONY: setup run clean test

setup:
	@mkdir -p certs
	@echo "Generating self-signed certificate..."
	@openssl req -x509 \
		-newkey rsa:4096 \
		-keyout certs/local.key \
		-out certs/local.crt \
		-days 365 \
		-nodes \
		-subj "/CN=api.local.dev" 2>/dev/null
	@echo "Starting database..."
	@docker-compose up -d db
	@echo "Installing dependencies..."
	@go mod tidy

run:
	@go run cmd/api/main.go

clean:
	@docker-compose down -v
	@rm -rf certs

test:
	@go test ./...