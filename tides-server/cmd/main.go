// main declares the CLI that spins up the server of
// our API.
// It takes some arguments, validates if they're valid
// and match the expected type and then intiialize the
// server.
package main

import (
	"tides-server/pkg/config"
	"tides-server/pkg/controller"
	"tides-server/pkg/restapi/operations"

	"github.com/go-openapi/loads"

	"tides-server/pkg/restapi"

	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	config := config.GetConfig()
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewCloudTidesAPI(swaggerSpec)
	server := restapi.NewServer(api)
	server.ConfigureAPI()
	defer server.Shutdown()

	server.Port, _ = strconv.Atoi(config.Port)
	fmt.Println(server.Port)

	name, err := os.Hostname()
	fmt.Println(name)

	server.Host = "localhost"
	controller.InitController()

	/*
		// Implement the handler functionality.
		// As all we need to do is give an implementation to the interface
		// we can just override the `api` method giving it a method with a valid
		// signature (we didn't need to have this implementation here, it could
		// even come from a different package).
		api.GetHostnameHandler = operations.GetHostnameHandlerFunc(
			func(params operations.GetHostnameParams) middleware.Responder {
				response, err := os.Hostname()
				if err != nil {
					return operations.NewGetHostnameDefault(500).WithPayload(&models.Error{
						Code: 500,
						Message: swag.String("failed to retrieve hostname"),
					})
				}

				return operations.NewGetHostnameOK().WithPayload(response)
			})
	*/
	// Start listening using having the handlers and port
	// already set up.
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

}
