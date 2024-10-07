package main

import (
	"flag"
	"os"
	"strconv"

	comicverse "forge.capytal.company/capytalcode/project-comicverse"
)

var port *int

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
}

func main() {
	flag.Parse()
	comicverse.Run(*port)
}
