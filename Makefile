include .env

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## run: run the whole project as a docker compose command
.PHONY: run
run:
	docker compose up

## db/pslq: connect to the database
.PHONY: db/psql
db/psql:
	@echo 'Connecting to the database...'
	docker exec -it pg-db psql ${DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${PUBLIC_DB_URL} up

## db/migrations/down n=$1: apply n down database migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Running ${n} down migrations'
	migrate -path ./migrations -database ${PUBLIC_DB_URL} down ${n}

## db/migrations/version: print the current database version
.PHONY: db/migrations/version
db/migrations/version:
	@echo 'Getting the current database version...'
	migrate -path ./migrations -database ${PUBLIC_DB_URL} version

## db/migrations/drop: drop the database
.PHONY: db/migrations/drop
db/migrations/drop: confirm
	@echo 'Dropping the database...'
	migrate -path ./migrations -database ${PUBLIC_DB_URL} drop

## db/migrations/force n=$1: force the database to the latest version
.PHONY: db/migrations/force
db/migrations/force: confirm
	@echo 'Forcing the database to the ${n} version...'
	migrate -path ./migrations -database ${PUBLIC_DB_URL} force ${n}

## auth/run: run the auth service
.PHONY: auth/run
auth/run:
	cd ./auth && DATABASE_URL=${PUBLIC_DB_URL} JWT_SECRET=${JWT_SECRET} ~/go/bin/air