package processor

import (
	img "github.com/echosoar/imgpro/core"
)

// FrameProcessor get image frame count
func FrameProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:         []string{"frame"},
		Precondition: []string{"type"},
		Runner:       frameRunner,
	})
}

func frameRunner(core *img.Core) map[string]img.Value {
	frame := 0
	imgType := core.Result["type"].String
	if imgType == "jpg" || imgType == "png" || imgType == "bmp" {
		frame = 1
	} else if imgType == "bmp" {

	} else if imgType == "webp" {

	}
	return map[string]img.Value{
		"frame": {
			Type: img.ValueTypeInt,
			Int:  frame,
		},
	}
}
