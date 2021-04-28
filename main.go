package main

import (
	"fmt"
	"log"
	"path"

	"github.com/fittipaldi/go-r3-d3/app"
	"github.com/fittipaldi/go-r3-d3/config"

	"github.com/joho/godotenv"
)

func main() {

	projectRootDir := "./"
	err := godotenv.Load(path.Join(projectRootDir, ".env"))
	if err != nil {
		log.Fatal(fmt.Println("Missing the [.env] config file"))
	}

	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)

	app.Run(config.Web.HttpPorts)
}
