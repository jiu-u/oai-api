.PHONY: init
init:
	go install github.com/google/wire/cmd/wire@latest


.PHONY: wire-server
wire-server:
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


.PHONY: server
server:
	go run ./cmd/api_server/main.go