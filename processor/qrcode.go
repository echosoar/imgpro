package processor

import (
	"fmt"
	"math"
	"sort"

	img "github.com/echosoar/imgpro/core"
	method "github.com/echosoar/imgpro/method"
)

// QRCode device
func QRCodeProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:          []string{"qrcode"},
		PreConditions: []string{"rgba", "frame", "width", "height"},
		Runner:        qrCodeRunner,
	})
}

func qrCodeRunner(core *img.Core) map[string]img.Value {
	rgba := core.Result["rgba"].Rgba
	width := core.Result["width"].Int
	height := core.Result["height"].Int
	frame := core.Result["frame"].Int
	for frameIndex := 0; frameIndex < frame; frameIndex++ {
		qr := QRCode{
			Width:  width,
			Height: height,
			Pixels: rgba[frameIndex],
		}
		qr.Run()
	}
	item := make(map[string]img.Value)
	item["Content"] = img.Value{
		Type:   img.ValueTypeString,
		String: "test",
	}
	item["Rect"] = img.Value{
		Type: img.ValueTypeRect,
		Rect: []img.ValuePosition{
			{
				X: 0,
				Y: 0,
			},
			{
				X: 0,
				Y: 0,
			},
			{
				X: 0,
				Y: 0,
			},
			{
				X: 0,
				Y: 0,
			},
		},
	}
	return map[string]img.Value{
		"qrcode": {
			Type: img.ValueTypeList,
			List: []img.Value{
				{
					Values: item,
				},
			},
		},
	}
}

type QRCode struct {
	Width      int
	Height     int
	Pixels     []img.RGBA
	greyPixels []uint8
	regions    []*QRCodeRegion
	corners    []*QRCodeCorner
	codeItems  []*QRCodeItem
}

type QRCodeRegion struct {
	posi         img.ValuePosition
	size         int
	cornersIndex int
}
type QRCodeCorner struct {
	corners      []img.ValuePosition
	matrix       []float64
	center       img.ValuePosition
	attachToCode int
}

type QRCodeItem struct {
	corners []*QRCodeCorner
	matrix  []float64
	version int
	size    int
	Pixels  []int
}

func (qr *QRCode) Run() {
	qr.grayscale()
	qr.binarization()
	qr.findCorners()
	qr.polymerizate()
	qr.output()
	fmt.Println("qr", len(qr.greyPixels))
}

// 灰度化
func (qr *QRCode) grayscale() {
	greyPixels := make([]uint8, len(qr.Pixels))
	for index, rgba := range qr.Pixels {
		greyPixels[index] = method.RGBAToGrey(rgba)
	}
	qr.greyPixels = greyPixels
}

// 去噪
func (qr *QRCode) denoising() {

}

// 畸变矫正
func (qr *QRCode) distortionCorrection() {

}

const (
	qrThresholdSDen = 8
	qrThresholdT    = 5
)

const (
	qrPixelWhite      = 0
	qrPixelBlack      = 1
	qrPixelRegion     = 2
	qrRegionIndexDiff = 2
	qrMaxRegion       = 256
	qrCornerSize      = 7
)

// 二值化
func (qr *QRCode) binarization() {
	avgW := 0
	avgU := 0
	thresholds := qr.Width / qrThresholdSDen
	if thresholds < 1 {
		thresholds = 1
	}

	for line := 0; line < qr.Height; line++ {
		lineStartPixelIndex := qr.Width * line
		rowGreyPixels := qr.greyPixels[lineStartPixelIndex : lineStartPixelIndex+qr.Width]
		rowAverage := make([]int, qr.Width)
		for row := 0; row < qr.Width; row++ {
			var w, u int
			if line&1 == 1 {
				w = row
				u = qr.Width - 1 - row
			} else {
				w = qr.Width - 1 - row
				u = row
			}

			avgW = (avgW*(thresholds-1))/thresholds + int(rowGreyPixels[w])
			avgU = (avgU*(thresholds-1))/thresholds + int(rowGreyPixels[u])

			rowAverage[w] += avgW
			rowAverage[u] += avgU
		}

		for row := 0; row < qr.Width; row++ {
			if int(rowGreyPixels[row]) < rowAverage[row]*(100-qrThresholdT)/(200*thresholds) {
				rowGreyPixels[row] = qrPixelBlack
			} else {
				rowGreyPixels[row] = qrPixelWhite
			}
		}
	}
}

