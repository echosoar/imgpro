package processor

import (
	img "github.com/echosoar/imgpro/core"
	method "github.com/echosoar/imgpro/method"
)

// HUEProcessor processor
func HUEProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:          []string{"hue"},
		PreConditions: []string{"rgba", "frame"},
		Runner:        hueRunner,
	})
}

const hueProcessSize = 200

func hueExec(rgbaList []img.RGBA) []img.RGBA {
	// var hueResult []img.RGBA
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

	// using canapy calc color size
	canopyInstance := method.Canopy{
		AllPoints: samplingRGBA,
		T1:        30,
		T2:        20,
	}
	canopyInstance.Run()
	canopyResult := canopyInstance.Result(2)

	kmeans := method.KMeans{
		K:         len(canopyResult),
		Center:    canopyResult,
		AllPoints: samplingRGBA,
	}
	kmeans.Run()
	return kmeans.GetResult()
}

func hueRunner(core *img.Core) map[string]img.Value {
	rgba := core.Result["rgba"].Frames

	frame := core.Result["frame"].Int
	var hueFrames []img.Value
	for frameIndex := 0; frameIndex < frame; frameIndex++ {
		hueFrames = append(hueFrames, img.Value{
			Type: img.ValueTypeRGBA,
			Rgba: hueExec(rgba[frameIndex].Rgba),
		})
	}

	return map[string]img.Value{
		"hue": {
			Type:   img.ValueTypeFrames,
			Frames: hueFrames,
		},
	}
}
