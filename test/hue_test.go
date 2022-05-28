package test

import (
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestHUEPng(t *testing.T) {
	result := Run("./imgs/go.png", []string{"hue"})
	if len(result["hue"].Rgba[0]) != 3 {
		t.Fatal("hue error", result["hue"].Rgba[0])
	}

	if result["hue"].Rgba[0][0].R != 254 || result["hue"].Rgba[0][0].G != 254 || result["hue"].Rgba[0][0].B != 254 || result["hue"].Rgba[0][0].A != 255 {
		t.Fatal("hue 0 error", result["hue"])
	}

	if result["hue"].Rgba[0][1].R != 106 || result["hue"].Rgba[0][1].G != 215 || result["hue"].Rgba[0][1].B != 229 || result["hue"].Rgba[0][1].A != 255 {
		t.Fatal("hue 1 error")
	}

	if result["hue"].Rgba[0][2].R != 5 || result["hue"].Rgba[0][2].G != 5 || result["hue"].Rgba[0][2].B != 5 || result["hue"].Rgba[0][1].A != 255 {
		t.Fatal("hue 2 error")
	}
}

func TestHUEJpg(t *testing.T) {
	result := Run("./imgs/exif.jpg", []string{"hue"})
	if len(result["hue"].Rgba[0]) != 6 {
		t.Fatal("hue jpeg err", result["hue"].Rgba[0])
	}
}
