
# go-graphql-mongodb-boilerplate

Lightweight, easy-to-develop graphql server that includes all you need to build amazing projects 🔥


- 🔮 **Gqlgen** — Generated, type safe Graphql for Go
- 👽 **Mongo Driver** - The official Mongodb driver for Go
- 🐶 **Dataloaden** — Generated type safe data loaders for Go
- 📄 **Echo** - High performance, extensible, minimalist Go web framework

## Features:

Server:
- Using Labstack Echo
- Restapi example
- prometheus implementation
- healthz - basic (demo only) health probe implementation

Graphql:
- using the latest (0.13.0) Gqlgen version
- playground security with http header password (Disable Introspection)
- custom scalar example
- dataloader examples (https://github.com/vektah/dataloaden) for n+1 problems

MongoDB
- model examples
- multiple db connections implementation

Other
- Config package for reading and caching ENV-s from global env, kubernetes, or docker swarm
- Using Makefile, Docker-compose for faster development


## 🚀 Getting started


### Development Mode
step1 - build local environment with docker-compose
```console
  make start
```

step2 - start the server (inside docker)
```console
  make run
```

step3 - restart the server (inside docker)
```console
  CTRL + C
  make run
```

regenerate gqlgen files
```console
  CTRL + C
  make generate
  make run
```

open `http://localhost:9090`.


### Enable Graphql Documentation:
 ```add "Playground-Password": <GRAPHQL_PLAYGROUND_PASS> to request header ```

### Working with Auth
```add "Authorization": Bearer <JWT TOKEN> to request header ```

### Create the first user
  - step 1 - create new user
  ```graphql
    mutation {
      authBasicStrategySignUp(data: {
        email: "x@x.x"
        password: "xxxxx"
        username: "x"
      }) {
        token
      }
    }
  ```
  - step 2 - check created user
    ```graphql
    query {
      userMe {
        email
        username
        role
      }
    }
    ```
  - step 2 - add ADMIN role to the generated user (through db admin [mongo atlas, robo3t, etc...])

### Create New CRUD - TODO
  - step 1 - create the model
  - TODO



## Production Mode
step0 - make sure to provide envs (copy .env file to build folder or provide global envs)

step1 - binary build
```console
  make generate
  make build
```

step2 - run
```console
  ./build/go-graphql-mongodb-boilerplate
```
