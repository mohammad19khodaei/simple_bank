postgres:
	docker container run -d --name postgres17 --network simplebank-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 postgres:17-alpine

createdb:
	docker exec postgres17 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec postgres17 dropdb --username=root simple_bank

migrateup:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/mohammad19khodaei/simple_bank/db/sqlc Store