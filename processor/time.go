package processor

import (
	img "github.com/echosoar/imgpro/core"
)

// TimeProcessor time
func TimeProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:          []string{"time"},
		PreConditions: []string{"exif"},
		Runner:        timeRunner,
	})
}

func timeRunner(core *img.Core) map[string]img.Value {
	time := ""

	exif := core.Result["exif"].Values

	modifyDate, exists := exif["ModifyDate"]

	if exists {
		time = modifyDate.String
	}

	return map[string]img.Value{
		"time": {
			Type:   img.ValueTypeString,
			String: time,
		},
	}
}
