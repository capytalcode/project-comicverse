package capytalcodecomicverse

import (
	"io/fs"
	"net/http"
)

type App struct {
	templates fs.FS
}

func NewApp(templates fs.FS) *App {
	return &App{
		templates: templates,
	}
}
