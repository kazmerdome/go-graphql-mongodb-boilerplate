FROM golang:alpine as dev
WORKDIR /go-graphql-mongodb-boilerplate/
RUN apk add --update make
EXPOSE 9090 9229 9230
COPY . /go-graphql-mongodb-boilerplate/
ENV CGO_ENABLED 0
RUN make generate
RUN make build

FROM alpine:latest as prod
RUN apk --no-cache add ca-certificates
WORKDIR /run/
COPY --from=dev /go-graphql-mongodb-boilerplate/build/go-graphql-mongodb-boilerplate .
CMD ["./go-graphql-mongodb-boilerplate"]
