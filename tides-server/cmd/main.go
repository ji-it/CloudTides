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

	"fmt"
	"log"
	"os"
	"tides-server/pkg/restapi"

	"github.com/joho/godotenv"
)

func main() {
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

	if !*boolPtr {
		err := godotenv.Load() // load .env file
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	config.InitConfig()
	conf := config.GetConfig()
	conf.Debug = *boolPtr
	server.Host = os.Getenv("envkey_SERVER_IP")
	server.Port, err = strconv.Atoi(os.Getenv("envkey_SERVER_PORT"))

	name, err := os.Hostname()
	fmt.Println(name)

	controller.InitController()

	// Start listening using having the handlers and port
	// already set up.
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

}
