.PHONY: bootstrap add-migration migrate-db create-run-db brun import-data generate-coverage test migrate-db-test brun-elastic

test_db_cred := DB_HOST=localhost DB_USER=postgres DB_PASS=postgres DB_NAME=got_test DB_PORT=5433
dev_db_cred := DB_HOST=localhost DB_USER=postgres DB_PASS=postgres DB_NAME=got DB_PORT=5433
goose_dbstring := GOOSE_DBSTRING="postgres://postgres:postgres@gotdb:5433/got_test?sslmode=disable"

bootstrap:
	go mod download
	go install github.com/matryer/moq@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	@make create-run-db
	@make migrate-db

add-migration:
	goose -dir migrations/ create $(name) sql
migrate-db:
	docker-compose -f docker-compose.local-dev.yml up migrate --build
migrate-db-test:
	$(goose_dbstring) docker-compose -f docker-compose.local-dev.yml up migrate --build

create-run-db:
	docker-compose -f docker-compose.local-dev.yml up db --build -d

brun:
	docker-compose -f docker-compose.local-dev.yml up got-app --build

import-data:
	$(dev_db_cred) go run cmd/import/main.go

generate-coverage:
	(set -a; source local.test; set +a; go test -coverprofile=coverage.out `go list ./... | grep -v ./mocks`)
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out | grep total
    
test:
	$(test_db_cred) go test ./...

make brun-elastic:
	curl -o jdbc/postgresql-42.2.18.jar https://jdbc.postgresql.org/download/postgresql-42.2.18.jar
	docker-compose -f docker-compose.local-dev.yml up elasticsearch -d
	docker-compose -f docker-compose.local-dev.yml up logstash -d
