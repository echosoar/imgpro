package processor

import (
	"bytes"
	"image"
	"image/draw"
	"image/gif"

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

func readImageRGBA(width, height int, bounds *image.Rectangle, rgbaFrameRef *[]img.RGBA, rgba *image.RGBA) {
	colorIndex := 0
	for line := 0; line < height; line++ {
		for col := 0; col < width; col++ {
			itemIndex := line*width + col
			if line < bounds.Min.Y || line >= bounds.Max.Y || col < bounds.Min.X || col >= bounds.Max.X {
				(*rgbaFrameRef)[itemIndex] = img.RGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 0,
				}
			} else {
				index := colorIndex * 4
				colorIndex++
				(*rgbaFrameRef)[itemIndex] = img.RGBA{
					R: int(rgba.Pix[index]),
					G: int(rgba.Pix[index+1]),
					B: int(rgba.Pix[index+2]),
					A: int(rgba.Pix[index+3]),
				}
			}

		}
	}
}

func rgbaRunner(core *img.Core) map[string]img.Value {
	imgType := core.Result["type"].String
	width := core.Result["width"].Int
	height := core.Result["height"].Int
	frame := 0
	frameRGBAs := []img.Value{}

	if imgType == "png" || imgType == "jpg" {
		frame = 1
		frameRGBAs = make([]img.Value, 1)
		originalImage, _, err := image.Decode(bytes.NewReader(core.FileBinary))
		if err != nil {
			panic(err)
		}
		rgbaFrame := make([]img.RGBA, height*width)
		bounds := originalImage.Bounds()
		rgba := image.NewRGBA(bounds)
		draw.Draw(rgba, bounds, originalImage, bounds.Min, draw.Src)
		readImageRGBA(width, height, &bounds, &rgbaFrame, rgba)
		frameRGBAs[0] = img.Value{
			Type: img.ValueTypeRGBA,
			Rgba: rgbaFrame,
		}
	} else if imgType == "gif" {
		gifInstance, err := gif.DecodeAll(bytes.NewReader(core.FileBinary))
		if err != nil {
			panic(err)
		}
		frame = len(gifInstance.Image)
		frameRGBAs = make([]img.Value, frame)
		for frameIndex, imageInstance := range gifInstance.Image {
			rgbaFrame := make([]img.RGBA, height*width)
			bounds := imageInstance.Bounds()
			rgba := image.NewRGBA(bounds)
			draw.Draw(rgba, bounds, imageInstance, bounds.Min, draw.Src)
			readImageRGBA(width, height, &bounds, &rgbaFrame, rgba)
			frameRGBAs[frameIndex] = img.Value{
				Type: img.ValueTypeRGBA,
				Rgba: rgbaFrame,
			}
		}
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
