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

To run `dev` server, create a `.env` file in the root of `tides-server` directory with following configurations:
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

## API Implementation

Add new REST API:
* Add new API specifications in `./swagger/swagger.yml`
* `make gen` to generate corresponding code
* Implement API logics in `./pkg/handler`
* In `./pkg/restapi/configure_cloud_tides.go`, configure server API with your handler implementation in `configureAPI` function.

## API Testing

[Postman](https://www.postman.com/) is commonly used for API testing. It simulates client-server communication in a web app.

## Future Improvements

- Tests should be added. More sophosticated tests should be added in CI workflow.
- Implement cost estimator to inform users about their contribution.
- Deleting a resource should destroy all the VMs deployed on the resource. Not implemented now.
