// main declares the CLI that spins up the server of
// our API.
// It takes some arguments, validates if they're valid
// and match the expected type and then intiialize the
// server.
package main

import (
	"flag"
	"strconv"
	"tides-server/pkg/config"
	"tides-server/pkg/controller"
	"tides-server/pkg/restapi/operations"

	"github.com/go-openapi/loads"

	"tides-server/pkg/restapi"

	"fmt"
	"log"
	"os"
)

func main() {
	conf := config.GetConfig()
	boolPtr := flag.Bool("local", false, "a bool")
	flag.Parse()
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewCloudTidesAPI(swaggerSpec)
	server := restapi.NewServer(api)
	server.ConfigureAPI()
	defer server.Shutdown()

	server.Host = "0.0.0.0"
	conf.Debug = *boolPtr
	server.Port, _ = strconv.Atoi(conf.Port)
	config.StartDB()

	name, err := os.Hostname()
	fmt.Println(name)

	controller.InitController()

	// Start listening using having the handlers and port
	// already set up.
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

}
