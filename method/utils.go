package method

import (
	"math"

	img "github.com/echosoar/imgpro/core"
)

func ManhattanDistance(pointA img.RGBA, pointB img.RGBA) float64 {
	var distance float64 = 0.0
	distance += math.Pow(float64(pointA.R-pointB.R), 2)
	distance += math.Pow(float64(pointA.G-pointB.G), 2)
	distance += math.Pow(float64(pointA.B-pointB.B), 2)
	distance += math.Pow(float64(pointA.A-pointB.A), 2)
	distance = math.Sqrt(distance)
	return distance
}

func AverageColor(list []img.RGBA) img.RGBA {

	listLen := len(list)
	if listLen == 0 {
		return img.RGBA{}
	}
	sumR := 0
	sumG := 0
	sumB := 0
	sumA := 0

	for i := 0; i < listLen; i++ {
		sumR += list[i].R
		sumG += list[i].G
		sumB += list[i].B
		sumA += list[i].A
	}

	sumR = sumR / listLen
	sumG = sumG / listLen
	sumB = sumB / listLen
	sumA = sumA / listLen
	return img.RGBA{
		R: sumR,
		G: sumG,
		B: sumB,
		A: sumA,
	}
}
