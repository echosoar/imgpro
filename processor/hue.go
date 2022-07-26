package processor

import (
	img "github.com/echosoar/imgpro/core"
	method "github.com/echosoar/imgpro/method"
)

// HUEProcessor processor
func HUEProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:          []string{"hue"},
		PreConditions: []string{"rgba", "frame", "width", "height"},
		Runner:        hueRunner,
	})
}

func hueExec(width int, height int, rgbaList []img.RGBA) []img.RGBA {

	minSize := width
	if height < width {
		minSize = height
	}

	// var hueResult []img.RGBA
	var samplingStep int = minSize / 100
	var samplingRGBA []img.RGBA

	if samplingStep < 1 {
		samplingStep = 1
	}
	for hIndex := 0; hIndex < height; hIndex += samplingStep {
		for wIndex := 0; wIndex < width; wIndex += samplingStep {
			index := hIndex*width + wIndex
			samplingRGBA = append(samplingRGBA, rgbaList[index])
		}
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
		Center:    make([]img.RGBA, 0),
		AllPoints: samplingRGBA,
	}
	kmeans.Run()
	return kmeans.GetResult()
}

func hueRunner(core *img.Core) map[string]img.Value {
	rgba := core.Result["rgba"].Frames
	width := core.Result["width"].Int
	height := core.Result["height"].Int

	frame := core.Result["frame"].Int
	var hueFrames []img.Value
	for frameIndex := 0; frameIndex < frame; frameIndex++ {
		hueFrames = append(hueFrames, img.Value{
			Type: img.ValueTypeRGBA,
			Rgba: hueExec(width, height, rgba[frameIndex].Rgba),
		})
	}

	return map[string]img.Value{
		"hue": {
			Type:   img.ValueTypeFrames,
			Frames: hueFrames,
		},
	}
}
