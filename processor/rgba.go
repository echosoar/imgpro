package processor

import (
	"bufio"
	"image"
	"os"

	img "github.com/echosoar/imgpro/core"
)

// RGBAProcessor bin size processor
func RGBAProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:         []string{"rgba", "frame"},
		Precondition: []string{"type", "wh"},
		Runner:       rgbaRunner,
	})
}

func uint32ToInt(num uint32) int {
	return int(num >> 8)
}

func rgbaRunner(core *img.Core) map[string]img.Value {
	imgType := core.Result["type"].String
	width := core.Result["width"].Int
	height := core.Result["height"].Int
	frame := 1
	rgba := [][]img.RGBA{}
	f, err := os.Open(core.FilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)

	originalImage, _, err := image.Decode(reader)
	if imgType == "png" || imgType == "jpg" {
		rgbaFrame := []img.RGBA{}
		for line := 0; line < height; line++ {
			for col := 0; col < width; col++ {
				r, g, b, a := originalImage.At(col, line).RGBA()
				rgbaFrame = append(rgbaFrame, img.RGBA{
					R: uint32ToInt(r),
					G: uint32ToInt(g),
					B: uint32ToInt(b),
					A: uint32ToInt(a),
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
