# Go-swagger server

## Go-swagger

The server side uses [go-swagger](https://github.com/go-swagger/go-swagger) to generate REST APIs. OpenAPI specifications can be found at [OpenAPI](https://swagger.io/specification/v2/).

To generate REST API,

```
make gen
```

will check your OpenAPI specifications and then generate corresponding code.

Generate API doc:
```
swagger serve ./swagger/swagger.yml
```

## Gorm

The server uses [gorm](https://github.com/go-gorm/gorm) to map Go structs to database schemas and interact with [PostgreSQL](https://www.postgresql.org/). The doc can be found [here](https://gorm.io/docs/).

## Development

Run following commands to install Golang dependencies:
```
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.io,direct
export GO111MODULE=on
go get -v all
```

Run `dev` server:

```
go run ./cmd/main.go
```

To run `test` server in local env, set env variables listed in `.env` and run
```
go run ./cmd/main.go -local
```

To add new features, add new API specifications in `./swagger/swagger.yml` and implement API logics in `./pkg/handler`.