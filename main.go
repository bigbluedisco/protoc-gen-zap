package main

import (
	"github.com/bigbluedisco/protoc-gen-zap/module"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func main() {
	pgs.
		Init(pgs.DebugEnv("DEBUG_PGV")).
		RegisterModule(module.Zap()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
