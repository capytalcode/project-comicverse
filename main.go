package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"

	"forge.capytal.company/capytalcode/project-comicverse/app"
)

var (
	port *int
	dev  *bool
)

func init() {
	portEnv := os.Getenv("COMICVERSE_PORT")
	if portEnv == "" {
		portEnv = "8080"
	}
	p, err := strconv.Atoi(portEnv)
	if err != nil {
		p = 8080
	}
	port = flag.Int("port", p, "The port to the server to listen on")

	devEnv := os.Getenv("COMICVERSE_DEV")
	if devEnv == "" {
		devEnv = "false"
	}
	d, err := strconv.ParseBool(devEnv)
	if err != nil {
		d = false
	}
	dev = flag.Bool("dev", d, "Run the application in development mode")
}

func main() {
	flag.Parse()

	app := app.NewApp(app.AppOpts{
		Port:   port,
		Dev:    dev,
		Assets: http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))),
	})

	app.Run()
}