// 寻找用于定位的角
func (qr *QRCode) findCorners() {
	for line := 0; line < qr.Height; line++ {
		lineStartPixelIndex := qr.Width * line
		rowGreyPixels := qr.greyPixels[lineStartPixelIndex : lineStartPixelIndex+qr.Width]
		lastColor := 0
		runLength := 0
		runCount := 0
		pb := make([]int, 5)
		// 宽度按照 1:1:3:1:1 的规模，哪怕倾斜之后也是如此
		check := [5]int{1, 1, 3, 1, 1}

		for row := 0; row < qr.Width; row++ {
			color := 0
			if rowGreyPixels[row] > 0 {
				color = 1
			}

			if row > 0 && color != lastColor {
				for i := 0; i < 4; i++ {
					pb[i] = pb[i+1]
				}
				pb[4] = runLength
				runLength = 0
				runCount++

				if color == 0 && runCount >= 5 {

					var avg, err int
					ok := true
					avg = (pb[0] + pb[1] + pb[3] + pb[4]) / 4

					err = avg * 3 / 4

					for i := 0; i < 5; i++ {
						if pb[i] < check[i]*avg-err || pb[i] > check[i]*avg+err {
							ok = false
							break
						}
					}
					if ok {
						qr.checkIsCorner(line, row, pb)
					}
				}
			}
			runLength++
			lastColor = color
		}
	}
}

func (qr *QRCode) checkIsCorner(line, row int, pb []int) {
	regionRightLineIndex := qr.regionCode(line, row-pb[4])
	regionCenterIndex := qr.regionCode(line, row-pb[4]-pb[3]-pb[2])
	regionLeftLineIndex := qr.regionCode(line, row-pb[4]-pb[3]-pb[2]-pb[1]-pb[0])

	if regionRightLineIndex < 0 || regionCenterIndex < 0 || regionLeftLineIndex < 0 || regionLeftLineIndex != regionRightLineIndex || regionLeftLineIndex == regionCenterIndex {
		return
	}

	centerRegion := qr.getRegionByRegionIndex(regionCenterIndex)
	lineRegion := qr.getRegionByRegionIndex(regionRightLineIndex)
	// 已经放在某一个 corner 之中了
	if centerRegion.cornersIndex >= 0 || lineRegion.cornersIndex >= 0 {
		return
	}

	ratio := centerRegion.size * 100 / lineRegion.size
	if ratio < 10 || ratio > 70 {
		return
	}
	qr.newCorner(regionRightLineIndex, lineRegion, centerRegion)
}

func (qr *QRCode) regionCode(line, row int) int {
	if row < 0 || line < 0 || row >= qr.Width || line >= qr.Height {
		return -1
	}
	index := line*qr.Width + row
	pixel := int(qr.greyPixels[index])
	if pixel >= 2 {
		return pixel
	}
	if pixel == qrPixelWhite {
		return -1
	}
	if len(qr.regions) >= qrMaxRegion {
		return -1
	}

	region := &QRCodeRegion{}

	qr.regions = append(qr.regions, region)
	regionIndex := len(qr.regions) - 1 + qrRegionIndexDiff // 避免与像素重合

	region.posi.X = row
	region.posi.Y = line
	region.size = 0
	region.cornersIndex = -1

	qr.fillRegion(row, line, pixel, regionIndex, region, 0)
	return regionIndex
}

