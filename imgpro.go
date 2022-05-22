package imgpro

import (
	img "github.com/echosoar/imgpro/core"
	pro "github.com/echosoar/imgpro/processor"
)

func initial(features []string) *img.Core {
	core := img.New(features)
	// size
	pro.SizeProcessor(core)
	// type
	pro.TypeProcessor(core)
	// width/height
	pro.WHProcessor(core)
	// rgba/frame
	pro.RGBAProcessor(core)
	// hue
	pro.HUEProcessor(core)
	// exif
	pro.ExifProcessor(core)
	// time
	pro.TimeProcessor(core)
	// device
	pro.DeviceProcessor(core)
	// position
	pro.PositionProcessor(core)
	return core
}

// Run run
func Run(filePath string, features []string) img.Result {
	core := initial(features)
	core.Run(filePath)
	return core.GetResult()
}

// RunBinary run
func RunBinary(imgFileBinary []byte, features []string) img.Result {
	core := initial(features)
	core.RunBinary(imgFileBinary)
	return core.GetResult()
}
