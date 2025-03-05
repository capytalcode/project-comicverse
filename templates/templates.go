package templates

import (
	"embed"
	"html/template"
)

//go:embed *.html test/*.html
var embedded embed.FS

var temps = template.Must(template.ParseFS(embedded,
	"*.html",
	"test/*.html",
))

func Templates() *template.Template {
	return temps
}
