package method

import (
	"math"

	"github.com/echosoar/imgpro/core"
)

type LineInfo struct {
	K   float64
	B   float64
	IsX bool
	X   float64
}

// 两点求线段，返回 k 和 b 以及是否为横线（即 y = kx + b） 或 x = num
func PointsToLine(aX, aY, bX, bY float64) LineInfo {
	if aX == bX {
		return LineInfo{0.0, 0.0, true, aX}
	}
	k := (bY - aY) / (bX - aX)
	b := aY - k*(aX)
	return LineInfo{k, b, false, 0.0}
}

// 两条线段是否有交点，以及交点位置
func LineIntesect(lineA LineInfo, lineB LineInfo) (bool, core.ValuePosition) {
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

func PointDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

func PointDistanceInt(x1, y1, x2, y2 int) int {
	return int(math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2)))
}
