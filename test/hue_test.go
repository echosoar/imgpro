package test

import (
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestHUEPng(t *testing.T) {
	result := Run("./imgs/go.png", []string{"hue"})
	if len(result["hue"].Rgba[0]) != 2 {
		t.Fatal("hue error")
	}

	if result["hue"].Rgba[0][0].R != 255 || result["hue"].Rgba[0][0].G != 255 || result["hue"].Rgba[0][0].B != 255 || result["hue"].Rgba[0][0].A != 255 {
		t.Fatal("hue error")
	}

	if result["hue"].Rgba[0][1].R != 106 || result["hue"].Rgba[0][1].G != 215 || result["hue"].Rgba[0][1].B != 229 || result["hue"].Rgba[0][1].A != 255 {
		t.Fatal("hue error")
	}
}

func TestHUEJpg(t *testing.T) {
	result := Run("./imgs/go.jpg", []string{"hue"})
	if len(result["hue"].Rgba[0]) != 3 {
		t.Fatal("hue jpeg err", len(result["hue"].Rgba[0]))
	}
}
