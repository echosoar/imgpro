package test

import (
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestExifJpg(t *testing.T) {
	result := Run("./imgs/exif.jpg", []string{"exif"})
	if result["exif"].Values["ModifyDate"].String != "2021:02:16 16:56:06" {
		t.Fatal("exif ModifyDate error", len(result["exif"].Values["ModifyDate"].String))
	}
}
