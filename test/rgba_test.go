package test

import (
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestRGBAPng(t *testing.T) {
	result := Run("./imgs/go.png", []string{"rgba"})
	if len(result["rgba"].Frames[0].Rgba) != 1402640 {
		t.Fatal("rgba png error")
	}
	if result["rgba"].Frames[0].Rgba[600000].R != 106 {
		t.Fatal("rgba png r error")
	}
	if result["rgba"].Frames[0].Rgba[600000].G != 215 {
		t.Fatal("rgba png g error")
	}
	if result["rgba"].Frames[0].Rgba[600000].B != 229 {
		t.Fatal("rgba png b error")
	}
}

func TestRGBAJpg(t *testing.T) {
	result := Run("./imgs/go.jpg", []string{"rgba"})
	if len(result["rgba"].Frames[0].Rgba) != 1402640 {
		t.Fatal("rgba jpeg error")
	}
	if result["rgba"].Frames[0].Rgba[600000].R != 107 {
		t.Fatal("rgba jpeg r error")
	}
	if result["rgba"].Frames[0].Rgba[600000].G != 215 {
		t.Fatal("rgba jpeg g error")
	}
	if result["rgba"].Frames[0].Rgba[600000].B != 229 {
		t.Fatal("rgba jpeg b error")
	}
}

func TestRGBAGif(t *testing.T) {
	result := Run("./imgs/cool.gif", []string{"rgba", "frame"})
	if result["frame"].Int != 13 {
		t.Fatal("rgba gif frames error", result["frame"].Int)
	}
}

func TestRGBABmp(t *testing.T) {
	result := Run("./imgs/go_24.bmp", []string{"rgba"})
	if len(result["rgba"].Frames[0].Rgba) != 1402640 {
		t.Fatal("rgba bmp error")
	}
	if result["rgba"].Frames[0].Rgba[600000].R != 106 {
		t.Fatal("rgba bmp r error")
	}
	if result["rgba"].Frames[0].Rgba[600000].G != 215 {
		t.Fatal("rgba bmp g error")
	}
	if result["rgba"].Frames[0].Rgba[600000].B != 229 {
		t.Fatal("rgba bmp b error")
	}
}

func TestRGBAWebp(t *testing.T) {
	result := Run("./imgs/go_32.webp", []string{"rgba"})
	if len(result["rgba"].Frames[0].Rgba) != 1402640 {
		t.Fatal("rgba webp error")
	}
	if result["rgba"].Frames[0].Rgba[600000].R != 107 {
		t.Fatal("rgba webp r error")
	}
	if result["rgba"].Frames[0].Rgba[600000].G != 201 {
		t.Fatal("rgba webp g error")
	}
	if result["rgba"].Frames[0].Rgba[600000].B != 213 {
		t.Fatal("rgba webp b error")
	}
}

func TestRGBAApng(t *testing.T) {
	result := Run("./imgs/cool.apng", []string{"rgba", "frame"})
	if result["frame"].Int != 3 {
		t.Fatal("rgba apng frames error", result["frame"].Int)
	}
	if len(result["rgba"].Frames[0].Rgba) != 16 {
		t.Fatal("rgba apng frame size error")
	}
	// Frame 0 is solid red
	if result["rgba"].Frames[0].Rgba[0].R != 255 {
		t.Fatal("rgba apng frame0 r error")
	}
	if result["rgba"].Frames[0].Rgba[0].G != 0 {
		t.Fatal("rgba apng frame0 g error")
	}
	if result["rgba"].Frames[0].Rgba[0].B != 0 {
		t.Fatal("rgba apng frame0 b error")
	}
	// Frame 1 is solid green
	if result["rgba"].Frames[1].Rgba[0].G != 255 {
		t.Fatal("rgba apng frame1 g error")
	}
	// Frame 2 is solid blue
	if result["rgba"].Frames[2].Rgba[0].B != 255 {
		t.Fatal("rgba apng frame2 b error")
	}
}
