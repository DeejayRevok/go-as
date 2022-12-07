setup-env:
	@cp .env.example .env
	@docker-compose build --build-arg DEVELOPMENT=false as

lint:
	@docker-compose run --rm as go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run -v

test:
	@docker-compose run --rm as go test ./...
