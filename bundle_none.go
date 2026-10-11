//go:build !bundled

package main

import "embed"

// Builds without the bundled tag embed nothing and use tools from PATH.
var bundledFS embed.FS

const bundledDir = ""
