package method

import (
	"fmt"

	"github.com/echosoar/imgpro/core"
)

type LineInfo struct {
	K   float64
	B   float64
	IsX bool
	X   float64
}

// 两点求线段，返回 k 和 b 以及是否为横线（即 y = kx + b） 或 x = num
func PointsToLine(pointA core.ValuePosition, pointB core.ValuePosition) LineInfo {
	if pointA.X == pointB.X {
		return LineInfo{0.0, 0.0, true, float64(pointA.X)}
	}
	k := float64(pointB.Y-pointA.Y) / float64(pointB.X-pointA.X)
	b := float64(pointA.Y) - k*float64(pointA.X)
	return LineInfo{k, b, false, 0.0}
}

// 两条线段是否有交点，以及交点位置
func LineIntesect(lineA LineInfo, lineB LineInfo) (bool, core.ValuePosition) {
	fmt.Println("lineA", lineA, lineB)
	if lineA.IsX && lineB.IsX {
		return false, core.ValuePosition{}
	}
	x, y := 0.0, 0.0
	if lineA.IsX {
		x = lineA.X
		y = lineB.K*lineA.X + lineB.B
	} else if lineB.IsX {
		x = lineB.X
		y = lineA.K*lineB.X + lineA.B
	} else {
		if lineA.K == lineB.K {
			return false, core.ValuePosition{}
		}
		// y = k1 * x + b1
		// y = k2 * x + b2
		// k2 * x + b2 = k1 * x + b1
		// x = (b1 - b2) / (k2 - k1)
		x = (lineA.B - lineB.B) / (lineB.K - lineA.K)
		y = lineA.K*x + lineA.B
	}
	return true, core.ValuePosition{X: int(x), Y: int(y)}
}
