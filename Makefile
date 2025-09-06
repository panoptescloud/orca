DC_PROJECT=pc_orca
DC=docker compose -f ./docker-compose.yml -p $(DC_PROJECT)
NOW=$(shell date '+%Y-%m-%d %H:%M:%S')
TESTS?=./...
.PHONY: build fmt vet lint-last-commit dc-up dc-down test gen-mocks


gen-mocks:
	@mockery

build: fmt
	@VERSION="dev" COMMIT="wip-commit" DATE="$(NOW)" ./scripts/build.sh

build-global: fmt
	@OUTDIR=/usr/local/bin/orca VERSION="dev" COMMIT="wip-commit" DATE="$(NOW)" ./scripts/build.sh

fmt:
	@./scripts/fmt.sh

test:
	@./scripts/test.sh $(TESTS)

vet:
	@./scripts/vet.sh

lint-last-commit:
	@npx commitlint --last

dc-up:
	$(DC) up -d

dc-down:
	$(DC) down --remove-orphans