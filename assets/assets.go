package assets

import (
	"embed"
	"io/fs"
)

//go:generate tailwindcss -o assets/stylesheets/out.css

//go:embed stylesheets/out.css
var files embed.FS

func Files(local ...bool) fs.FS {
	var l bool
	if len(local) > 0 {
		l = local[0]
	}

	if !l {
		return files
	}

	return files
}
