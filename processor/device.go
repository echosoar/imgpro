package processor

import (
	img "github.com/echosoar/imgpro/core"
)

// DeviceProcessor device
func DeviceProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:          []string{"device"},
		PreConditions: []string{"exif"},
		Runner:        deviceRunner,
	})
}

func deviceRunner(core *img.Core) map[string]img.Value {
	device := ""

	exif := core.Result["exif"].Values

	attributes := []string{"Model", "Software"}
	for _, attr := range attributes {
		model, exists := exif[attr]
		if exists {
			device = model.String
			break
		}
	}

	return map[string]img.Value{
		"device": {
			Type:   img.ValueTypeString,
			String: device,
		},
	}
}
