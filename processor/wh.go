package processor

import (
	"bytes"
	"encoding/binary"
	"image"

	// Register to image
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	img "github.com/echosoar/imgpro/core"
)

// WHProcessor get image width and height
func WHProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:          []string{"width", "height", "wh"},
		PreConditions: []string{"type", "size"},
		Runner:        whRunner,
	})
}

func whRunner(core *img.Core) map[string]img.Value {
	imgType := core.Result["type"].String
	size := core.Result["size"].Int
	width := 0
	height := 0
	fileBytes := core.FileBinary

	if imgType == "png" || imgType == "jpg" || imgType == "gif" {
		image, _, err := image.DecodeConfig(bytes.NewReader(core.FileBinary))
		if err != nil {
			panic(err)
		}
		width = image.Width
		height = image.Height
	} else if imgType == "bmp" {
		if size >= 18 {
			headersize := binary.LittleEndian.Uint32(fileBytes[14:18])
			if headersize >= 40 && size >= 26 {
				width = int(binary.LittleEndian.Uint32(fileBytes[18:22]))
				height = int(binary.LittleEndian.Uint32(fileBytes[22:26]))
			} else if headersize == 12 && size >= 22 {
				width = int(binary.LittleEndian.Uint16(fileBytes[18:20]))
				height = int(binary.LittleEndian.Uint16(fileBytes[20:22]))
			}
		}
	} else if imgType == "webp" {
		if size >= 15 {

			switch fileBytes[15] {
			case 32:
				width = int(binary.LittleEndian.Uint16(fileBytes[26:28])) & 0x3fff
				height = int(binary.LittleEndian.Uint16(fileBytes[28:30])) & 0x3fff
			case 76:
				// 后 14 位 为 width
				// 前 14 位 为 height
				bit := int(binary.LittleEndian.Uint32(fileBytes[21:25]))
				width = bit&0x3fff + 1
				height = bit>>14&0x3fff + 1
			case 88:
				width = int(fileBytes[26])<<16 | int(fileBytes[25])<<8 | int(fileBytes[24]) + 1
				height = int(fileBytes[29])<<16 | int(fileBytes[28])<<8 | int(fileBytes[27]) + 1
			}
		}
	}
	return map[string]img.Value{
		"width": {
			Type: img.ValueTypeInt,
			Int:  width,
		},
		"height": {
			Type: img.ValueTypeInt,
			Int:  height,
		},
	}
}