func (qr *QRCode) fillRegion(row, line, pixel, regionIndex int, region *QRCodeRegion, depth int) {
	left := row
	right := row
	rowGreyPixels := qr.greyPixels[line*qr.Width : line*qr.Width+qr.Width]

	if depth >= 1024 {
		return
	}

	// 如果左边的像素和自己一样，就向左移动
	for left > 0 && int(rowGreyPixels[left-1]) == pixel {
		left--
	}
	// 如果右边的像素和自己一样，就向左移动
	for right < qr.Width-1 && int(rowGreyPixels[right+1]) == pixel {
		right++
	}

	changedPixel := 0
	// 填充
	for i := left; i <= right; i++ {
		newPixel := uint8(regionIndex)
		if rowGreyPixels[i] != newPixel {
			rowGreyPixels[i] = uint8(regionIndex)
			changedPixel++
		}
	}

	if changedPixel <= 0 {
		return
	}
	region.size += changedPixel

	// 往上面一行、下面一行进行检测
	lineAutoCheck := []int{-1, 1}
	for _, lineDiff := range lineAutoCheck {
		newLine := line + lineDiff
		if newLine <= 0 || newLine >= qr.Height {
			continue
		}
		rowGreyPixels := qr.greyPixels[newLine*qr.Width : newLine*qr.Width+qr.Width]
		for i := left; i <= right; i++ {
			if int(rowGreyPixels[i]) == pixel {
				qr.fillRegion(i, newLine, pixel, regionIndex, region, depth+1)
			}
		}
	}
}

func (qr *QRCode) getRegionByRegionIndex(index int) *QRCodeRegion {
	return qr.regions[index-qrRegionIndexDiff]
}

func (qr *QRCode) newCorner(lineRegionIndex int, lineRegion *QRCodeRegion, centerRegion *QRCodeRegion) {
	corner := &QRCodeCorner{
		attachToCode: -1,
	}
	qr.corners = append(qr.corners, corner)
	cornerIndex := len(qr.corners)
	lineRegion.cornersIndex = cornerIndex
	centerRegion.cornersIndex = cornerIndex
	lineRegionIndexUint := uint8(lineRegionIndex)
	// 寻找四个顶点
	allPixelPositions := make([]img.ValuePosition, 0)

	firstCornerMaxDistance := 0
	firstCornerIndex := 0

	for line := 0; line < qr.Height; line++ {
		for row := 0; row < qr.Width; row++ {
			index := line*qr.Width + row
			if qr.greyPixels[index] == lineRegionIndexUint {
				allPixelPositions = append(allPixelPositions, img.ValuePosition{
					X: row,
					Y: line,
				})
				pixelIndex := len(allPixelPositions) - 1
				distance := int(math.Pow(float64(row-lineRegion.posi.X), 2.0) + math.Pow(float64(line-lineRegion.posi.Y), 2.0))
				if distance > firstCornerMaxDistance {
					firstCornerMaxDistance = distance
					firstCornerIndex = pixelIndex
				}
			}
		}
	}

	firstPixelPosi := allPixelPositions[firstCornerIndex]

	secondCornerMaxDistance := 0
	secondCornerIndex := 0
	for index, posi := range allPixelPositions {
		distance := int(math.Pow(float64(posi.X-firstPixelPosi.X), 2.0) + math.Pow(float64(posi.Y-firstPixelPosi.Y), 2.0))
		if distance > secondCornerMaxDistance {
			secondCornerMaxDistance = distance
			secondCornerIndex = index
		}
	}

	secondPixelPosi := allPixelPositions[secondCornerIndex]

	thirdCornerMaxDistance := 0
	thirdCornerIndex := 0
	for index, posi := range allPixelPositions {
		distance := int(math.Sqrt(math.Pow(float64(posi.X-firstPixelPosi.X), 2.0)+math.Pow(float64(posi.Y-firstPixelPosi.Y), 2.0)) + math.Sqrt(math.Pow(float64(posi.X-secondPixelPosi.X), 2.0)+math.Pow(float64(posi.Y-secondPixelPosi.Y), 2.0)))
		if distance > thirdCornerMaxDistance {
			thirdCornerMaxDistance = distance
			thirdCornerIndex = index
		}
	}

	thirdPixelPosi := allPixelPositions[thirdCornerIndex]

	lineInfo := method.PointsToLine(firstPixelPosi, secondPixelPosi)
	if lineInfo.IsX {
		secondPixelPosi, thirdPixelPosi = thirdPixelPosi, secondPixelPosi
	}

	k, b := lineInfo.K, lineInfo.B
	thirdIsUnderLine := (float64(thirdPixelPosi.Y) - k*float64(thirdPixelPosi.X)) > b

	fouthCornerMaxDistance := 0
	fouthCornerIndex := 0
	for index, posi := range allPixelPositions {
		posiIsUnderLine := (float64(posi.Y) - k*float64(posi.X)) > b
		if thirdIsUnderLine == posiIsUnderLine {
			continue
		}
		distance := int(math.Sqrt(math.Pow(float64(posi.X-firstPixelPosi.X), 2.0)+math.Pow(float64(posi.Y-firstPixelPosi.Y), 2.0)) + math.Sqrt(math.Pow(float64(posi.X-secondPixelPosi.X), 2.0)+math.Pow(float64(posi.Y-secondPixelPosi.Y), 2.0)))
		if distance > fouthCornerMaxDistance {
			fouthCornerMaxDistance = distance
			fouthCornerIndex = index
		}
	}

	fouthCornerPosi := allPixelPositions[fouthCornerIndex]

	corners := []img.ValuePosition{
		firstPixelPosi,
		secondPixelPosi,
		thirdPixelPosi,
		fouthCornerPosi,
	}

	// 按照顺时针方向，左上角为第一个
	sort.SliceStable(corners, func(i, j int) bool {
		return corners[j].Y > corners[i].Y
	})
	if corners[0].X > corners[1].X {
		corners[1], corners[0] = corners[0], corners[1]
	}
	if corners[3].X > corners[2].X {
		corners[3], corners[2] = corners[2], corners[3]
	}
	corner.corners = corners
	corner.matrix = method.PerspectiveMap(&corner.corners, qrCornerSize, qrCornerSize)
	corner.center = method.PerspectiveTransform(&corner.matrix, qrCornerSize/2, qrCornerSize/2)
}

