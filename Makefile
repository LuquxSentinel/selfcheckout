build:
	@go build -o ./bin/checkout

run:build
	@./bin/checkout

test:
	@go test ./...