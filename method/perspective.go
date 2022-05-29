package method

import "github.com/echosoar/imgpro/core"

func PerspectiveMap(posis *[]core.ValuePosition, targetWidth float64, targetHeight float64) []float64 {
	x0 := float64((*posis)[0].X)
	y0 := float64((*posis)[0].Y)
	x1 := float64((*posis)[1].X)
	y1 := float64((*posis)[1].Y)
	x2 := float64((*posis)[2].X)
	y2 := float64((*posis)[2].Y)
	x3 := float64((*posis)[3].X)
	y3 := float64((*posis)[3].Y)
	widtheightDenominatorominator := targetWidth * (x2*y3 - x3*y2 + (x3-x2)*y1 + x1*(y2-y3))
	heightDenominator := targetHeight * (x2*y3 + x1*(y2-y3) - x3*y2 + (x3-x2)*y1)
	matrix := make([]float64, 8)
	matrix[0] = (x1*(x2*y3-x3*y2) + x0*(-x2*y3+x3*y2+(x2-x3)*y1) + x1*(x3-x2)*y0) / widtheightDenominatorominator
	matrix[1] = -(x0*(x2*y3+x1*(y2-y3)-x2*y1) - x1*x3*y2 + x2*x3*y1 + (x1*x3-x2*x3)*y0) / heightDenominator
	matrix[2] = x0
	matrix[3] = (y0*(x1*(y3-y2)-x2*y3+x3*y2) + y1*(x2*y3-x3*y2) + x0*y1*(y2-y3)) / widtheightDenominatorominator
	matrix[4] = (x0*(y1*y3-y2*y3) + x1*y2*y3 - x2*y1*y3 + y0*(x3*y2-x1*y2+(x2-x3)*y1)) / heightDenominator
	matrix[5] = y0
	matrix[6] = (x1*(y3-y2) + x0*(y2-y3) + (x2-x3)*y1 + (x3-x2)*y0) / widtheightDenominatorominator
	matrix[7] = (-x2*y3 + x1*y3 + x3*y2 + x0*(y1-y2) - x3*y1 + (x2-x1)*y0) / heightDenominator
	return matrix
}

func PerspectiveTransform(matrix *[]float64, targetU float64, targetV float64) core.ValuePosition {
	denominator := (*matrix)[6]*targetU + (*matrix)[7]*targetV + 1.0
	newPosi := core.ValuePosition{X: 0, Y: 0}
	newPosi.X = int(((*matrix)[0]*targetU + (*matrix)[1]*targetV + (*matrix)[2]) / denominator)
	newPosi.Y = int(((*matrix)[3]*targetU + (*matrix)[4]*targetV + (*matrix)[5]) / denominator)
	return newPosi
}

func PerspectiveTransformBack(matrix *[]float64, posi core.ValuePosition) core.ValuePosition {
	x := float64(posi.X)
	y := float64(posi.Y)
	newPosi := core.ValuePosition{X: 0, Y: 0}
	den := -(*matrix)[0]*(*matrix)[7]*y + (*matrix)[1]*(*matrix)[6]*y + ((*matrix)[3]*(*matrix)[7]-(*matrix)[4]*(*matrix)[6])*x + (*matrix)[0]*(*matrix)[4] - (*matrix)[1]*(*matrix)[3]
	newPosi.X = int(-((*matrix)[1]*(y-(*matrix)[5]) - (*matrix)[2]*(*matrix)[7]*y + ((*matrix)[5]*(*matrix)[7]-(*matrix)[4])*x + (*matrix)[2]*(*matrix)[4]) / den)
	newPosi.Y = int(((*matrix)[0]*(y-(*matrix)[5]) - (*matrix)[2]*(*matrix)[6]*y + ((*matrix)[5]*(*matrix)[6]-(*matrix)[3])*x + (*matrix)[2]*(*matrix)[3]) / den)
	return newPosi
}
