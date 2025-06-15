tests:
	go test ./...

# Starts a database container for local development
run-db:
	./db/local/start-local-db.sh cosmos-mongo 27017 cosmos-mongo-data

migrate:
	migrate -source file://./db/migrations -database mongodb://localhost:27017/cosmos up

migrate-down:
	migrate -source file://./db/migrations -database mongodb://localhost:27017/cosmos down