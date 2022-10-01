postgres:
	docker run --name postgres12 -e POSTGRES_PASSWORD=tempPassword -d -p 5432:5432 postgres 

createdb:
	docker exec -it postgres12 createdb --username=postgres --owner=postgres lift_tracker

dropdb:
	docker exec -it postgres12 dropdb --username=postgres lift_tracker

migrateup:
	migrate -path db/migration -database "postgresql://postgres:tempPassword@localhost:5432/lift_tracker?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:tempPassword@localhost:5432/lift_tracker?sslmode=disable" -verbose down

sqlc:
	docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc generate

pgshell:
	docker exec -it postgres12 bash

test: 
	go test -v -cover ./...

ctest:
	go clean -testcache && go test -v -cover ./...

fmt:
	go fmt ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc pgshell

