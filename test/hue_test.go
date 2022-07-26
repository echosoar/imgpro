package test

import (
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestHUEPng(t *testing.T) {
	result := Run("./imgs/go.png", []string{"hue"})
	if len(result["hue"].Frames[0].Rgba) != 3 {
		t.Fatal("hue error", result["hue"].Frames[0].Rgba)
	}

	if result["hue"].Frames[0].Rgba[0].R != 254 || result["hue"].Frames[0].Rgba[0].G != 254 || result["hue"].Frames[0].Rgba[0].B != 254 || result["hue"].Frames[0].Rgba[0].A != 255 {
		t.Fatal("hue 0 error", result["hue"])
	}

	if result["hue"].Frames[0].Rgba[1].R != 105 || result["hue"].Frames[0].Rgba[1].G != 214 || result["hue"].Frames[0].Rgba[1].B != 228 || result["hue"].Frames[0].Rgba[1].A != 255 {
		t.Fatal("hue 1 error")
	}

	if result["hue"].Frames[0].Rgba[2].R != 2 || result["hue"].Frames[0].Rgba[2].G != 2 || result["hue"].Frames[0].Rgba[2].B != 2 || result["hue"].Frames[0].Rgba[1].A != 255 {
		t.Fatal("hue 2 error")
	}
}

func TestHUEJpg(t *testing.T) {
	result := Run("./imgs/exif.jpg", []string{"hue"})
	if len(result["hue"].Frames[0].Rgba) != 4 {
		t.Fatal("hue jpeg err", result["hue"].Frames[0].Rgba)
	}
}
