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
	frame := 0
	frameRGBAs := []img.Value{}

	if imgType == "png" || imgType == "jpg" {
		frame = 1
		originalImage, _, err := image.Decode(bytes.NewReader(core.FileBinary))
		if err != nil {
			panic(err)
		}
		rgbaFrame := make([]img.RGBA, height*width)
		for line := 0; line < height; line++ {
			for col := 0; col < width; col++ {
				index := line*width + col
				r, g, b, a := originalImage.At(col, line).RGBA()
				rgbaFrame[index] = img.RGBA{
					R: utils.Uint32ToInt(r),
					G: utils.Uint32ToInt(g),
					B: utils.Uint32ToInt(b),
					A: utils.Uint32ToInt(a),
				}
			}
		}
		frameRGBAs = append(frameRGBAs, img.Value{
			Type: img.ValueTypeRGBA,
			Rgba: rgbaFrame,
		})
	} else if imgType == "gif" {

	}

	return map[string]img.Value{
		"frame": {
			Type: img.ValueTypeInt,
			Int:  frame,
		},
		"rgba": {
			Type:   img.ValueTypeFrames,
			Frames: frameRGBAs,
		},
	}
}
