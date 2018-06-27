package main

import (
	"github.com/bigbluedisco/protoc-gen-zap/module"
	"github.com/lyft/protoc-gen-star"
)

func main() {
	pgs.
		Init(pgs.DebugEnv("DEBUG_PGV")).
		RegisterModule(module.Zap()).
		RegisterPostProcessor(pgs.GoFmt()).
		Render()
}
