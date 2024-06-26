# Load environment variables from .env file
include ./config/.env
export

.PHONY: all build run migrate clean

all: build run

run-dev:
	go run ./main.go

build:
	go build -o blockchain-medical-record .

clean:
	rm -f blockchain-medical-record

migrate-up:
	migrate -database "postgres://$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)" -source=file://./datasource/migrations up

migrate-down:
	migrate -database "postgres://$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)" -source=file://./datasource/migrations down