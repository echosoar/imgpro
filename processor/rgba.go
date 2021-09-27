package processor

import (
	"bytes"
	"image"

	img "github.com/echosoar/imgpro/core"
	utils "github.com/echosoar/imgpro/utils"
)

// RGBAProcessor bin size processor
func RGBAProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:          []string{"rgba", "frame"},
		PreConditions: []string{"type", "wh"},
		Runner:        rgbaRunner,
	})
}

func rgbaRunner(core *img.Core) map[string]img.Value {
	imgType := core.Result["type"].String
	width := core.Result["width"].Int
	height := core.Result["height"].Int
	frame := 1
	rgba := [][]img.RGBA{}

	originalImage, _, err := image.Decode(bytes.NewReader(core.FileBinary))
	if err != nil {
		panic(err)
	}
	if imgType == "png" || imgType == "jpg" {
		rgbaFrame := []img.RGBA{}
		for line := 0; line < height; line++ {
			for col := 0; col < width; col++ {
				r, g, b, a := originalImage.At(col, line).RGBA()
				rgbaFrame = append(rgbaFrame, img.RGBA{
					R: utils.Uint32ToInt(r),
					G: utils.Uint32ToInt(g),
					B: utils.Uint32ToInt(b),
					A: utils.Uint32ToInt(a),
				})
			}
		}
		rgba = append(rgba, rgbaFrame)
	} else if imgType == "gif" {

	}

	return map[string]img.Value{
		"frame": {
			Type: img.ValueTypeInt,
			Int:  frame,
		},
		"rgba": {
			Type: img.ValueTypeRGBA,
			Rgba: rgba,
		},
	}
}
