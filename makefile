postgres:
	docker run --name postgresLatest --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:latest

createdb:
	docker exec -it postgresLatest createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgresLatest dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:hWE6fEPG6nsZ0d6jod9X@simple-bank2.c0cjzrswxdtv.eu-west-1.rds.amazonaws.com:5432/simplebank" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:hWE6fEPG6nsZ0d6jod9X@simple-bank2.c0cjzrswxdtv.eu-west-1.rds.amazonaws.com:5432/simplebank" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:hWE6fEPG6nsZ0d6jod9X@simple-bank2.c0cjzrswxdtv.eu-west-1.rds.amazonaws.com:5432/simplebank" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:hWE6fEPG6nsZ0d6jod9X@simple-bank2.c0cjzrswxdtv.eu-west-1.rds.amazonaws.com:5432/simplebank" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main/main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go simpleBank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1 migratedown1