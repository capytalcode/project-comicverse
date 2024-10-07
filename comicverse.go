package comicverse

import (
	"fmt"
	"log"
	"net/http"

	"forge.capytal.company/capytalcode/project-comicverse/pages"
	"forge.capytal.company/capytalcode/project-comicverse/router"
)

func Run(port int) {
	router := router.NewRouter()

	router.HandleRoutes(pages.PAGES)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), router); err != nil {
		log.Fatal(err)
	}
}
