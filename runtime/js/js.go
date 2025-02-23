package js

import (
	_ "embed"
)

//go:embed js.wasm
var Runtime []byte
