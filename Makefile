postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=test -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:test@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migrateupall:
	migrate -path db/migration -database "postgresql://postgres:test@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:test@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

migratedownall:
	migrate -path db/migration -database "postgresql://postgres:test@localhost:5432/simple_bank?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

server:
	go run main.go

testdb:
	mockgen -package testdb -destination db/test/store.go github.com/sssaang/simplebank/db/sqlc Store

.PHONY: createdb dropdb postgres migrateup migrateupall migratedown migratedownall test server testdb