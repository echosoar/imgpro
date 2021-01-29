package imgpro

import (
	img "github.com/echosoar/imgpro/core"
	pro "github.com/echosoar/imgpro/processor"
)

// Run run
func Run(filePath string, features []string) img.Result {
	core := img.New(features)
	// size
	pro.SizeProcessor(core)
	// type
	pro.TypeProcessor(core)
	// width/height
	pro.WHProcessor(core)
	core.Run(filePath)
	return core.GetResult()
}
