##
## QUALITY
## -------
## 
unit: ## run unit tests
	go test -race ./...

integration: ## run integration tests
	go test -race --tags=integrations ./...

lint: generate **/*.go ## lint code
	golangci-lint -v run --build-tags=integration

test: lint unit integration ## run all tests

##
## RUN
## ---
##
run: lint generate
	go run cmd/api/main.go
##
## BUILD
## -----
##
generate: storage/migrations/*.sql storage/queries.sql ## cenerate models from the SQL schema
	sqlc generate
#    PGPASSWORD=root psql -h localhost -d hunt -U root -w -f models/schema.sql

fixture: ## import fixtures
    PGPASSWORD=root psql -h localhost -d hunt -U root -w -f fixtures/testdata.sql

.PHONY=fixture

sqlc: ## install tool to run migrations
	CGO_ENABLED=1 GOOS=linux go install github.com/kyleconroy/sqlc/cmd/sqlc@latest

build: test ## build binary
	go build -o dist/local/hunt cmd/hunt/main.go

##
## DATABASE
## --------
##

tern: ## Install tool to run migrations
	go install github.com/jackc/tern/v2@latest

migrate: storage/migrations/*.sql # Run database migrations
	tern migrate --config tern.conf --migrations storage/migrations

resetdb: ## complete reset local db 
	docker compose rm -sfv adminer db
	docker volume rm -f hunt_data01
	docker compose up -d db adminer

##
## HELP
## ----
## 

help: ## Makefile help
	@grep -E '(^[a-zA-Z_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

.DEFAULT_GOAL=help
