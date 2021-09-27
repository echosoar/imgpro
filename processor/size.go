package processor

import (
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
	size := len(core.FileBinary)
	return map[string]img.Value{
		"size": {
			Type: img.ValueTypeInt,
			Int:  int(size),
		},
	}
}
