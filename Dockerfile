FROM golang:alpine as dev
WORKDIR /aery-graphql/
RUN apk add --update make
EXPOSE 9090
COPY . /aery-graphql/
ENV CGO_ENABLED 0
RUN go install
RUN make generate
RUN make build

FROM alpine:latest as prod
RUN apk --no-cache add ca-certificates
WORKDIR /run/
COPY --from=dev /aery-graphql/build/aery-graphql .
CMD ["./aery-graphql"]