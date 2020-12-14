# Server-Side Development

## Start server alone in local mode
Set up requires [Golang](https://golang.org/).

Run following commands to install Golang dependencies:
```
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.io,direct
export GO111MODULE=on
go get -v all
```

Change the server ip to `localhost` in `cmd/main.go`. Change database configuration in `pkg/config/type.go`
```
cd tides-server/cmd
go run main.go
```

Future improvements:
- The way of starting local mode is not convenient. Better to read from a config file.

## Development

Add new REST API:
* Add new API specifications in `./swagger/swagger.yml`
* `make gen` to generate corresponding code
* Implement API logics in `./pkg/handler`
* In `./pkg/restapi/configure_cloud_tides.go`, configure server API with your handler implementation in `configureAPI` function.

