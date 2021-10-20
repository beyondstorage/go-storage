package definitions

import "embed"

//go:embed *.toml
var Bindata embed.FS
