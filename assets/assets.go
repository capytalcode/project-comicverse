package assets

import (
	"embed"
)

//go:embed css fonts javascript
var ASSETS embed.FS
