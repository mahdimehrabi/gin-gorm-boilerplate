include .env
MIGRATE=docker-compose exec web migrate -path=migration -database "postgres://${DBUsername}:${DBPassword}@${DBHost}:${DBPort}/${DBName}?sslmode=disable" -verbose

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

crud:
	bash automate/scripts/crud.sh

test-all:
	docker-compose exec web go test ./tests/...

test-all-debugger:
	docker-compose exec web dlv test ./tests --headless --listen=:4000 --api-version=2 --accept-multiclient 

kill-test-debugger:
	docker-compose exec web pkill -f "dlv test"

.PHONY: migrate-up migrate-down force goto drop create

.PHONY: migrate-up migrate-down force goto drop create auto-create