func (qr *QRCode) output() {
	rgba := make([]img.RGBA, len(qr.greyPixels))
	for index, pixel := range qr.greyPixels {
		if pixel <= 0 {
			rgba[index] = method.ColorListToRGBA([]int{255, 255, 255, 255})
		} else {
			rgba[index] = method.ColorListToRGBA([]int{0, 0, 0, 255})
		}
	}
	method.OutputToImg("./ignore_qrout.jpg", qr.Width, qr.Height, rgba)
}

// 聚合多个 corner，确认一张二维码
// 一个 corner 有横向和纵向相对应的 corner的时候，才有可能是一个二维码
func (qr *QRCode) polymerizate() {
	for index, corner := range qr.corners {
		if corner.attachToCode >= 0 {
			continue
		}
		horizontalOppositeCorners := make([]*QRCodeCorner, 0)
		horizontalOppositeCornerDistance := make([]float64, 0)
		verticalOppositeCorners := make([]*QRCodeCorner, 0)
		verticalOppositeCornerDistance := make([]float64, 0)
		for testIndex, testCorner := range qr.corners {
			if testIndex == index || testCorner.attachToCode >= 0 {
				continue
			}
			newPosi := method.PerspectiveTransformBack(&corner.matrix, testCorner.center)

			newPosi.X = int(math.Abs(float64(newPosi.X) - qrCornerSize/2))
			newPosi.Y = int(math.Abs(float64(newPosi.Y) - qrCornerSize/2))
			if float64(newPosi.X) < 0.2*float64(newPosi.Y) {
				horizontalOppositeCorners = append(horizontalOppositeCorners, testCorner)
				horizontalOppositeCornerDistance = append(horizontalOppositeCornerDistance, float64(newPosi.Y))
			}
			if float64(newPosi.Y) < 0.2*float64(newPosi.X) {
				verticalOppositeCorners = append(verticalOppositeCorners, testCorner)
				verticalOppositeCornerDistance = append(verticalOppositeCornerDistance, float64(newPosi.X))
			}
		}
		if len(horizontalOppositeCorners) == 0 || len(verticalOppositeCorners) == 0 {
			continue
		}

		bestScore := 0.0
		bestH := -1
		bestV := -1
		for hIndex := range horizontalOppositeCorners {
			for vIndex := range verticalOppositeCorners {
				score := math.Abs(1.0 - horizontalOppositeCornerDistance[hIndex]/verticalOppositeCornerDistance[vIndex])
				if score > 2.5 {
					continue
				}
				if bestH < 0 || score < bestScore {
					bestH = hIndex
					bestV = vIndex
					bestScore = score
				}
			}
		}

		if bestH < 0 || bestV < 0 {
			continue
		}
		qr.polymerizateExec(corner, horizontalOppositeCorners[bestH], verticalOppositeCorners[bestV])
	}
}

