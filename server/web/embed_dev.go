//go:build !embed

package web

import "embed"

var adminFS embed.FS
var userFS embed.FS

const embedEnabled = false
