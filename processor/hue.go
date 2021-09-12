package processor

import (
	img "github.com/echosoar/imgpro/core"
)

// HUEProcessor processor
func HUEProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:          []string{"hue"},
		PreConditions: []string{"rgba", "frame"},
		Runner:        hueRunner,
	})
}

const hueProcessSize = 20

func hueExec(rgbaList []img.RGBA) []img.RGBA {
	var hueResult []img.RGBA
	var samplingStep int = 1
	var samplingRGBA []img.RGBA

	size := len(rgbaList)
	if size > hueProcessSize {
		// int auto floor
		samplingStep = size / hueProcessSize
	}
	for index := 0; index < size; index += samplingStep {
		samplingRGBA = append(samplingRGBA, rgbaList[index])
	}
	return hueResult
}

func hueRunner(core *img.Core) map[string]img.Value {
	rgba := core.Result["rgba"].Rgba
	frame := core.Result["frame"].Int
	var hueResult [][]img.RGBA
	for frameIndex := 0; frameIndex < frame; frameIndex++ {
		hueResult = append(hueResult, hueExec(rgba[frameIndex]))
	}

	return map[string]img.Value{
		"hue": {
			Type: img.ValueTypeRGBA,
			Rgba: hueResult,
		},
	}
}
