package templates

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
)

//go:embed *.html layouts/*.html
var embedded embed.FS

var temps = template.Must(template.New("templates").Funcs(template.FuncMap{
	"args": func(pairs ...any) (map[string]any, error) {
		if len(pairs)%2 != 0 {
			return nil, errors.New("misaligned map in template arguments")
		}

		m := make(map[string]any, len(pairs)/2)

		for i := 0; i < len(pairs); i += 2 {
			key, ok := pairs[i].(string)
			if !ok {
				return nil, fmt.Errorf("cannot use type %T as map key", pairs[i])
			}

			m[key] = pairs[i+1]
		}

		return m, nil
	},
}).ParseFS(embedded,
	"*.html",
	"layouts/*.html",
))

func Templates() *template.Template {
	return temps
}
