#!make
include .env
export $(shell sed 's/=.*//' .env)
RED='\033[0;31m'        #  ${RED}
GREEN='\033[0;32m'      #  ${GREEN}
BOLD='\033[1;m'          #  ${BOLD}
WARNING=\033[37;1;41m  #  ${WARNING}
END_COLOR='\033[0m'       #  ${END_COLOR}

.PHONY: rebuild up stop restart status console-app console-db logs logs-app logs-db help

deploy: rebuild up

rebuild: stop
	@echo "\n\033[1;m Rebuilding containers... \033[0m"
	@docker-compose ${COMPOSE_PROJECT_NAME} build --no-cache

up:
	@echo "\n\033[1;m Spinning up containers for ${ENVIRONMENT} environment... \033[0m"
	@docker-compose ${COMPOSE_PROJECT_NAME} up -d
	@$(MAKE) --no-print-directory status

stop:
	@echo "\n\033[1;m Halting containers... \033[0m"
	@docker-compose  ${COMPOSE_PROJECT_NAME} stop

restart: stop up

down:
	echo "\n\033[1;m Removing containers... \033[0m"
	@docker-compose  ${COMPOSE_PROJECT_NAME} down

status:
	@echo "\n\033[1;m Containers statuses \033[0m"
	@docker-compose  ${COMPOSE_PROJECT_NAME} ps

app-console:
	docker exec -it server sh

db-console:
	docker exec -it database sh

generate-rpc:
	protoc --go_out=./pkg/grpc \
        --go-grpc_out=./pkg/grpc protos/user-service/user.proto

logs:
	@docker-compose logs --tail=1000 -f

migration-create:
#NAME - is name of migrations
	goose -dir ./pkg/migrations create $(NAME) sql

migrations:
	goose -dir ./internal/migrations postgres "user=${POSTGRES_USER} dbname=${POSTGRES_DB} host=${POSTGRES_HOST} password=${POSTGRES_PASSWORD} port=${POSTGRES_PORT} sslmode=disable" up

migrations-container:
	docker exec -it auth-service goose -dir ./pkg/migrations postgres "user=${POSTGRES_USER} dbname=${POSTGRES_DB} host=${POSTGRES_HOST} password=${POSTGRES_PASSWORD} port=${POSTGRES_PORT} sslmode=disable" up

help:
	@echo "\033[1;32mdocker-env\t\t- Main scenario, used by default\033[0m"

	@echo "\n\033[1mMain section\033[0m"
	@echo "clone\t\t\t- clone app repo"

	@echo "rebuild\t\t\t- build containers w/o cache"
	@echo "up\t\t\t- start project"
	@echo "stop\t\t\t- stop project"
	@echo "restart\t\t\t- restart containers"
	@echo "status\t\t\t- show status of containers"
