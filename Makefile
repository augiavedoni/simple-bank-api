postgres:
	docker run --name postgres-1 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=PaSSw0rD -d postgres

createdb:
	docker exec -it postgres-1 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres-1 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:PaSSw0rD@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:PaSSw0rD@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc