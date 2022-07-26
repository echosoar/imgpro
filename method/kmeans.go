package method

import (
	"math/rand"
	"sort"

	img "github.com/echosoar/imgpro/core"
)

type KMeans struct {
	AllPoints  []img.RGBA
	Center     []img.RGBA
	K          int
	valueRange []kMeansValueRange
	groupPoint [][]img.RGBA
}

type kMeansValueRange struct {
	Max int
	Min int
}

func (k *KMeans) Run() {
	k.initRange()
	k.initCenter()
	k.exec()
}

func (k *KMeans) initRange() {
	valueRange := make([]kMeansValueRange, 4)
	for _, point := range k.AllPoints {
		for attrIndex, value := range IterateRGBA(point) {
			attrRange := &valueRange[attrIndex]
			if value < attrRange.Min {
				attrRange.Min = value
			}
			if value > attrRange.Max {
				attrRange.Max = value
			}
		}
	}
	k.valueRange = valueRange
}

func (k *KMeans) exec() {
	isNewCenter := true
	for isNewCenter {
		k.group()
		isNewCenter = k.calcCenter()
	}
}

func (k *KMeans) group() {
	groupPoint := make([][]img.RGBA, k.K)
	for _, point := range k.AllPoints {
		minDistance := -1.0
		minDistanceCenterIndex := 0
		for centerIndex, center := range k.Center {
			distance := ManhattanDistance(center, point)
			if minDistance < 0 || distance < minDistance {
				minDistance = distance
				minDistanceCenterIndex = centerIndex
			}
		}
		groupPoint[minDistanceCenterIndex] = append(groupPoint[minDistanceCenterIndex], point)
	}
	k.groupPoint = groupPoint
}

func (k *KMeans) calcCenter() bool {
	isNewCenter := false
	for i := 0; i < k.K; i++ {
		clusterPoints := k.groupPoint[i]
		if len(clusterPoints) == 0 {
			pointColors := make([]int, 4)
			for index := range pointColors {
				attrRange := k.valueRange[index]
				pointColors[index] = attrRange.Min + rand.Intn(attrRange.Max-attrRange.Min)
			}
			k.Center[i] = ColorListToRGBA(pointColors)
			isNewCenter = true
			continue
		}

		newPoint := AverageColor(clusterPoints)
		if !IsSameColor(newPoint, k.Center[i]) {
			k.Center[i] = newPoint
			isNewCenter = true
		}
	}
	return isNewCenter
}

func (k *KMeans) initCenter() {
	need := k.K - len(k.Center)
	for i := 0; i < need; i++ {
		point := make([]int, 4)
		for attrIndex, attrRange := range k.valueRange {
			point[attrIndex] = attrRange.Min + rand.Intn(attrRange.Max-attrRange.Min)
		}
		k.Center = append(k.Center, ColorListToRGBA(point))
	}
}

func (k *KMeans) GetResult() []img.RGBA {
	limit := len(k.AllPoints) / 100
	curCenters := k.Center
	resultIndex := make([]int, 0)
	for index := range curCenters {
		if len(k.groupPoint[index]) < limit {
			continue
		}
		resultIndex = append(resultIndex, index)
	}
	sort.SliceStable(resultIndex, func(i, j int) bool {
		indexI := resultIndex[i]
		indexJ := resultIndex[j]
		return len(k.groupPoint[indexI]) > len(k.groupPoint[indexJ])
	})
	result := make([]img.RGBA, 0)
	for _, index := range resultIndex {
		result = append(result, curCenters[index])
	}

	return result
}
