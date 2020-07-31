## Go-swagger server

### Go-swagger

The server side uses [go-swagger](https://github.com/go-swagger/go-swagger) to generate REST APIs. OpenAPI specifications can be found at [openapi](https://swagger.io/specification/v2/).

To generate REST API,

```
make gen
```

will check your OpenAPI specifications and then generate corresponding code.

### Gorm

The server uses [gorm](https://github.com/jinzhu/gorm) to map Go structs to database schemas and interact with [Postgresql](https://www.postgresql.org/). The doc can be found [here](https://gorm.io/docs/).

### Development

`glide` install will install Golang dependencies in the repo.

To start the server,

```
go run main.go
```

To add new features, add new API specifications in `swagger.yml` and implement API logics in `/pkg/handler`.