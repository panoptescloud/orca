DC_PROJECT=pc_orca
DC=docker compose -f ./docker-compose.yml -p $(DC_PROJECT)
NOW=$(shell date '+%Y-%m-%d %H:%M:%S')

.PHONY: build
build:
	VERSION="dev" COMMIT="wip-commit" DATE="$(NOW)" ./scripts/build.sh

.PHONY: fmt
fmt:
	./scripts/fmt.sh

.PHONY: vet
vet:
	./scripts/vet.sh

.PHONY: lint-last-commit
lint-last-commit:
	npx commitlint --last

.PHONY: dc-up
dc-up:
	$(DC) up -d

.PHONY: dc-down
dc-down:
	$(DC) down --remove-orphans