NAME=go-graphql-mongodb-boilerplate
VERSION=1.0.0

# OUTSIDE DOCKER
################
start: #generate #generate needed because of the local linter and git
	docker-compose run --service-ports $(NAME)

restart:
	docker container prune -f
	docker-compose down --volumes --rmi all
	docker-compose run --service-ports $(NAME)

stop:
	docker container prune -f
	docker-compose down --volumes --rmi all


# INSIDE DOCKER
################
.PHONY: init
init:
	@go mod init $(NAME)

.PHONY: build
build:
	@go build -o build/$(NAME)

.PHONY: run
run: build
	@./build/$(NAME)

.PHONY: clean
clean:
	@rm -r build

.PHONY: generate
generate:
	@if [ -d tmp ]; then rm -r tmp; fi;
	@go run hack/gqlgen-generator.go
