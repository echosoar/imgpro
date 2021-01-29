package processor

import (
	"image"

	// Register to image
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
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
	if imgType == "png" || imgType == "jpg" || imgType == "gif" {
		image, _, err := image.DecodeConfig(f)
		if err != nil {
			panic(err)
		}
		width = image.Width
		height = image.Height
	} else if imgType == "bmp" {
		// refs: https://github.com/scardine/image_size/blob/master/get_image_size.py
		// fileBytes := make([]byte, 26)
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
