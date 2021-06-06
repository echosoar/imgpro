package test

import (
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestRGBAPng(t *testing.T) {
	result := Run("./imgs/go.png", []string{"rgba"})
	if len(result["rgba"].Rgba[0]) != 1402640 {
		t.Fatal("rgba png error")
	}
	if result["rgba"].Rgba[0][600000].R != 106 {
		t.Fatal("rgba png r error")
	}
	if result["rgba"].Rgba[0][600000].G != 215 {
		t.Fatal("rgba png g error")
	}
	if result["rgba"].Rgba[0][600000].B != 229 {
		t.Fatal("rgba png b error")
	}
}

func TestRGBAJpg(t *testing.T) {
	result := Run("./imgs/go.jpg", []string{"rgba"})
	if len(result["rgba"].Rgba[0]) != 1402640 {
		t.Fatal("rgba jpeg error")
	}
	if result["rgba"].Rgba[0][600000].R != 107 {
		t.Fatal("rgba jpeg r error")
	}
	if result["rgba"].Rgba[0][600000].G != 215 {
		t.Fatal("rgba jpeg g error")
	}
	if result["rgba"].Rgba[0][600000].B != 229 {
		t.Fatal("rgba jpeg b error")
	}
}
