package processor

import (
	"bytes"
	"image"
	"image/draw"

	img "github.com/echosoar/imgpro/core"
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
		bounds := originalImage.Bounds()
		rgba := image.NewRGBA(originalImage.Bounds())
		draw.Draw(rgba, bounds, originalImage, bounds.Min, draw.Src)
		rgbaFrame := make([]img.RGBA, height*width)
		for line := 0; line < height; line++ {
			for col := 0; col < width; col++ {
				itemIndex := line*width + col
				index := itemIndex * 4
				rgbaFrame[itemIndex] = img.RGBA{
					R: int(rgba.Pix[index]),
					G: int(rgba.Pix[index+1]),
					B: int(rgba.Pix[index+2]),
					A: int(rgba.Pix[index+3]),
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
