SHELL := /bin/zsh

.PHONY: setup run clean test

setup:
	@echo "Generating certificates..."
	@chmod +x ../scripts/generate-certs.sh
	@../scripts/generate-certs.sh
	@echo "Starting database..."
	@docker-compose up -d db
	@echo "Installing dependencies..."
	@go mod tidy
	@echo "Trust the local CA certificate..."
	@if [ "$(shell uname)" = "Darwin" ]; then \
		sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain certs/rootCA.crt; \
	elif [ "$(shell uname)" = "Linux" ]; then \
		sudo cp certs/rootCA.crt /usr/local/share/ca-certificates/ && \
		sudo update-ca-certificates; \
	else \
		echo "Please manually trust certs/rootCA.crt"; \
	fi

run:
	@if [ "$(shell uname)" = "Darwin" ] || [ "$(shell uname)" = "Linux" ]; then \
		sudo -E go run cmd/api/main.go; \
	else \
		go run cmd/api/main.go; \
	fi

clean:
	@docker-compose down -v
	@rm -rf certs

test:
	@go test ./...