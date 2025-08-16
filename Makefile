DC_PROJECT=pc_orca
DC=docker compose -f ./docker-compose.yml -p $(DC_PROJECT)

.PHONY: build
build:
	./scripts/build.sh

.PHONY: lint-last-commit
lint-last-commit:
	npx commitlint --last

.PHONY: dc-up
dc-up:
	$(DC) up -d

.PHONY: dc-down
dc-down:
	$(DC) down --remove-orphans