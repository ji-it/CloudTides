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

To run `dev` server, create a `.env` file in the root of this directory with following configurations:
```
SERVER_IP=
SERVER_PORT=
POSTGRES_HOST=
POSTGRES_PORT=
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
ADMIN_USER=
ADMIN_PASSWORD=
SECRET_KEY=
```

If running in local environment, `SERVER_IP` and `POSTGRES_HOST` would be `localhost`. Then start the server:
```
go run ./cmd/main.go
```

To add new features, add new API specifications in `./swagger/swagger.yml` and implement API logics in `./pkg/handler`. For more details, refer to [developer guide](https://cloudtides.github.io/CloudTides/#/).