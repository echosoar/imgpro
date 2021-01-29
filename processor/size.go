package processor

import (
	"os"

	img "github.com/echosoar/imgpro/core"
)

// SizeProcessor bin size processor
func SizeProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:   []string{"size"},
		Runner: sizeRunner,
	})
}

func sizeRunner(core *img.Core) map[string]img.Value {
	file, err := os.Stat(core.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			panic("File \"" + core.FilePath + "\" not exists")
		}
		panic(err)
	}
	// get the size
	size := file.Size()
	return map[string]img.Value{
		"size": {
			Type: img.ValueTypeInt,
			Int:  int(size),
		},
	}
}
