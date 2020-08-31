# defining environment variables
GOLANGCI_VERSION=1.26.0
GOLANGCI_COMMIT=6bd10d01fde78697441d9c11e2235f0dbb1e2822
PROJECT_PATH=$(shell dirname $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))
MIGRATIONS_DIR=db/migrations
DATABASE_URL=postgres://postgres:abcd@127.0.0.1:5432/test_database?sslmode=disable
name=default_migration_name

migrate-new: ## Create new migration
	dbmate --migrations-dir ${MIGRATIONS_DIR} new ${name}

migrate-create: ## Upgrade all migrations
	dbmate --migrations-dir ${MIGRATIONS_DIR} create ${name}

migrate-drop: ## Upgrade all migrations
	dbmate --migrations-dir ${MIGRATIONS_DIR} drop ${name}

migrate-up: ## Upgrade all migrations
	dbmate --migrations-dir ${MIGRATIONS_DIR} up

migrate-down: ## Downgrade all migrations
	dbmate --migrations-dir ${MIGRATIONS_DIR} down

code-quality:
	golangci-lint run 

test:
	go test -cover ./model
	go test -cover ./web/country
	go test -cover ./web/user
	go test -cover ./web/health

mock-generate:
	go generate ./...
