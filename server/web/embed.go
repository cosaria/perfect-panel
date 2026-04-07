//go:build embed

package web

import "embed"

//go:embed all:admin-dist
var adminFS embed.FS

//go:embed all:user-dist
var userFS embed.FS

const embedEnabled = true
