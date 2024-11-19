.PHONY: dev qc prod

# Run the project in development mode
dev:
	@echo "Running in development mode..."
	@go run main.go

# Run the project in quality check (qc) mode
qc:
	@echo "Running in quality check mode..."
	@GIN_MODE=release go run main.go

# Run the project in production mode
prod:
	@echo "Running in production mode..."
	@GIN_MODE=release go build -o app . && ./app
