package processor

import (
	"bufio"
	"image/gif"
	"os"

	img "github.com/echosoar/imgpro/core"
)

// FrameProcessor get image frame count
func FrameProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:         []string{"frame"},
		Precondition: []string{"type"},
		Runner:       frameRunner,
	})
}

func frameRunner(core *img.Core) map[string]img.Value {
	frame := 0
	imgType := core.Result["type"].String
	if imgType == "jpg" || imgType == "png" || imgType == "bmp" {
		frame = 1
	} else if imgType == "gif" {
		fp, err := os.Open(core.FilePath)
		if err != nil {
			panic(err)
		}
		defer fp.Close()
		reader := bufio.NewReader(fp)
		gif, err := gif.DecodeAll(reader)
		if err != nil {
			panic(err)
		}
		frame = len(gif.Image)
	} else if imgType == "webp" {

	}
	return map[string]img.Value{
		"frame": {
			Type: img.ValueTypeInt,
			Int:  frame,
		},
	}
}
