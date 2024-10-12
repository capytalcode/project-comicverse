package main

import (
	"embed"
	"flag"
	"net/http"
	"os"
	"strconv"

	"forge.capytal.company/capytalcode/project-comicverse/app"
)

//go:embed assets
var assetsFolder embed.FS

var (
	port   *int
	dev    *bool
	assets *string
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

	assetsEnv := os.Getenv("COMICVERSE_ASSETS")
	if assetsEnv == "" {
		assetsEnv = "./assets"
	}
	assets = flag.String("assets", assetsEnv, "The directory for the development assets")
}

func main() {
	flag.Parse()

	var assetsFS http.Handler
	if *dev {
		assetsFS = http.StripPrefix("/assets/", http.FileServer(http.Dir(*assets)))
	} else {
		assetsFS = http.FileServerFS(assetsFolder)
	}

	app := app.NewApp(app.AppOpts{
		Port:   port,
		Dev:    dev,
		Assets: assetsFS,
	})

	app.Run()
}
