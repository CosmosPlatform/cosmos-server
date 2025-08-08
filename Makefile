tests:
	go test ./...

# Starts a database container for local development
run-db:
	./db/local/start-local-db.sh cosmos-postgres 5432 cosmos-postgres-data

migrate:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/cosmos?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/cosmos?sslmode=disable" down

migrate_uri_up:
	migrate -path db/migrations -database "$(uri)" up

migrate_uri_down:
	migrate -path db/migrations -database "$(uri)" down

