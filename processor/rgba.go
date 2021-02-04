package processor

import (
	"bufio"
	"fmt"
	"image"
	"os"

	img "github.com/echosoar/imgpro/core"
)

// RGBAProcessor bin size processor
func RGBAProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:         []string{"rgba", "frame"},
		Precondition: []string{"type"},
		Runner:       rgbaRunner,
	})
}

func rgbaRunner(core *img.Core) map[string]img.Value {
	imgType := core.Result["type"].String
	f, err := os.Open(core.FilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)

	originalImage, _, err := image.Decode(reader)
	if imgType == "png" {

	}
	r, g, b, a := originalImage.At(0, 0).RGBA()
	fmt.Print(r, g, b, a)
	return map[string]img.Value{}
}
