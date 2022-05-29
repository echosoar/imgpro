package method

import (
	img "github.com/echosoar/imgpro/core"
)

func RGBAToGrey(rgba img.RGBA) uint8 {
	// ref: https://www.dcode.fr/grayscale-image
	Y := 0.2125*float64(rgba.R) + 0.7154*float64(rgba.G) + 0.0721*float64(rgba.B)
	return uint8(Y)
}

func IterateRGBA(point img.RGBA) []int {
	return []int{point.R, point.G, point.B, point.A}
}

func ColorListToRGBA(colorList []int) img.RGBA {
	return img.RGBA{
		R: colorList[0],
		G: colorList[1],
		B: colorList[2],
		A: colorList[3],
	}
}

func IsSameColor(pointA img.RGBA, pointB img.RGBA) bool {
	return pointA.R == pointB.R && pointA.G == pointB.G && pointA.B == pointB.B && pointA.A == pointB.A
}
