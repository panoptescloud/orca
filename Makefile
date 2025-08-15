.PHONY: build
build:
	./scripts/build.sh

.PHONY: lint-last-commit
lint-last-commit:
	npx commitlint --last