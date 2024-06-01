generate:
	buf generate ./grpc/
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
docker:
	# docker-compose -f docker-compose.dev.yml up
	docker-compose up