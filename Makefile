include .env

postgres:
	docker run --name postgres12 -e POSTGRES_PASSWORD=${PG_PASS} -d -p 5432:5432 ${PG_USER} 

createdb:
	docker exec -it postgres12 createdb --username=${PG_USER} --owner=${PG_USER} ${PG_DB}

dropdb:
	docker exec -it postgres12 dropdb --username=${PG_USER} ${PG_DB}

migrateup:
	migrate -path db/migration -database "postgresql://${PG_USER}:${PG_PASS}@localhost:5432/${PG_DB}?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://${PG_USER}:${PG_PASS}@localhost:5432/${PG_DB}?sslmode=disable" -verbose down

sqlc:
	docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc

