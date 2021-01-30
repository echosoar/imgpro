package test

import (
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestWHPng(t *testing.T) {
	result := Run("./imgs/go.png", []string{"width", "height"})
	if result["width"].Int != 1576 {
		t.Fatal("png width error")
	}
	if result["height"].Int != 890 {
		t.Fatal("png height error")
	}
}

func TestWHJpg(t *testing.T) {
	result := Run("./imgs/lincoln.jpg", []string{"width", "height"})
	if result["width"].Int != 645 {
		t.Fatal("jpg width error")
	}
	if result["height"].Int != 796 {
		t.Fatal("jpg height error")
	}
}

func TestWHGif(t *testing.T) {
	result := Run("./imgs/cool.gif", []string{"width", "height"})
	if result["width"].Int != 300 {
		t.Fatal("gif width error")
	}
	if result["height"].Int != 300 {
		t.Fatal("gif height error")
	}
}

func TestWHBmp(t *testing.T) {
	result := Run("./imgs/iojs.bmp", []string{"width", "height"})
	if result["width"].Int != 367 {
		t.Fatal("bmp width error")
	}
	if result["height"].Int != 187 {
		t.Fatal("bmp height error")
	}
}

func TestWHWebp(t *testing.T) {
	result := Run("./imgs/cool.webp", []string{"width", "height"})
	if result["width"].Int != 300 {
		t.Fatal("webp width error")
	}
	if result["height"].Int != 300 {
		t.Fatal("webp height error")
	}
}
