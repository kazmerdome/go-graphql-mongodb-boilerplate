NAME=aery-graphql
VERSION=0.0.1

# OUTSIDE DOCKER
start:
	docker-compose run --service-ports $(NAME)

restart:
	docker container prune -f
	docker-compose down --volumes --rmi all
	docker-compose run --service-ports $(NAME)

stop:
	docker container prune -f
	docker-compose down --volumes --rmi all

# INSIDE DOCKER
.PHONY: init
init:
	@go mod init $(NAME)

.PHONY: build
build:
	@go build -o build/$(NAME)

.PHONY: run
run: build
	@./build/$(NAME) -env development

.PHONY: clean
clean:
	@rm -r build

.PHONY: generate
generate:
	@if [ -d tmp ]; then rm -r tmp; fi;
	@go run scripts/gqlgen.go
