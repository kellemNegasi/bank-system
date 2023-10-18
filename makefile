postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root basic_bank
dropdb:
	docker exec -it postgres12 dropdb basic_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/basic_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/basic_bank?sslmode=disable" -verbose down
migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/basic_bank?sslmode=disable" -verbose up 1
migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/basic_bank?sslmode=disable" -verbose down 1

test:
	go test -v -cover ./...
sqlc:
	sqlc generate
server:
	go run main.go
mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/kellemNegasi/bank-system/db/sqlc Store
.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc server mock