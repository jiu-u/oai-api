.PHONY: init
init:
	go install github.com/google/wire/cmd/wire@latest

.PHONY: wire-server
wire-server:
	wire gen ./cmd/server/wire

.PHONY: wire-api
wire-api:
	wire gen ./cmd/api_server/wire/
	wire gen ./cmd/api_server/wire_load/

.PHONY: wire-migrate
wire-migrate:
	wire gen ./cmd/migrate/wire

.PHONY: wire-load
wire-load:
	wire gen ./cmd/dataload/wire/

.PHONY: migrate
migrate:
	go run ./cmd/migrate/main.go

.PHONY: load
load:
	go run ./cmd/dataload/

.PHONY: api
api:
	go run ./cmd/api_server/main.go

.PHONY: server
server:
	go run ./cmd/server/main.go

.PHONY: docker
docker:
	docker build --build-arg APP_CONF=config/prod.yaml --build-arg  APP_RELATIVE_PATH=./cmd/server/ -t jiu/oai-api:v1 .
	docker run --rm -p 18081:8080 jiu/oai-api:v1

.PHONY: build-docker
build-docker:
	docker build --build-arg APP_CONF=config/prod.yaml --build-arg  APP_RELATIVE_PATH=./cmd/server/ -t jiu/oai-api:v1 .
	#docker build -t jiu/oai-api:v1 --build-arg APP_CONF=config/config.yaml --build-arg  APP_RELATIVE_PATH=./cmd/server/ .

.PHONY: build
build:
	go build -ldflags="-s -w" -trimpath -o ./data/bin/server ./cmd/server/