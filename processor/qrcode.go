package processor

import (
	"fmt"

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
	corners    []*QRCodeRegion
}

type QRCodeRegion struct {
	posi         img.ValuePosition
	size         int
	cornersIndex int
}

func (qr *QRCode) Run() {
	qr.grayscale()
	qr.binarization()
	qr.findCorners()
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
	qrMaxRegion       = 20
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
	cornerRightLine := qr.regionCode(line, row-pb[4])
	cornerCenter := qr.regionCode(line, row-pb[4]-pb[3]-pb[2])
	cornerLeftLine := qr.regionCode(line, row-pb[4]-pb[3]-pb[2]-pb[1]-pb[0])

	if cornerRightLine < 0 || cornerCenter < 0 || cornerLeftLine < 0 || cornerLeftLine != cornerRightLine || cornerLeftLine == cornerCenter {
		return
	}

	centerRegion := qr.getRegionByRegionIndex(cornerCenter)
	lineRegion := qr.getRegionByRegionIndex(cornerRightLine)

	// 已经放在某一个 corner 之中了
	if centerRegion.cornersIndex >= 0 || lineRegion.cornersIndex >= 0 {
		return
	}

	ratio := centerRegion.size * 100 / lineRegion.size
	if ratio < 10 || ratio > 70 {
		return
	}
	qr.newCorner(lineRegion, centerRegion)
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

func (qr *QRCode) newCorner(lineRegion *QRCodeRegion, centerRegion *QRCodeRegion) {
	lineRegion.cornersIndex = 1
	centerRegion.cornersIndex = 1
	fmt.Println("new corner", lineRegion, centerRegion)
}
