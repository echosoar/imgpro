package processor

import (
	"errors"
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

	rgba := core.Result["rgba"].Frames
	width := core.Result["width"].Int
	height := core.Result["height"].Int
	frame := core.Result["frame"].Int
	results := []img.Value{}
	for frameIndex := 0; frameIndex < frame; frameIndex++ {
		qr := QRCode{
			Width:  width,
			Height: height,
			Pixels: rgba[frameIndex].Rgba,
		}
		qr.Run()
		result := qr.GetResult()
		results = append(results, result)
	}
	return map[string]img.Value{
		"qrcode": {
			Type:   img.ValueTypeFrames,
			Frames: results,
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
	corners              []*QRCodeCorner
	matrix               []float64
	version              int
	size                 int
	Pixels               []int
	errorCorrectionLevel int
	mask                 int
	blocksData           []int
	result               img.Value
	success              bool
	currentReadIndex     int
}

type QRDataInfo struct {
	alignmentPatterns []int
}

// refs: https://www.thonky.com/qr-code-tutorial/error-correction-table
type QRDataECC struct {
	tdc     int // Total Number of Data Codewords for this Version and EC Level
	ecc     int // EC Codewords Per Block
	blocks1 int // Number of Blocks in Group 1
	dc1     int // Number of Data Codewords in Each of Group 1's Blocks
	blocks2 int // Number of Blocks in Group 2
	dc2     int // Number of Data Codewords in Each of Group 2's Blocks
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
	qrMustMatch       = true
	qrIsBlack         = true
	qrIsWhite         = false
)

const (
	qrDataTypeEnd      = 0
	qrDataTypeNumric   = 1
	qrDataType8BitByte = 4
	qrDataTypeECI      = 7
)

var QRDataInfoList = []QRDataInfo{
	{},
	{alignmentPatterns: []int{}},
	{alignmentPatterns: []int{6, 18}},
	{alignmentPatterns: []int{6, 22}},
	{alignmentPatterns: []int{6, 26}},
	{alignmentPatterns: []int{6, 30}},
	{alignmentPatterns: []int{6, 34}},
	{alignmentPatterns: []int{6, 22, 38}},
	{alignmentPatterns: []int{6, 24, 42}},
	{alignmentPatterns: []int{6, 26, 46}},
	{alignmentPatterns: []int{6, 28, 50}},
	{alignmentPatterns: []int{6, 30, 54}},
	{alignmentPatterns: []int{6, 32, 58}},
	{alignmentPatterns: []int{6, 34, 62}},
	{alignmentPatterns: []int{6, 26, 46, 66}},
	{alignmentPatterns: []int{6, 26, 48, 70}},
	{alignmentPatterns: []int{6, 26, 50, 74}},
	{alignmentPatterns: []int{6, 30, 54, 78}},
	{alignmentPatterns: []int{6, 30, 56, 82}},
	{alignmentPatterns: []int{6, 30, 58, 86}},
	{alignmentPatterns: []int{6, 34, 62, 90}},
	{alignmentPatterns: []int{6, 28, 50, 72, 94}},
	{alignmentPatterns: []int{6, 26, 50, 74, 98}},
	{alignmentPatterns: []int{6, 30, 54, 78, 102}},
	{alignmentPatterns: []int{6, 28, 54, 80, 106}},
	{alignmentPatterns: []int{6, 32, 58, 84, 110}},
	{alignmentPatterns: []int{6, 30, 58, 86, 114}},
	{alignmentPatterns: []int{6, 34, 62, 90, 118}},
	{alignmentPatterns: []int{6, 26, 50, 74, 98, 122}},
	{alignmentPatterns: []int{6, 30, 54, 78, 102, 126}},
	{alignmentPatterns: []int{6, 26, 52, 78, 104, 130}},
	{alignmentPatterns: []int{6, 30, 56, 82, 108, 134}},
	{alignmentPatterns: []int{6, 34, 60, 86, 112, 138}},
	{alignmentPatterns: []int{6, 30, 58, 86, 114, 142}},
	{alignmentPatterns: []int{6, 34, 62, 90, 118, 146}},
	{alignmentPatterns: []int{6, 30, 54, 78, 102, 126, 150}},
	{alignmentPatterns: []int{6, 24, 50, 76, 102, 128, 154}},
	{alignmentPatterns: []int{6, 28, 54, 80, 106, 132, 158}},
	{alignmentPatterns: []int{6, 32, 58, 84, 110, 136, 162}},
	{alignmentPatterns: []int{6, 26, 54, 82, 110, 138, 166}},
	{alignmentPatterns: []int{6, 30, 58, 86, 114, 142, 170}},
}

var qrEccList = [][]QRDataECC{
	{},
	{ // version 1: L/M/Q/H 四个纠错级别，1/0/3/2
		{19, 7, 1, 19, 0, 0},
		{16, 10, 1, 16, 0, 0},
		{13, 13, 1, 13, 0, 0},
		{9, 17, 1, 9, 0, 0},
	},
	{ // version 2
		{34, 10, 1, 34, 0, 0},
		{28, 16, 1, 28, 0, 0},
		{22, 22, 1, 22, 0, 0},
		{16, 28, 1, 16, 0, 0},
	},
	{ // version 3
		{55, 15, 1, 55, 0, 0},
		{44, 26, 1, 44, 0, 0},
		{34, 18, 2, 17, 0, 0},
		{26, 22, 2, 13, 0, 0},
	},
	{ // version 4
		{80, 20, 1, 80, 0, 0},
		{64, 18, 2, 32, 0, 0},
		{48, 26, 2, 24, 0, 0},
		{36, 16, 4, 9, 0, 0},
	},
	{ // version 5
		{108, 26, 1, 108, 0, 0},
		{86, 24, 2, 43, 0, 0},
		{62, 18, 2, 15, 2, 16},
		{46, 22, 2, 11, 2, 12},
	},
	{ // version 6
		{136, 18, 2, 68, 0, 0},
		{108, 16, 4, 27, 0, 0},
		{76, 24, 4, 19, 0, 0},
		{60, 28, 4, 15, 0, 0},
	},
	{ // version 7
		{156, 20, 2, 78, 0, 0},
		{124, 18, 4, 31, 0, 0},
		{88, 18, 2, 14, 4, 15},
		{66, 26, 4, 13, 1, 14},
	},
	{ // version 8
		{194, 24, 2, 97, 0, 0},
		{154, 22, 2, 38, 2, 39},
		{110, 22, 4, 18, 2, 19},
		{86, 26, 4, 14, 2, 15},
	},
	{ // version 9
		{232, 30, 2, 116, 0, 0},
		{182, 22, 3, 36, 2, 37},
		{132, 20, 4, 16, 4, 17},
		{100, 24, 4, 12, 4, 13},
	},
	{ // version 10
		{274, 18, 2, 68, 2, 69},
		{216, 26, 4, 43, 1, 44},
		{154, 24, 6, 19, 2, 20},
		{122, 28, 6, 15, 2, 16},
	},
	{ // version 11
		{324, 20, 4, 81, 0, 0},
		{254, 30, 1, 50, 4, 51},
		{180, 28, 4, 22, 4, 23},
		{140, 24, 3, 12, 8, 13},
	},
	{ // version 12
		{370, 24, 2, 92, 2, 93},
		{290, 22, 6, 36, 2, 37},
		{206, 26, 4, 20, 6, 21},
		{158, 28, 7, 14, 4, 15},
	},
	{ // version 13
		{428, 26, 4, 107, 0, 0},
		{334, 22, 8, 37, 1, 38},
		{244, 24, 8, 20, 4, 21},
		{180, 22, 12, 11, 4, 12},
	},
	{ // version 14
		{461, 30, 3, 115, 1, 116},
		{365, 24, 4, 40, 5, 41},
		{261, 20, 11, 16, 5, 17},
		{197, 24, 11, 12, 5, 13},
	},
	{ // version 15
		{523, 22, 5, 87, 1, 88},
		{415, 24, 5, 41, 5, 42},
		{295, 30, 5, 24, 7, 25},
		{223, 24, 11, 12, 7, 13},
	},
	{ // version 16
		{589, 24, 5, 98, 1, 99},
		{453, 28, 7, 45, 3, 46},
		{325, 24, 15, 19, 2, 20},
		{253, 30, 3, 15, 13, 16},
	},
	{ // version 17
		{647, 28, 1, 107, 5, 108},
		{507, 28, 10, 46, 1, 47},
		{367, 28, 1, 22, 15, 23},
		{283, 28, 2, 14, 17, 15},
	},
	{ // version 18
		{721, 30, 5, 120, 1, 121},
		{563, 26, 9, 43, 4, 44},
		{397, 28, 17, 22, 1, 23},
		{313, 28, 2, 14, 19, 15},
	},
	{ // version 19
		{795, 28, 3, 113, 4, 114},
		{627, 26, 3, 44, 11, 45},
		{445, 26, 17, 21, 4, 22},
		{341, 26, 9, 13, 16, 14},
	},
	{ // version 20
		{861, 28, 3, 107, 5, 108},
		{669, 26, 3, 41, 13, 42},
		{485, 30, 15, 24, 5, 25},
		{385, 28, 15, 15, 10, 16},
	},
	{ // version 21
		{932, 28, 4, 116, 4, 117},
		{714, 26, 17, 42, 0, 0},
		{512, 28, 17, 22, 6, 23},
		{406, 30, 19, 16, 6, 17},
	},
	{ // version 22
		{1006, 28, 2, 111, 7, 112},
		{782, 28, 17, 46, 0, 0},
		{568, 30, 7, 24, 16, 25},
		{442, 24, 34, 13, 0, 0},
	},
	{ // version 23
		{1094, 30, 4, 121, 5, 122},
		{860, 28, 4, 47, 14, 48},
		{614, 30, 11, 24, 14, 25},
		{464, 30, 16, 15, 14, 16},
	},
	{ // version 24
		{1174, 30, 6, 117, 4, 118},
		{914, 28, 6, 45, 14, 46},
		{664, 30, 11, 24, 16, 25},
		{514, 30, 30, 16, 2, 17},
	},
	{ // version 25
		{1276, 26, 8, 106, 4, 107},
		{1000, 28, 8, 47, 13, 48},
		{718, 30, 7, 24, 22, 25},
		{538, 30, 22, 15, 13, 16},
	},
	{ // version 26
		{1370, 28, 10, 114, 2, 115},
		{1062, 28, 19, 46, 4, 47},
		{754, 28, 28, 22, 6, 23},
		{596, 30, 33, 16, 4, 17},
	},
	{ // version 27
		{1468, 30, 8, 122, 4, 123},
		{1128, 28, 22, 45, 3, 46},
		{808, 30, 8, 23, 26, 24},
		{628, 30, 12, 15, 28, 16},
	},
	{ // version 28
		{1531, 30, 3, 117, 10, 118},
		{1193, 28, 3, 45, 23, 46},
		{871, 30, 4, 24, 31, 25},
		{661, 30, 11, 15, 31, 16},
	},
	{ // version 29
		{1631, 30, 7, 116, 7, 117},
		{1267, 28, 21, 45, 7, 46},
		{911, 30, 1, 23, 37, 24},
		{701, 30, 19, 15, 26, 16},
	},
	{ // version 30
		{1735, 30, 5, 115, 10, 116},
		{1373, 28, 19, 47, 10, 48},
		{985, 30, 15, 24, 25, 25},
		{745, 30, 23, 15, 25, 16},
	},
	{ // version 31
		{1843, 30, 13, 115, 3, 116},
		{1455, 28, 2, 46, 29, 47},
		{1033, 30, 42, 24, 1, 25},
		{793, 30, 23, 15, 28, 16},
	},
	{ // version 32
		{1955, 30, 17, 115, 0, 0},
		{1541, 28, 10, 46, 23, 47},
		{1115, 30, 10, 24, 35, 25},
		{845, 30, 19, 15, 35, 16},
	},
	{ // version 33
		{2071, 30, 17, 115, 1, 116},
		{1631, 28, 14, 46, 21, 47},
		{1171, 30, 29, 24, 19, 25},
		{901, 30, 11, 15, 46, 16},
	},
	{ // version 34
		{2191, 30, 13, 115, 6, 116},
		{1725, 28, 14, 46, 23, 47},
		{1231, 30, 44, 24, 7, 25},
		{961, 30, 59, 16, 1, 17},
	},
	{ // version 35
		{2306, 30, 12, 121, 7, 122},
		{1812, 28, 12, 47, 26, 48},
		{1286, 30, 39, 24, 14, 25},
		{986, 30, 22, 15, 41, 16},
	},
	{ // version 36
		{2434, 30, 6, 121, 14, 122},
		{1914, 28, 6, 47, 34, 48},
		{1354, 30, 46, 24, 10, 25},
		{1054, 30, 2, 15, 64, 16},
	},
	{ // version 37
		{2566, 30, 17, 122, 4, 123},
		{1992, 28, 29, 46, 14, 47},
		{1426, 30, 49, 24, 10, 25},
		{1096, 30, 24, 15, 46, 16},
	},
	{ // version 38
		{2702, 30, 4, 122, 18, 123},
		{2102, 28, 13, 46, 32, 47},
		{1502, 30, 48, 24, 14, 25},
		{1142, 30, 42, 15, 32, 16},
	},
	{ // version 39
		{2812, 30, 20, 117, 4, 118},
		{2216, 28, 40, 47, 7, 48},
		{1582, 30, 43, 24, 22, 25},
		{1222, 30, 10, 15, 67, 16},
	},
	{ // version 40
		{2956, 30, 19, 118, 6, 119},
		{2334, 28, 18, 47, 31, 48},
		{1666, 30, 34, 24, 34, 25},
	},
}

func (qr *QRCode) Run() {
	qr.grayscale()
	qr.binarization()
	qr.findCorners()
	qr.polymerizate()
}

// 灰度化
func (qr *QRCode) grayscale() {
	greyPixels := make([]uint8, len(qr.Pixels))
	for index, rgba := range qr.Pixels {
		greyPixels[index] = method.RGBAToGrey(rgba)
	}
	qr.greyPixels = greyPixels
}

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
			// 避免定位点在最边上的情况
			if row == qr.Width-1 {
				color = 0
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

	lineInfo := method.PointsToLine(float64(firstPixelPosi.X), float64(firstPixelPosi.Y), float64(secondPixelPosi.X), float64(secondPixelPosi.Y))
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
		// distance := math.Abs(k*float64(posi.X)-float64(posi.Y)+b) / ab
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
		method.PointsToLine(float64(hCorner.corners[0].X), float64(hCorner.corners[0].Y), float64(hCorner.corners[1].X), float64(hCorner.corners[1].Y)),
		method.PointsToLine(float64(vCorner.corners[0].X), float64(vCorner.corners[0].Y), float64(vCorner.corners[3].X), float64(vCorner.corners[3].Y)),
	)
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

	// 获取二维码版本
	qrCodeItem.version = qr.getVersion(qrCodeItem)
	qrCodeItem.size = qrCodeItem.version*4 + 17
	qrCodeItem.Pixels = make([]int, qrCodeItem.size*qrCodeItem.size)

	// 超过7版本，就有 alignPattern 了
	if qrCodeItem.version >= 7 {
		areaRange := method.PointDistanceInt(centerCorner.corners[0].X, centerCorner.corners[0].Y, hCorner.corners[3].X, hCorner.corners[3].Y) * 3 / qrCodeItem.size
	checkPointInfo:
		for i := 0; i < areaRange; i++ {
			for j := 0; j < areaRange; j++ {
				// changeList := []int{i}
				// if i != 0 {
				// 	changeList = []int{-i, i}
				// }
				newPosi := img.ValuePosition{X: pointInfo.X + i, Y: pointInfo.Y + j}
				posiIndex := newPosi.X + newPosi.Y*qr.Width
				if qr.greyPixels[posiIndex] == 0 {
					continue
				}
				matched, newX, newY := qr.check111(qrCodeItem, newPosi.X, newPosi.Y)
				if matched {
					pointInfo.X = newX
					pointInfo.Y = newY
					break checkPointInfo
				}
			}
		}
	}

	// 获取图像透视矩阵
	posies := []img.ValuePosition{corners[1].corners[0], corners[2].corners[0], pointInfo, corners[0].corners[0]}
	matrix := method.PerspectiveMap(&posies, float64(qrCodeItem.size)-float64(qrCornerSize), float64(qrCodeItem.size)-float64(qrCornerSize))
	qrCodeItem.matrix = matrix
	qr.autoAdjustmentMatrix(qrCodeItem, &matrix)
	// 获取图像透视像素点
	for line := 0; line < qrCodeItem.size; line++ {
		for row := 0; row < qrCodeItem.size; row++ {
			pixelIndex := line*qrCodeItem.size + row
			qrCodeItem.Pixels[pixelIndex] = qrPixelWhite
			// 之所以加0.5是为了获取到像素中心
			point := method.PerspectiveTransform(&qrCodeItem.matrix, float64(row)+0.5, float64(line)+0.5)
			// default isWhite
			if point.Y >= 0 && point.Y < qr.Height && point.X >= 0 && point.X < qr.Width {
				index := point.Y*qr.Width + point.X
				if qr.greyPixels[index] != 0 {
					qrCodeItem.Pixels[pixelIndex] = qrPixelBlack
				}
			}
		}
	}
	qrCodeItem.decode()
}

// 检测是否是 1:1:1
func (qr *QRCode) check111(qrCodeItem *QRCodeItem, x, y int) (bool, int, int) {
	xi := x
	xj := x
	yi := y
	yj := y
	colorChangeXIndex := make([]int, 0)
	colorChangeYIndex := make([]int, 0)

	lastColor := qrIsBlack
	colorChangedCount := 0
	for xi >= 0 {
		xi--
		index := xi + y*qr.Width
		curColor := qrIsBlack
		if qr.greyPixels[index] == 0 {
			curColor = qrIsWhite
		}
		if curColor != lastColor {
			lastColor = curColor
			colorChangedCount++
			colorChangeXIndex = method.ConcatArray([]int{xi}, colorChangeXIndex)
			if colorChangedCount == 2 {
				break
			}
		}
	}
	lastColor = qrIsBlack
	colorChangedCount = 0

	for xj < qr.Width {
		xj++
		index := xj + y*qr.Width
		curColor := qrIsBlack
		if qr.greyPixels[index] == 0 {
			curColor = qrIsWhite
		}
		if curColor != lastColor {
			lastColor = curColor
			colorChangedCount++
			colorChangeXIndex = append(colorChangeXIndex, xj)
			if colorChangedCount == 2 {
				break
			}
		}
	}
	diff := []float64{
		float64(colorChangeXIndex[1] - colorChangeXIndex[0]),
		float64(colorChangeXIndex[2] - colorChangeXIndex[1]),
		float64(colorChangeXIndex[3] - colorChangeXIndex[2]),
	}
	if math.Abs(diff[1]/diff[0]-1) > 0.2 || math.Abs(diff[2]/diff[1]-1) > 0.2 {
		return false, 0, 0
	}

	lastColor = qrIsBlack
	colorChangedCount = 0
	for yi >= 0 {
		yi--
		index := x + yi*qr.Width
		curColor := qrIsBlack
		if qr.greyPixels[index] == 0 {
			curColor = qrIsWhite
		}
		if curColor != lastColor {
			lastColor = curColor
			colorChangedCount++
			colorChangeYIndex = method.ConcatArray([]int{yi}, colorChangeYIndex)
			if colorChangedCount == 2 {
				break
			}
		}
	}

	lastColor = qrIsBlack
	colorChangedCount = 0

	for yj < qr.Height {
		yj++
		index := x + yj*qr.Width
		curColor := qrIsBlack
		if qr.greyPixels[index] == 0 {
			curColor = qrIsWhite
		}
		if curColor != lastColor {
			lastColor = curColor
			colorChangedCount++
			colorChangeYIndex = append(colorChangeYIndex, yj)
			if colorChangedCount == 2 {
				break
			}
		}
	}
	diffY := []float64{
		float64(colorChangeYIndex[1] - colorChangeYIndex[0]),
		float64(colorChangeYIndex[2] - colorChangeYIndex[1]),
		float64(colorChangeYIndex[3] - colorChangeYIndex[2]),
	}
	if math.Abs(diffY[1]/diffY[0]-1) > 0.2 || math.Abs(diffY[2]/diffY[1]-1) > 0.2 {
		return false, 0, 0
	}

	xRange := (colorChangeXIndex[3] - colorChangeXIndex[0]) / 2
	yRange := (colorChangeYIndex[3] - colorChangeYIndex[0]) / 2

	startX := x - xRange
	endX := x + xRange
	if startX < 0 {
		startX = 0
	}
	if endX > qr.Width {
		endX = qr.Width
	}

	startY := y - yRange
	endY := y + yRange
	if startY < 0 {
		startY = 0
	}
	if endY > qr.Height {
		endY = qr.Height
	}
	pointMap := make([]int, (endX-startX)*(endY-startY))
	distance, newX, newY := qr.check111Black(qrCodeItem, &pointMap, startX, endX, startY, endY, x, y)
	if distance > 0 {
		return true, newX, newY
	}
	return false, x, y
}

func (qr *QRCode) check111Black(qrCodeItem *QRCodeItem, pointMap *[]int, startX, endX, startY, endY, x, y int) (float64, int, int) {
	index := x + y*qr.Width
	pointIndex := (endX-startX)*(y-startY) + (x - startX)
	if (*pointMap)[pointIndex] != 0 {
		return 0.0, 0, 0
	}
	if qr.greyPixels[index] == 0 {
		(*pointMap)[pointIndex] = 2 // white
		return 0.0, 0, 0
	}
	(*pointMap)[pointIndex] = 1 // black;
	distance := method.PointDistance(float64(x), float64(y), float64(qrCodeItem.corners[1].corners[2].X), float64(qrCodeItem.corners[1].corners[2].Y))
	resX := x
	resY := y
	offsetRange := []int{-1, 1}
	for _, xOffset := range offsetRange {
		for _, yOffset := range offsetRange {
			newX := x + xOffset
			newY := y + yOffset
			if newX < startX || newX >= endX || newY < startY || newY >= endY {
				continue
			}
			checkDistance, checkX, checkY := qr.check111Black(qrCodeItem, pointMap, startX, endX, startY, endY, newX, newY)
			if checkDistance > 0 && checkDistance < distance {
				distance = checkDistance
				resX = checkX
				resY = checkY
			}

		}
	}
	return distance, resX, resY
}

// 获取二维码版本
func (qr *QRCode) getVersion(qrItem *QRCodeItem) int {

	centerCornerX := float64(qrItem.corners[1].corners[2].X)
	centerCornerY := float64(qrItem.corners[1].corners[2].Y)
	centerCornerPerPixel := method.PointDistance(centerCornerX, centerCornerY, float64(qrItem.corners[1].corners[3].X), float64(qrItem.corners[1].corners[3].Y)) / 14
	bottomCornerX := float64(qrItem.corners[0].corners[1].X)
	bottomCornerY := float64(qrItem.corners[0].corners[1].Y)
	bottomCornerPerPixel := method.PointDistance(bottomCornerX, bottomCornerY, float64(qrItem.corners[0].corners[0].X), float64(qrItem.corners[0].corners[0].Y)) / 14
	rightCornerX := float64(qrItem.corners[2].corners[3].X)
	rightCornerY := float64(qrItem.corners[2].corners[3].Y)
	rightCornerPerPixel := method.PointDistance(rightCornerX, rightCornerY, float64(qrItem.corners[2].corners[0].X), float64(qrItem.corners[2].corners[0].Y)) / 14

	cornerCorners := [][]float64{
		{centerCornerX - centerCornerPerPixel, centerCornerY - centerCornerPerPixel, bottomCornerX - bottomCornerPerPixel, bottomCornerY + bottomCornerPerPixel},
		{centerCornerX - centerCornerPerPixel, centerCornerY - centerCornerPerPixel, rightCornerX + rightCornerPerPixel, rightCornerY - rightCornerPerPixel},
	}
	curMatchedBlackPoint := 0
	for _, corners := range cornerCorners {
		xChange := math.Abs(corners[0] - corners[2])
		yChange := math.Abs(corners[1] - corners[3])
		lineInfo := method.PointsToLine(corners[0], corners[1], corners[2], corners[3])
		calcIndexFun := func(num float64) (int, float64, float64) {
			return 0, 0.0, 0.0
		}
		var start, max float64
		if xChange > yChange {
			// 变化x，求y
			start = math.Min(corners[0], corners[2])
			max = start + xChange
			calcIndexFun = func(x float64) (int, float64, float64) {
				y := lineInfo.K*x + lineInfo.B
				return int(y)*qr.Width + int(x), x, y
			}
		} else {
			// 变化y 求x
			start = math.Min(corners[1], corners[3])
			max = start + yChange
			calcIndexFun = func(y float64) (int, float64, float64) {
				x := 0.0
				if lineInfo.IsX {
					x = lineInfo.X
				} else {
					x = (y - lineInfo.B) / lineInfo.K
				}
				return int(y)*qr.Width + int(x), x, y
			}
		}
		matchedBlackPoint := 0
		lastColor := qrIsBlack
		for i := start; i < max; i++ {
			index, _, _ := calcIndexFun(i)
			if qr.greyPixels[index] > 0 {
				// isBlack
				if lastColor == qrIsWhite {
					matchedBlackPoint++
				}
				lastColor = qrIsBlack
			} else {
				lastColor = qrIsWhite
			}
		}
		if matchedBlackPoint > curMatchedBlackPoint {
			curMatchedBlackPoint = matchedBlackPoint
		}
	}

	size := (curMatchedBlackPoint)*2 + 13

	version := (size - 15) / 4
	return version
}

var autoAdjustmentMatrixChangeIndex = [][]int{
	{0},
	{1},
	{2},
	{3},
	{4},
	{5},
	{6},
	{7},
	{0, 1},
	{0, 2},
	{0, 3},
	{0, 4},
	{0, 5},
	{0, 6},
	{0, 7},
	{1, 2},
	{1, 3},
	{1, 4},
	{1, 5},
	{1, 6},
	{1, 7},
	{2, 3},
	{2, 4},
	{2, 5},
	{2, 6},
	{2, 7},
	{3, 4},
	{3, 5},
	{3, 6},
	{3, 7},
	{4, 5},
	{4, 6},
	{4, 7},
	{5, 6},
	{5, 7},
	{6, 7},
	{0, 1, 2},
	{0, 1, 3},
	{0, 1, 4},
	{0, 1, 5},
	{0, 1, 6},
	{0, 1, 7},
	{0, 2, 3},
	{0, 2, 4},
	{0, 2, 5},
	{0, 2, 6},
	{0, 2, 7},
	{0, 3, 4},
	{0, 3, 5},
	{0, 3, 6},
	{0, 3, 7},
	{0, 4, 5},
	{0, 4, 6},
	{0, 4, 7},
	{0, 5, 6},
	{0, 5, 7},
	{0, 6, 7},
	{1, 2, 3},
	{1, 2, 4},
	{1, 2, 5},
	{1, 2, 6},
	{1, 2, 7},
	{1, 3, 4},
	{1, 3, 5},
	{1, 3, 6},
	{1, 3, 7},
	{1, 4, 5},
	{1, 4, 6},
	{1, 4, 7},
	{1, 5, 6},
	{1, 5, 7},
	{1, 6, 7},
	{2, 3, 4},
	{2, 3, 5},
	{2, 3, 6},
	{2, 3, 7},
	{2, 4, 5},
	{2, 4, 6},
	{2, 4, 7},
	{2, 5, 6},
	{2, 5, 7},
	{2, 6, 7},
	{3, 4, 5},
	{3, 4, 6},
	{3, 4, 7},
	{3, 5, 6},
	{3, 5, 7},
	{3, 6, 7},
	{4, 5, 6},
	{4, 5, 7},
	{4, 6, 7},
	{5, 6, 7},
	{0, 1, 2, 3},
	{0, 1, 2, 4},
	{0, 1, 2, 5},
	{0, 1, 2, 6},
	{0, 1, 2, 7},
	{0, 1, 3, 4},
	{0, 1, 3, 5},
	{0, 1, 3, 6},
	{0, 1, 3, 7},
	{0, 1, 4, 5},
	{0, 1, 4, 6},
	{0, 1, 4, 7},
	{0, 1, 5, 6},
	{0, 1, 5, 7},
	{0, 1, 6, 7},
	{0, 2, 3, 4},
	{0, 2, 3, 5},
	{0, 2, 3, 6},
	{0, 2, 3, 7},
	{0, 2, 4, 5},
	{0, 2, 4, 6},
	{0, 2, 4, 7},
	{0, 2, 5, 6},
	{0, 2, 5, 7},
	{0, 2, 6, 7},
	{0, 3, 4, 5},
	{0, 3, 4, 6},
	{0, 3, 4, 7},
	{0, 3, 5, 6},
	{0, 3, 5, 7},
	{0, 3, 6, 7},
	{0, 4, 5, 6},
	{0, 4, 5, 7},
	{0, 4, 6, 7},
	{0, 5, 6, 7},
	{1, 2, 3, 4},
	{1, 2, 3, 5},
	{1, 2, 3, 6},
	{1, 2, 3, 7},
	{1, 2, 4, 5},
	{1, 2, 4, 6},
	{1, 2, 4, 7},
	{1, 2, 5, 6},
	{1, 2, 5, 7},
	{1, 2, 6, 7},
	{1, 3, 4, 5},
	{1, 3, 4, 6},
	{1, 3, 4, 7},
	{1, 3, 5, 6},
	{1, 3, 5, 7},
	{1, 3, 6, 7},
	{1, 4, 5, 6},
	{1, 4, 5, 7},
	{1, 4, 6, 7},
	{1, 5, 6, 7},
	{2, 3, 4, 5},
	{2, 3, 4, 6},
	{2, 3, 4, 7},
	{2, 3, 5, 6},
	{2, 3, 5, 7},
	{2, 3, 6, 7},
	{2, 4, 5, 6},
	{2, 4, 5, 7},
	{2, 4, 6, 7},
	{2, 5, 6, 7},
	{3, 4, 5, 6},
	{3, 4, 5, 7},
	{3, 4, 6, 7},
	{3, 5, 6, 7},
	{4, 5, 6, 7},
	{0, 1, 2, 3, 4},
	{0, 1, 2, 3, 5},
	{0, 1, 2, 3, 6},
	{0, 1, 2, 3, 7},
	{0, 1, 2, 4, 5},
	{0, 1, 2, 4, 6},
	{0, 1, 2, 4, 7},
	{0, 1, 2, 5, 6},
	{0, 1, 2, 5, 7},
	{0, 1, 2, 6, 7},
	{0, 1, 3, 4, 5},
	{0, 1, 3, 4, 6},
	{0, 1, 3, 4, 7},
	{0, 1, 3, 5, 6},
	{0, 1, 3, 5, 7},
	{0, 1, 3, 6, 7},
	{0, 1, 4, 5, 6},
	{0, 1, 4, 5, 7},
	{0, 1, 4, 6, 7},
	{0, 1, 5, 6, 7},
	{0, 2, 3, 4, 5},
	{0, 2, 3, 4, 6},
	{0, 2, 3, 4, 7},
	{0, 2, 3, 5, 6},
	{0, 2, 3, 5, 7},
	{0, 2, 3, 6, 7},
	{0, 2, 4, 5, 6},
	{0, 2, 4, 5, 7},
	{0, 2, 4, 6, 7},
	{0, 2, 5, 6, 7},
	{0, 3, 4, 5, 6},
	{0, 3, 4, 5, 7},
	{0, 3, 4, 6, 7},
	{0, 3, 5, 6, 7},
	{0, 4, 5, 6, 7},
	{1, 2, 3, 4, 5},
	{1, 2, 3, 4, 6},
	{1, 2, 3, 4, 7},
	{1, 2, 3, 5, 6},
	{1, 2, 3, 5, 7},
	{1, 2, 3, 6, 7},
	{1, 2, 4, 5, 6},
	{1, 2, 4, 5, 7},
	{1, 2, 4, 6, 7},
	{1, 2, 5, 6, 7},
	{1, 3, 4, 5, 6},
	{1, 3, 4, 5, 7},
	{1, 3, 4, 6, 7},
	{1, 3, 5, 6, 7},
	{1, 4, 5, 6, 7},
	{2, 3, 4, 5, 6},
	{2, 3, 4, 5, 7},
	{2, 3, 4, 6, 7},
	{2, 3, 5, 6, 7},
	{2, 4, 5, 6, 7},
	{3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5},
	{0, 1, 2, 3, 4, 6},
	{0, 1, 2, 3, 4, 7},
	{0, 1, 2, 3, 5, 6},
	{0, 1, 2, 3, 5, 7},
	{0, 1, 2, 3, 6, 7},
	{0, 1, 2, 4, 5, 6},
	{0, 1, 2, 4, 5, 7},
	{0, 1, 2, 4, 6, 7},
	{0, 1, 2, 5, 6, 7},
	{0, 1, 3, 4, 5, 6},
	{0, 1, 3, 4, 5, 7},
	{0, 1, 3, 4, 6, 7},
	{0, 1, 3, 5, 6, 7},
	{0, 1, 4, 5, 6, 7},
	{0, 2, 3, 4, 5, 6},
	{0, 2, 3, 4, 5, 7},
	{0, 2, 3, 4, 6, 7},
	{0, 2, 3, 5, 6, 7},
	{0, 2, 4, 5, 6, 7},
	{0, 3, 4, 5, 6, 7},
	{1, 2, 3, 4, 5, 6},
	{1, 2, 3, 4, 5, 7},
	{1, 2, 3, 4, 6, 7},
	{1, 2, 3, 5, 6, 7},
	{1, 2, 4, 5, 6, 7},
	{1, 3, 4, 5, 6, 7},
	{2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6},
	{0, 1, 2, 3, 4, 5, 7},
	{0, 1, 2, 3, 4, 6, 7},
	{0, 1, 2, 3, 5, 6, 7},
	{0, 1, 2, 4, 5, 6, 7},
	{0, 1, 3, 4, 5, 6, 7},
	{0, 2, 3, 4, 5, 6, 7},
	{1, 2, 3, 4, 5, 6, 7},
	{0, 1, 2, 3, 4, 5, 6, 7},
}

// 自动调整矩阵
func (qr *QRCode) autoAdjustmentMatrix(qrItem *QRCodeItem, matrix *[]float64) int {
	score := qr.scoreMatrix(qrItem, matrix)

	matrixValue := *matrix
	step := 0.02
	for _, changeIndexList := range autoAdjustmentMatrixChangeIndex {
		for changePercent := -0.5; changePercent < 0.5; changePercent += step {
			newMatrixValue := make([]float64, len(matrixValue))
			copy(newMatrixValue, matrixValue)
			for _, index := range changeIndexList {
				newMatrixValue[index] *= changePercent
			}

			testScore := qr.scoreMatrix(qrItem, &newMatrixValue)
			if testScore > score {
				score = testScore
				*matrix = newMatrixValue
			}
		}
	}
	qrItem.matrix = *matrix
	return score
}

// 给透视矩阵打分
func (qr *QRCode) scoreMatrix(qrItem *QRCodeItem, matrix *[]float64) int {
	score := 0
	// 给定位角打分
	score += qr.scorePositioningAngle(qrItem, matrix)
	// 给定位虚线打分
	score += qr.scoreDashedLine(qrItem, matrix)
	// 给 align pattern 打分
	score += qr.scoreAlignPattern(qrItem, matrix)
	return score
}

type QRScoreArea struct {
	FromX   int
	FromY   int
	TargetX int
	TargetY int
	Color   bool
}

// 给定位角打分
func (qr *QRCode) scorePositioningAngle(qrItem *QRCodeItem, matrix *[]float64) int {
	score := 0
	scoreAreaList := []QRScoreArea{
		// 定位角中心
		{2, qrItem.size - 5, 5, qrItem.size - 2, qrIsBlack},
		{2, 2, 5, 5, qrIsBlack},
		{qrItem.size - 5, 2, qrItem.size - 2, 5, qrIsBlack},
		// 定位角边缘空白
		{0, qrItem.size - 8, 8, qrItem.size - 7, qrIsWhite},
		{7, qrItem.size - 7, 8, qrItem.size, qrIsWhite},
		{qrItem.size - 8, 0, qrItem.size - 7, 8, qrIsWhite},
		{qrItem.size - 7, 7, qrItem.size, 8, qrIsWhite},
		// 定位角边缘黑线
		{0, qrItem.size - 7, 1, qrItem.size, qrIsBlack},
		{6, qrItem.size - 7, 7, qrItem.size, qrIsBlack},
		{1, qrItem.size - 7, 6, qrItem.size - 6, qrIsBlack},
		{1, qrItem.size - 1, 6, qrItem.size, qrIsBlack},
		{qrItem.size - 7, 0, qrItem.size - 6, 7, qrIsBlack},
		{qrItem.size - 1, 0, qrItem.size, 7, qrIsBlack},
		{qrItem.size - 6, 0, qrItem.size - 1, 1, qrIsBlack},
		{qrItem.size - 6, 6, qrItem.size - 1, 7, qrIsBlack},
		// 定位角内部空白
		{1, qrItem.size - 6, 2, qrItem.size - 1, qrIsWhite},
		{5, qrItem.size - 6, 6, qrItem.size - 1, qrIsWhite},
		{2, qrItem.size - 6, 5, qrItem.size - 5, qrIsWhite},
		{2, qrItem.size - 2, 5, qrItem.size - 1, qrIsWhite},
		{qrItem.size - 6, 1, qrItem.size - 1, 2, qrIsWhite},
		{qrItem.size - 6, 5, qrItem.size - 1, 6, qrIsWhite},
		{qrItem.size - 6, 2, qrItem.size - 5, 5, qrIsWhite},
		{qrItem.size - 2, 2, qrItem.size - 1, 5, qrIsWhite},
	}

	for _, areaInfo := range scoreAreaList {
		score += qr.scoreArea(qrItem, matrix, areaInfo.FromX, areaInfo.FromY, areaInfo.TargetX, areaInfo.TargetY, areaInfo.Color)
	}
	return score
}

func (qr *QRCode) scoreDashedLine(qrItem *QRCodeItem, matrix *[]float64) int {
	score := 0
	for i := 7; i < qrItem.size-7; i++ {
		color := qrIsWhite
		if i%2 == 0 {
			color = qrIsBlack
		}
		score += qr.scoreArea(qrItem, matrix, i, 6, i+1, 7, color)
		score += qr.scoreArea(qrItem, matrix, 6, i, 7, i+1, color)
	}
	return score
}

func (qr *QRCode) scoreAlignPattern(qrItem *QRCodeItem, matrix *[]float64) int {
	return 0
}

// 给区域打分
func (qr *QRCode) scoreArea(qrItem *QRCodeItem, matrix *[]float64, fromX, fromY, targetX, targetY int, isBlack bool) int {
	score := 0
	for x := fromX; x < targetX; x++ {
		for y := fromY; y < targetY; y++ {
			point := method.PerspectiveTransform(matrix, float64(x)+0.5, float64(y)+0.5)
			if point.Y < 0 || point.Y >= qr.Height || point.X < 0 || point.X >= qr.Width {
				return 0
			}
			index := point.Y*qr.Width + point.X
			if int(qr.greyPixels[index]) != 0 {
				if !isBlack {
					score++
				} else {
					// score--
				}
			} else {
				if isBlack {
					// score--
				} else {
					score++
				}
			}
		}
	}
	return score
}

func (qr *QRCode) GetResult() img.Value {
	results := []img.Value{}
	for _, item := range qr.codeItems {
		if !item.success {
			continue
		}

		result := img.Value{
			Type: img.ValueTypeMap,
			Values: map[string]img.Value{
				"value": item.result,
				"posi": {
					Type: img.ValueTypeRect,
					Rect: []img.ValuePosition{item.corners[0].corners[3], item.corners[2].corners[1]},
				},
			},
		}
		results = append(results, result)
	}
	return img.Value{
		Type: img.ValueTypeList,
		List: results,
	}
}

func (qrItem *QRCodeItem) decode() error {
	// 获取格式信息
	err := qrItem.getFormatData()
	if err != nil {
		return err
	}
	// https://www.jianshu.com/p/3cf1862552f8

	qrItem.readData()
	return nil
}

func (qrItem *QRCodeItem) getFormatData() error {
	format := make([]int, 15)
	formatMask := []int{1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0}
	for x := 0; x < 8; x++ {
		index := 8*qrItem.size + (qrItem.size - x - 1)
		format[x] = qrItem.Pixels[index]
	}
	for y := 0; y < 7; y++ {
		index := (qrItem.size-7+y)*qrItem.size + 8
		format[y+8] = qrItem.Pixels[index]
	}
	format = method.XOR(method.ReverseArray(format), formatMask)
	formatInt := method.BinaryToInt(format)
	decodeFormatRes, newFormat := method.BCH_Decode_Format(formatInt)
	if decodeFormatRes == -1 {
		ys := [15]int{0, 1, 2, 3, 4, 5, 7, 8, 8, 8, 8, 8, 8, 8, 8}
		xs := [15]int{8, 8, 8, 8, 8, 8, 8, 8, 7, 5, 4, 3, 2, 1, 0}
		for i := 0; i < 15; i++ {
			index := ys[i]*qrItem.size + xs[i]
			format[i] = qrItem.Pixels[index]
		}
		format = method.XOR(method.ReverseArray(format), formatMask)
		formatInt := method.BinaryToInt(format)
		decodeFormatRes, newFormat2 := method.BCH_Decode_Format(formatInt)
		if decodeFormatRes == -1 {
			return errors.New("get format info error")
		}
		newFormat = newFormat2
	}
	format = method.IntToBinary(newFormat, 15)

	qrItem.errorCorrectionLevel = method.BinaryToInt(format[0:2])
	qrItem.mask = method.BinaryToInt(format[2:5])
	return nil
}

func (qrItem *QRCodeItem) readData() error {
	// 从右向左，从低向上
	y := qrItem.size - 1
	x := qrItem.size - 1
	step := -1
	pixelCount := 0
	eccInfo := qrEccList[qrItem.version]
	errLevel := 0
	switch qrItem.errorCorrectionLevel {
	case 1:
		errLevel = 0
	case 0:
		errLevel = 1
	case 3:
		errLevel = 2
	case 2:
		errLevel = 3
	}
	eccItem := eccInfo[errLevel]
	parts := make([][]int, eccItem.tdc+eccItem.ecc*(eccItem.blocks1+eccItem.blocks2))
foreachPixel:
	for x > 0 {
		if x == 6 {
			x--
		}

		xs := []int{x, x - 1}
		for _, xItem := range xs {
			pixelIndex := y*qrItem.size + xItem
			if qrItem.checkNotIsData(xItem, y) {
				continue
			}
			partIndex := pixelCount / 8
			if partIndex >= len(parts) {
				break foreachPixel
			}
			pixelCount++
			bit := qrItem.Pixels[pixelIndex]
			if qrItem.dataMask(xItem, y) > 0 {
				bit ^= 1
			}

			parts[partIndex] = append(parts[partIndex], bit)
		}

		y += step
		if y < 0 || y >= qrItem.size {
			step = -step
			x -= 2
			y += step
		}
	}

	allBlocks := eccItem.blocks1 + eccItem.blocks2
	blocksData := make([]int, 0)

	for blockIndex := 0; blockIndex < allBlocks; blockIndex++ {

		dc := eccItem.dc1
		if blockIndex >= eccItem.blocks1 {
			dc = eccItem.dc2
		}
		blockLength := dc + eccItem.ecc
		blockData := make([]int, blockLength)

		for bitIndex := 0; bitIndex < dc; bitIndex++ {
			// 因为是多个 block 横向排列，然后纵向读取之后再放入数据中的
			valueInBitIndex := blockIndex + bitIndex*allBlocks
			blockData[bitIndex] = method.BinaryToInt(parts[valueInBitIndex])
		}

		for eccBitIndex := 0; eccBitIndex < eccItem.ecc; eccBitIndex++ {
			// 因为是多个 block 横向排列，然后纵向读取之后再放入数据中的
			eccInBitIndex := eccItem.tdc + blockIndex + eccBitIndex*allBlocks
			blockData[dc+eccBitIndex] = method.BinaryToInt(parts[eccInBitIndex])
		}
		res, err := method.RS_Error_Correct(blockData, eccItem.ecc, []int{}, method.Galois_GF_256)
		if err != nil {
			fmt.Println("err", err)
			return err
		}
		dataRes := res[:dc]
		dataBit := make([]int, 0)
		for _, number := range dataRes {
			// 再转换回二进制
			dataBit = method.ConcatArray(dataBit, method.IntToBinary(number, 8))
		}
		blocksData = method.ConcatArray(blocksData, dataBit)
	}

	qrItem.blocksData = blocksData

	var err error
dataTypeDecode:
	for qrItem.currentReadIndex < len(qrItem.blocksData)-4 {
		dataType := method.BinaryToInt(blocksData[qrItem.currentReadIndex : qrItem.currentReadIndex+4])
		qrItem.currentReadIndex += 4
		switch dataType {
		case qrDataTypeEnd:
			break dataTypeDecode
		case qrDataTypeNumric:
			err = qrItem.readNumricData()
		case qrDataType8BitByte:
			err = qrItem.read8BitByteData()
		case qrDataTypeECI:
			err = qrItem.readECIData()
		default:
			return errors.New("this data type is not currently supported")
		}
		if err != nil {
			return err
		}
	}
	qrItem.success = true

	return nil
}

// 检测是否不是数据块
func (qrItem *QRCodeItem) checkNotIsData(x, y int) bool {
	finderSize := 9

	if x == 6 || y == 6 {
		return true
	}

	// 三个角的定位点
	if x < finderSize && y < finderSize || x > qrItem.size-finderSize && y < finderSize || x < finderSize && y > qrItem.size-finderSize {
		return true
	}
	// 版本信息
	if qrItem.version >= 7 {
		if x < 6 && y > qrItem.size-finderSize-3 {
			return true
		}

		if x > qrItem.size-finderSize-3 && y < 6 {
			return true
		}
	}

	// 版本信息区域

	// Alignment Patterns
	qrDataInfo := QRDataInfoList[qrItem.version]
	apLastIndex := len(qrDataInfo.alignmentPatterns) - 1
	for xi, alignPattern := range qrDataInfo.alignmentPatterns {
		for yi, alignPatternY := range qrDataInfo.alignmentPatterns {
			if (math.Abs(float64(alignPattern-x))) < 3 && (math.Abs(float64(alignPatternY-y))) < 3 {
				if (xi == 0 && yi == 0) || (xi == apLastIndex && yi == 0) || (yi == apLastIndex && xi == 0) {
					return false
				}
				return true
			}
		}
	}

	return false
}

// 数据掩码
func (qrItem *QRCodeItem) dataMask(j, i int) int {
	k := 0
	switch qrItem.mask {
	case 0:
		k = (i + j) % 2
	case 1:
		k = i % 2
	case 2:
		k = j % 3
	case 3:
		k = (i + j) % 3
	case 4:
		k = ((i / 2) + (j / 3)) % 2
	case 5:
		k = (i*j)%2 + (i*j)%3
	case 6:
		k = ((i*j)%2 + (i*j)%3) % 2
	case 7:
		k = ((i*j)%3 + (i+j)%2) % 2
	default:
		return 0
	}
	if k != 0 {
		return 0
	}
	return 1
}

func (qrItem *QRCodeItem) readNumricData() error {
	dataLengthBitLen := 10
	if qrItem.version >= 10 && qrItem.version <= 26 {
		dataLengthBitLen = 12
	} else if qrItem.version > 26 {
		dataLengthBitLen = 14
	}
	dataStartIndex := qrItem.currentReadIndex + dataLengthBitLen
	length := method.BinaryToInt(qrItem.blocksData[qrItem.currentReadIndex:dataStartIndex])
	part := length / 3
	remainder := length % 3

	result := []int{}
	for i := 0; i < part; i++ {
		data := method.BinaryToInt(qrItem.blocksData[dataStartIndex : dataStartIndex+10])
		result = append(result, (method.IntToIntList(data))...)
		dataStartIndex += 10
	}
	if remainder == 2 {
		data := method.BinaryToInt(qrItem.blocksData[dataStartIndex : dataStartIndex+7])
		result = append(result, (method.IntToIntList(data))...)
		dataStartIndex += 7
	} else if remainder == 1 {
		data := method.BinaryToInt(qrItem.blocksData[dataStartIndex : dataStartIndex+4])
		result = append(result, (method.IntToIntList(data))...)
		dataStartIndex += 4
	}
	qrItem.currentReadIndex = dataStartIndex
	qrItem.result = img.Value{
		Type:   img.ValueTypeString,
		String: qrItem.result.String + method.IntListToString(result),
	}
	return nil
}

func (qrItem *QRCodeItem) read8BitByteData() error {
	dataLengthBitLen := 8
	if qrItem.version >= 10 {
		dataLengthBitLen = 16
	}
	dataStartIndex := qrItem.currentReadIndex + dataLengthBitLen
	length := method.BinaryToInt(qrItem.blocksData[qrItem.currentReadIndex:dataStartIndex])

	result := []byte{}
	for i := 0; i < length; i++ {
		data := method.BinaryToInt(qrItem.blocksData[dataStartIndex : dataStartIndex+8])
		result = append(result, byte(data))
		dataStartIndex += 8
	}
	qrItem.currentReadIndex = dataStartIndex
	// UTF-8BOM 头
	if result[0] == byte(239) && result[1] == byte(187) && result[2] == byte(191) {
		result = result[3:]
	}
	qrItem.result = img.Value{
		Type:   img.ValueTypeString,
		String: qrItem.result.String + string(result),
	}
	return nil
}

func (qrItem *QRCodeItem) readECIData() error {
	dataStartIndex := qrItem.currentReadIndex + 8
	indicator := method.BinaryToInt(qrItem.blocksData[qrItem.currentReadIndex:dataStartIndex])

	if indicator != 26 {
		return errors.New("only support utf-8 0001 1010")
	}
	qrItem.currentReadIndex = dataStartIndex
	return nil
}
