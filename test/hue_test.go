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

	if result["hue"].Rgba[0][0].R != 106 || result["hue"].Rgba[0][0].G != 215 || result["hue"].Rgba[0][0].B != 229 || result["hue"].Rgba[0][0].A != 255 {
		t.Fatal("hue 0 error")
	}

	if result["hue"].Rgba[0][1].R != 0 || result["hue"].Rgba[0][1].G != 0 || result["hue"].Rgba[0][1].B != 0 || result["hue"].Rgba[0][1].A != 255 {
		t.Fatal("hue 1 error")
	}

	if result["hue"].Rgba[0][2].R != 246 || result["hue"].Rgba[0][2].G != 210 || result["hue"].Rgba[0][2].B != 162 || result["hue"].Rgba[0][1].A != 255 {
		t.Fatal("hue 2 error")
	}
}

func TestHUEJpg(t *testing.T) {
	result := Run("./imgs/go.jpg", []string{"hue"})
	if len(result["hue"].Rgba[0]) != 3 {
		t.Fatal("hue jpeg err", len(result["hue"].Rgba[0]))
	}
}
