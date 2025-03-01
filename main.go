package capytalcodecomicverse

import (
	"flag"
)

var (
	debug = flag.Bool("dev", false, "Run the server in debug mode.")
	port  = flag.Int("port", 8080, "Port to be used for the server.")
)

func init() {
	flag.Parse()
}

func main() {

}
