package test

import (
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestTypePng(t *testing.T) {
	result := Run("./imgs/go.png", []string{"type"})
	if result["type"].String != "png" {
		t.Fatal("type error")
	}
}
func TestTypeJpg(t *testing.T) {
	result := Run("./imgs/lincoln.jpg", []string{"type"})
	if result["type"].String != "jpg" {
		t.Fatal("type error")
	}
}

func TestTypeBmp(t *testing.T) {
	result := Run("./imgs/iojs.bmp", []string{"type"})
	if result["type"].String != "bmp" {
		t.Fatal("type error")
	}
}

func TestTypeGif(t *testing.T) {
	result := Run("./imgs/cool.gif", []string{"type"})
	if result["type"].String != "gif" {
		t.Fatal("type error")
	}
}

func TestTypeWebp(t *testing.T) {
	result := Run("./imgs/cool.webp", []string{"type"})
	if result["type"].String != "webp" {
		t.Fatal("type error")
	}
}
