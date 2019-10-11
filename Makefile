NAME=aery-graphql
VERSION=0.0.1

.PHONY: init
init:
	@go mod init $(NAME)

.PHONY: build
build:
	@go build -o build/$(NAME)

.PHONY: run
run: build
	@./build/$(NAME) -env development

.PHONY: run-prod
run-prod: build
	@./build/$(NAME) -env production

.PHONY: clean
clean:
	@rm -r build

.PHONY: generate
generate:
	@if [ -d tmp ]; then rm -r tmp; fi;
	@go run scripts/gqlgen.go
