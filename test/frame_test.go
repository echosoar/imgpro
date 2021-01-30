package test

import (
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestFramePng(t *testing.T) {
	result := Run("./imgs/go.png", []string{"frame"})
	if result["frame"].Int != 1 {
		t.Fatal("frame error")
	}
}
func TestFrameJpg(t *testing.T) {
	result := Run("./imgs/go.jpg", []string{"frame"})
	if result["frame"].Int != 1 {
		t.Fatal("frame error")
	}
}

func TestFrameBmp(t *testing.T) {
	result := Run("./imgs/go.bmp", []string{"frame"})
	if result["frame"].Int != 1 {
		t.Fatal("frame error")
	}
}
