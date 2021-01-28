package imgpro

import (
	img "github.com/echosoar/imgpro/core"
	pro "github.com/echosoar/imgpro/processor"
)

// Run run
func Run(filePath string, features []string) img.Result {
	core := img.New()
	core.Features = features
	pro.BindSize(core)
	return core.Run(filePath)
}
