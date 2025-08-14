.PHONY: build run clean infra-up infra-down migrate-up migrate-down migrate-reset producer all lint test

NAME=wb_order_service
MAIN=cmd/wb_order_service/main.go
MIGRATE=cmd/migrate/main.go
PRODUCER=scripts/producer.go

build:
	go build -o $(NAME) $(MAIN)

run: build
	./$(NAME)

producer:
	go run $(PRODUCER)

all: infra-up
	sleep 5
	$(MAKE) run

lint:
	golangci-lint run ./...

test:
	go test ./...

infra-up:
	docker-compose up -d

infra-down:
	docker-compose down

migrate-up:
	go run $(MIGRATE) -action=up

migrate-down:
	go run $(MIGRATE) -action=down

migrate-reset:
	$(MAKE) migrate-down
	$(MAKE) migrate-up

clean:
	rm -f $(NAME)
