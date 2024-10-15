package pages

import (
	"forge.capytal.company/capytalcode/project-comicverse/router"
)

var PAGES = []router.Route{
	{Pattern: "/_dev/colors", Handler: &Colors{}},
}
