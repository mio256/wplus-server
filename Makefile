BIN_DIR:=$(shell pwd)/bin

.PHONY: test
test:
	go test ./...

.PHONY: tools
tools:
	go get .
	GOBIN=$(BIN_DIR) go install github.com/sqldef/sqldef/cmd/psqldef@$(shell go list -m -f "{{.Version}}" github.com/sqldef/sqldef)
	GOBIN=$(BIN_DIR) go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(shell go list -m -f "{{.Version}}" github.com/sqlc-dev/sqlc)

.PHONY: migrate-db
migrate-db:
	$(BIN_DIR)/psqldef -U postgres -W postgres -p 5432 -f ./db/core.sql --enable-drop-table wplus

.PHONY: dry-migrate-db
dry-migrate-db:
	$(BIN_DIR)/psqldef -U postgres -W postgres -p 5432 --dry-run -f ./db/core.sql --enable-drop-table wplus

.PHONY: sqlc-gen
sqlc-gen:
	$(BIN_DIR)/sqlc generate

.PHONY: sqlc-lint
sqlc-lint:
	$(BIN_DIR)/sqlc vet
