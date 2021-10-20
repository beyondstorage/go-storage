//go:build tools
// +build tools

package definitions

import "embed"

//go:embed *.toml
var Bindata embed.FS