// 将3个corner处理成一个二维码
func (qr *QRCode) polymerizateExec(centerCorner *QRCodeCorner, hCorner *QRCodeCorner, vCorner *QRCodeCorner) {
	// 转换顺序，顺时针
	// kVH := float64(vCorner.center.Y-hCorner.center.Y) / float64(vCorner.center.X-hCorner.center.X)
	// kCH := float64(centerCorner.center.Y-hCorner.center.Y) / float64(centerCorner.center.X-hCorner.center.X)
	// 避免除以零，转换上述公式
	if float64(vCorner.center.Y-hCorner.center.Y)*float64(centerCorner.center.X-hCorner.center.X) < float64(centerCorner.center.Y-hCorner.center.Y)*float64(vCorner.center.X-hCorner.center.X) {
		hCorner, vCorner = vCorner, hCorner
	}

	// 横向的左边与纵向的 上边交点就是第四个点，为什么不能用横向的右边与纵向的下边，是因为避免外围影响
	isExistsPoint, pointInfo := method.LineIntesect(
		method.PointsToLine(hCorner.corners[0], hCorner.corners[1]),
		method.PointsToLine(vCorner.corners[0], vCorner.corners[3]),
	)
	// isExistsPoint, pointInfo := method.LineIntesect(
	// 	method.PointsToLine(hCorner.corners[2], hCorner.corners[3]),
	// 	method.PointsToLine(vCorner.corners[1], vCorner.corners[2]),
	// )
	if !isExistsPoint {
		return
	}
	corners := []*QRCodeCorner{hCorner, centerCorner, vCorner}

	qrCodeItem := &QRCodeItem{
		corners: corners,
	}
	qr.codeItems = append(qr.codeItems, qrCodeItem)
	codeItemIndex := len(qr.codeItems) - 1
	centerCorner.attachToCode = codeItemIndex
	hCorner.attachToCode = codeItemIndex
	vCorner.attachToCode = codeItemIndex

	// TODO: 检测横纵两条线是否在 center 相交
	// TODO: 获取二维码版本
	qrCodeItem.version = 5
	qrCodeItem.size = qrCodeItem.version*4 + 17
	qrCodeItem.Pixels = make([]int, qrCodeItem.size*qrCodeItem.size)

	// 获取图像透视矩阵
	posies := []img.ValuePosition{corners[1].corners[0], corners[2].corners[0], pointInfo, corners[0].corners[0]}
	matrix := method.PerspectiveMap(&posies, float64(qrCodeItem.size)-float64(qrCornerSize/2), float64(qrCodeItem.size)-float64(qrCornerSize/2))
	qrCodeItem.matrix = matrix
	qr.autoAdjustmentMatrix(qrCodeItem, &matrix)
	rgba := make([]img.RGBA, qrCodeItem.size*qrCodeItem.size)
	// 获取图像透视像素点
	for line := 0; line < qrCodeItem.size; line++ {
		for row := 0; row < qrCodeItem.size; row++ {
			pixelIndex := line*qrCodeItem.size + row
			qrCodeItem.Pixels[pixelIndex] = qrPixelWhite
			rgba[pixelIndex] = img.RGBA{R: 255, G: 255, B: 255, A: 255}
			// 之所以加0.5是为了获取到像素中心
			point := method.PerspectiveTransform(&qrCodeItem.matrix, float64(row)+0.5, float64(line)+0.5)
			// default isWhite
			if point.Y >= 0 && point.Y < qr.Height && point.X >= 0 && point.X < qr.Width {
				index := point.Y*qr.Width + point.X
				if qr.greyPixels[index] != 0 {
					qrCodeItem.Pixels[pixelIndex] = qrPixelBlack
					rgba[pixelIndex] = img.RGBA{R: 0, G: 0, B: 0, A: 255}
				}
			}
		}
	}
	method.OutputToImg("./ignore_target.jpg", qrCodeItem.size, qrCodeItem.size, rgba)
}

