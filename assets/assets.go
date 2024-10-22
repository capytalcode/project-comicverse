package assets

import (
	_ "embed"
)

//go:embed css/uno.css
var UNO_CSS []byte

//go:embed css/theme.css
var THEME_CSS []byte
