package method

import (
	img "github.com/echosoar/imgpro/core"
)

type CanopyPoint struct {
	Center img.RGBA
	Points []img.RGBA
}

func (c *CanopyPoint) Add(point img.RGBA) {
	c.Points = append(c.Points, point)
}

type Canopy struct {
	AllPoints []img.RGBA
	Canopies  []*CanopyPoint
	T1        float64
	T2        float64
}

func (c *Canopy) getT1T2() {
	center := AverageColor(c.AllPoints)
	var allDistance float64 = 0
	for i := 0; i < len(c.AllPoints); i++ {
		allDistance += ManhattanDistance(center, c.AllPoints[i])
	}
	c.T1 = allDistance / float64(len(c.AllPoints))
	c.T2 = c.T1 * 2 / 3
}

func (c *Canopy) Run() {
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
				distance := ManhattanDistance(canopyCenter, current)

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
				c.Canopies = append(c.Canopies, &CanopyPoint{
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

func (c *Canopy) getRandom(allPoints []img.RGBA) {
	index := len(allPoints) / 2
	current := allPoints[index]

	c.AllPoints = append(allPoints[:index], allPoints[index+1:]...)
	c.Canopies = append(c.Canopies, &CanopyPoint{
		Center: current,
		Points: []img.RGBA{current},
	})
}

func (c *Canopy) Result(size int) []img.RGBA {
	var centerPoints []img.RGBA
	for i := 0; i < len(c.Canopies); i++ {
		centerPoints = append(centerPoints, AverageColor(c.Canopies[i].Points))
	}
	if size > 0 {
		canopyInstance := Canopy{
			AllPoints: centerPoints,
		}
		canopyInstance.Run()
		return canopyInstance.Result(size - 1)
	}

	return centerPoints
}
