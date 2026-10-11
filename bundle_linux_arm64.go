//go:build bundled

package main

import "embed"

//go:embed deps/linux_arm64
var bundledFS embed.FS

const bundledDir = "deps/linux_arm64"
