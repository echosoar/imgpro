package method

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"

	img "github.com/echosoar/imgpro/core"
)

func OutputToImg(target string, width int, height int, pixels []img.RGBA) {
	bounds := image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: width,
			Y: height,
		},
	}
	imgItem := image.NewRGBA(bounds)
	for line := 0; line < height; line++ {
		for row := 0; row < width; row++ {
			pixel := pixels[line*width+row]
			multiPixel := color.RGBA{uint8(pixel.R), uint8(pixel.G), uint8(pixel.B), uint8(pixel.A)}
			imgItem.Set(row, line, multiPixel)
		}
	}

	imgItemFile, _ := os.Create(target)
	defer imgItemFile.Close()
	jpeg.Encode(imgItemFile, imgItem, nil)
}
