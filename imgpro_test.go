package imgpro

import (
	"testing"
)

func TestSize(t *testing.T) {
	result := Run("./test/imgs/go.png", []string{"size"})
	if result["size"].Int != 60746 {
		t.Fatal("size error")
	}
}

func TestTypePng(t *testing.T) {
	result := Run("./test/imgs/go.png", []string{"type"})
	if result["type"].String != "png" {
		t.Fatal("type error")
	}
}
func TestTypeJpeg(t *testing.T) {
	result := Run("./test/imgs/lincoln.jpg", []string{"type"})
	if result["type"].String != "jpeg" {
		t.Fatal("type error")
	}
}

func TestTypeBmp(t *testing.T) {
	result := Run("./test/imgs/iojs.bmp", []string{"type"})
	if result["type"].String != "bmp" {
		t.Fatal("type error")
	}
}

func TestTypeGif(t *testing.T) {
	result := Run("./test/imgs/cool.gif", []string{"type"})
	if result["type"].String != "gif" {
		t.Fatal("type error")
	}
}

func TestTypeWebp(t *testing.T) {
	result := Run("./test/imgs/cool.webp", []string{"type"})
	if result["type"].String != "webp" {
		t.Fatal("type error")
	}
}
