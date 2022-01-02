package processor

import (
	"math"
	"math/rand"

	img "github.com/echosoar/imgpro/core"
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
	canopyInstance := canopy{
		AllPoints: samplingRGBA,
		T1:        30,
		T2:        20,
	}
	canopyInstance.run()
	canopyResult := canopyInstance.result(2)
	return canopyResult
}

func hueRunner(core *img.Core) map[string]img.Value {
	rgba := core.Result["rgba"].Rgba

	frame := core.Result["frame"].Int
	var hueResult [][]img.RGBA
	for frameIndex := 0; frameIndex < frame; frameIndex++ {
		hueResult = append(hueResult, hueExec(rgba[frameIndex]))
	}

	return map[string]img.Value{
		"hue": {
			Type: img.ValueTypeRGBA,
			Rgba: hueResult,
		},
	}
}

type canopyPoint struct {
	Center img.RGBA
	Points []img.RGBA
}

func (c *canopyPoint) Add(point img.RGBA) {
	c.Points = append(c.Points, point)
}

type canopy struct {
	AllPoints []img.RGBA
	Canopies  []canopyPoint
	T1        float64
	T2        float64
}

func (c *canopy) getT1T2() {
	center := c.average(c.AllPoints)
	var allDistance float64 = 0
	for i := 0; i < len(c.AllPoints); i++ {
		current := c.AllPoints[i]
		distance := math.Pow((float64(center.R-current.R)), 2) + math.Pow((float64(center.G-current.G)), 2) + math.Pow((float64(center.B-current.B)), 2) + math.Pow((float64(center.A-current.A)), 2)

		allDistance += math.Sqrt(distance)
	}
	c.T1 = allDistance / float64(len(c.AllPoints))
	c.T2 = c.T1 * 2 / 3
}

func (c *canopy) run() {
	c.getT1T2()
	for len(c.AllPoints) != 0 {
		var newAllPoint []img.RGBA
		c.getRandom(c.AllPoints)
		for i := 0; i < len(c.AllPoints); i++ {
			current := c.AllPoints[i]
			isRemove := false
			index := 0
			for canopiesIndex := 0; canopiesIndex < len(c.Canopies); canopiesIndex++ {
				canopyPoint := c.Canopies[canopiesIndex]
				canopyCenter := canopyPoint.Center
				distance := c.manhattanDistance(canopyCenter, current)

				if distance <= c.T1 {
					canopyPoint.Add(current)
				} else {
					index++
				}

				if distance <= c.T2 {
					isRemove = true
				}
			}

			if index == len(c.Canopies) {
				c.Canopies = append(c.Canopies, canopyPoint{
					Center: current,
					Points: []img.RGBA{current},
				})
				continue
			}

			if !isRemove {
				newAllPoint = append(newAllPoint, current)
			}
		}
		c.AllPoints = newAllPoint
	}
}

func (c *canopy) getRandom(allPoints []img.RGBA) {
	index := rand.Intn(len(allPoints))
	current := allPoints[index]

	c.AllPoints = append(allPoints[:index], allPoints[index+1:]...)
	c.Canopies = append(c.Canopies, canopyPoint{
		Center: current,
		Points: []img.RGBA{current},
	})
}

func (c *canopy) manhattanDistance(pointA img.RGBA, pointB img.RGBA) float64 {
	var distance float64 = 0.0
	distance += math.Pow(float64(pointA.R-pointB.R), 2)
	distance += math.Pow(float64(pointA.G-pointB.G), 2)
	distance += math.Pow(float64(pointA.B-pointB.B), 2)
	distance += math.Pow(float64(pointA.A-pointB.A), 2)
	distance = math.Sqrt(distance)
	return distance
}

func (c *canopy) average(list []img.RGBA) img.RGBA {

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

func (c *canopy) result(size int) []img.RGBA {
	var centerPoints []img.RGBA
	for i := 0; i < len(c.Canopies); i++ {
		centerPoints = append(centerPoints, c.average(c.Canopies[i].Points))
	}
	if size > 0 {
		canopyInstance := canopy{
			AllPoints: centerPoints,
		}
		canopyInstance.run()
		return canopyInstance.result(size - 1)
	}

	return centerPoints
}
