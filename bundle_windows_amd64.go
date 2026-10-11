//go:build bundled

package main

import "embed"

//go:embed deps/windows_amd64
var bundledFS embed.FS

const bundledDir = "deps/windows_amd64"
