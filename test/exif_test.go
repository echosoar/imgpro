package test

import (
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestExifJpg(t *testing.T) {
	result := Run("./imgs/exif.jpg", []string{"exif"})
	t.Fatal(result)
}
