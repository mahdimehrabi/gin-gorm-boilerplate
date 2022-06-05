include .env
ifeq ($(Environment),development)
DOCKER_COMMAND=docker-compose
else 
DOCKER_COMMAND=docker-compose -f docker-compose.prod.yml
endif
MIGRATE=${DOCKER_COMMAND} exec web migrate -path=core/migrations -database "postgres://${DBUsername}:${DBPassword}@${DBHost}:${DBPort}/${DBName}?sslmode=disable" -verbose

migrate-up:
		$(MIGRATE) up
migrate-down:
		$(MIGRATE) down 
force:
		@read -p  "Which version do you want to force?" VERSION; \
		$(MIGRATE) force $$VERSION

goto:
		@read -p  "Which version do you want to migrate?" VERSION; \
		$(MIGRATE) goto $$VERSION

drop:
		$(MIGRATE) drop

create-migration:
		@read -p  "What is the name of migration?" NAME; \
		${MIGRATE} create -ext sql -seq -dir migration  $$NAME

test-all:
	${DOCKER_COMMAND} exec web go test ./tests/...

test-all-debugger:
	${DOCKER_COMMAND} exec web dlv test ./tests --headless --listen=:4000 --api-version=2 --accept-multiclient 

kill-test-debugger:
	${DOCKER_COMMAND} exec web pkill -f "dlv test"

create-admin:
	${DOCKER_COMMAND} exec web go run ./cmd/. create_admin

.PHONY: migrate-up migrate-down force goto drop create

.PHONY: migrate-up migrate-down force goto drop create auto-create
