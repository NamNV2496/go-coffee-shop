generate:
	buf generate ./api/
	go generate ./...
wire:
	wire ./internal/product/
	wire ./internal/counter/
	wire ./internal/kitchen/
	wire ./internal/batch/
	wire ./internal/authorization/
run-product:
	go run ./cmd/product/main.go
run-counter:
	go run ./cmd/counter/main.go
run-kitchen:
	go run ./cmd/kitchen/main.go
run-batch:
	go run ./cmd/batch/main.go
run-authorization:
	go run ./cmd/authorization/main.go
docker:
	# docker-compose -f docker-compose.dev.yml up
	docker-compose up
lint:
	golangci-lint run

coverage:
	go test ./... --coverprofile=./coverprofile

show-coverage:
	go tool cover -html=./coverprofile