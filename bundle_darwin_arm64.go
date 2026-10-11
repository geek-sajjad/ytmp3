//go:build bundled

package main

import "embed"

//go:embed deps/darwin_arm64
var bundledFS embed.FS

const bundledDir = "deps/darwin_arm64"
