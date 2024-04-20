package views

import "embed"

var (
	//go:embed *.html
	Views embed.FS
)