// 自动调整矩阵
func (qr *QRCode) autoAdjustmentMatrix(qrItem *QRCodeItem, matrix *[]float64) {
	score := qr.scoreMatrix(qrItem, matrix)
	matrixValue := *matrix
	matrixLen := len(matrixValue)
	adjustmentSteps := make([]float64, matrixLen)
	for index, matrixComponent := range matrixValue {
		adjustmentSteps[index] = matrixComponent * 0.1
	}

	for pass := 0; pass < 20; pass++ {
		for j := 0; j < matrixLen*2; j++ {
			i := j >> 1
			step := adjustmentSteps[i]
			newMatrix := (*matrix)[:]
			if j&1 == 1 {
				newMatrix[i] += step
			} else {
				newMatrix[i] -= step
			}

			testScore := qr.scoreMatrix(qrItem, &newMatrix)
			if testScore > score {

				score = testScore
				(*matrix)[i] = newMatrix[i]
			}
		}
		for i := 0; i < matrixLen; i++ {
			adjustmentSteps[i] *= 0.5
		}
	}
	fmt.Println("matrix", matrix, score)
	qrItem.matrix = *matrix
}

// 给透视矩阵打分
func (qr *QRCode) scoreMatrix(qrItem *QRCodeItem, matrix *[]float64) int {
	score := 0
	// 给定位角打分
	score += qr.scorePositioningAngle(qrItem, matrix)
	return score
}

// 给定位角打分
func (qr *QRCode) scorePositioningAngle(qrItem *QRCodeItem, matrix *[]float64) int {
	score := 0
	// 给定位角打分
	score += qr.scoreArea(qrItem, matrix, 2, 2, 5, 5, true)
	score += qr.scoreArea(qrItem, matrix, qrItem.size-5, 2, qrItem.size-2, 5, true)
	// 左下角
	score += qr.scoreArea(qrItem, matrix, 2, qrItem.size-5, 5, qrItem.size-2, true)
	// black
	score += qr.scoreArea(qrItem, matrix, 0, qrItem.size-7, 1, qrItem.size, true)
	score += qr.scoreArea(qrItem, matrix, 6, qrItem.size-7, 7, qrItem.size, true)
	// empty
	score += qr.scoreArea(qrItem, matrix, 1, qrItem.size-6, 2, qrItem.size-1, false)
	score += qr.scoreArea(qrItem, matrix, 5, qrItem.size-6, 6, qrItem.size-1, false)
	return score
}

// 给区域打分打分
func (qr *QRCode) scoreArea(qrItem *QRCodeItem, matrix *[]float64, fromX, fromY, targetX, targetY int, isBlack bool) int {
	score := 0
	// 给定位角打分
	for x := fromX; x < targetX; x++ {
		for y := fromY; y < targetY; y++ {
			point := method.PerspectiveTransform(matrix, float64(x)+0.5, float64(y)+0.5)
			if point.Y < 0 || point.Y >= qr.Height || point.X < 0 || point.X >= qr.Width {
				continue
			}
			index := point.Y*qr.Width + point.X
			if int(qr.greyPixels[index]) != 0 {
				if isBlack {
					score++
				} else {
					score--
				}

			} else {
				if isBlack {
					score--
				} else {
					score++
				}
			}
		}
	}
	return score
}
