#!/bin/bash

set -o allexport
source .env
set +o allexport

cp spec/authenticator.spec.yml swagger/

docker compose --env-file .env -f docker-compose-local.yml up -d

sleep 2

goapi-gen --package=spec --out spec/authenticator.gen.spec.go spec/authenticator.spec.yml
printf " \033[0;32m✔\033[0m API Specs \n"

tern migrate --migrations migrations --config migrations/tern.conf
printf " \033[0;32m✔\033[0m Migrations runned \n"

sqlc generate -f ./sqlc.yml
printf " \033[0;32m✔\033[0m Queries SQL compiled \n"

go run main.go
