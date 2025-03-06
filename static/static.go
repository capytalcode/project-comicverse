package static

import (
	"embed"
	"io/fs"
)

//go:generate tailwindcss -o static/css/wind.css

//go:embed css/*.css
var staticFiles embed.FS

func Files(local ...bool) fs.FS {
	var l bool
	if len(local) > 0 {
		l = local[0]
	}

	if !l {
		return staticFiles
	}

	return staticFiles
}
