package processor

import (
	img "github.com/echosoar/imgpro/core"
)

// QRCode device
func QRCodeProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:          []string{"qrcode"},
		PreConditions: []string{"rgba"},
		Runner:        qrCodeRunner,
	})
}

func qrCodeRunner(core *img.Core) map[string]img.Value {
	item := make(map[string]img.Value)
	item["Content"] = img.Value{
		Type:   img.ValueTypeString,
		String: "test",
	}
	item["Rect"] = img.Value{
		Type: img.ValueTypeRect,
		Rect: []img.ValuePosition{
			{
				X: 0,
				Y: 0,
			},
			{
				X: 0,
				Y: 0,
			},
			{
				X: 0,
				Y: 0,
			},
			{
				X: 0,
				Y: 0,
			},
		},
	}
	return map[string]img.Value{
		"qrcode": {
			Type: img.ValueTypeList,
			List: []img.Value{
				{
					Values: item,
				},
			},
		},
	}
}
