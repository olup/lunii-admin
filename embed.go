//go:build !bindings

package main

import "embed"

//go:embed all:frontend/dist
var assets embed.FS
