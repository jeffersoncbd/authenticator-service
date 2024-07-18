#!/bin/bash

set -o allexport
source .env
set +o allexport

cp internal/spec/authenticator.spec.yml swagger/

docker compose --env-file .env -f docker-compose-local.yml up -d

sleep 2

goapi-gen --package=spec --out internal/spec/authenticator.gen.spec.go internal/spec/authenticator.spec.yml
printf " \033[0;32m✔\033[0m API Specs \n"

tern migrate --migrations internal/databases/postgresql/migrations --config internal/databases/postgresql/migrations/tern.conf
printf " \033[0;32m✔\033[0m Migrations runned \n"

sqlc generate -f internal/databases/postgresql/sqlc.yml
printf " \033[0;32m✔\033[0m Queries SQL compiled \n"

go run application/application.go
