package processor

import (
	img "github.com/echosoar/imgpro/core"
)

// PositionProcessor device
func PositionProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:          []string{"position"},
		PreConditions: []string{"exif"},
		Runner:        positionRunner,
	})
}

func positionRunner(core *img.Core) map[string]img.Value {
	position := ""

	exif := core.Result["exif"].Values

	latitude, latitudeExists := exif["GPSLatitude"]
	longitude, longitudeExists := exif["GPSLongitude"]

	if latitudeExists && longitudeExists {
		position = latitude.String + " " + exif["GPSLatitudeRef"].String +  "," + longitude.String + " " + exif["GPSLongitudeRef"].String
	}

	return map[string]img.Value{
		"position": {
			Type:   img.ValueTypeString,
			String: position,
		},
	}
}
