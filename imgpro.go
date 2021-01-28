package imgpro

import (
	img "github.com/echosoar/imgpro/core"
	pro "github.com/echosoar/imgpro/processor"
)

// Run run
func Run(filePath string, features []string) img.Result {
	core := img.New(features)
	pro.SizeProcessor(core)
	pro.TypeProcessor(core)
	core.Run(filePath)
	return core.GetResult()
}
