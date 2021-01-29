package processor

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"image"

	// Register to image
	_ "image/jpeg"
	"os"

	img "github.com/echosoar/imgpro/core"
)

// WHProcessor get image width and height
func WHProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:         []string{"width", "height"},
		Precondition: []string{"type"},
		Runner:       whRunner,
	})
}

func whRunner(core *img.Core) map[string]img.Value {
	imgType := core.Result["type"].String
	width := 0
	height := 0

	f, err := os.Open(core.FilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	if imgType == "png" {
		fileBytes := make([]byte, 24)
		count, readErr := reader.Read(fileBytes)
		if readErr == nil {
			if count == 24 && bytes.HasPrefix(fileBytes[12:16], []byte("IHDR")) {
				width = int(binary.BigEndian.Uint32(fileBytes[16:20]))
				height = int(binary.BigEndian.Uint32(fileBytes[20:24]))
			} else if count >= 16 {
				width = int(binary.BigEndian.Uint32(fileBytes[8:12]))
				height = int(binary.BigEndian.Uint32(fileBytes[12:16]))
			}
		}
	} else if imgType == "jpg" {
		image, _, err := image.DecodeConfig(f)
		if err != nil {
			panic(err)
		}
		width = image.Width
		height = image.Height
	} else if imgType == "gif" {
		fileBytes := make([]byte, 10)
		_, readErr := reader.Read(fileBytes)
		if readErr == nil {
			width = int(binary.LittleEndian.Uint16(fileBytes[6:8]))
			height = int(binary.LittleEndian.Uint16(fileBytes[8:10]))
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
