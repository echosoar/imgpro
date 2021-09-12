package test

import (
	"encoding/json"
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestExifJpg(t *testing.T) {
	result := Run("./imgs/exif.jpg", []string{"exif"})

	if result["exif"].Values["ModifyDate"].String != "2021:02:16 16:56:06" {
		jsonData, _ := json.Marshal(result)
		t.Fatal("exif ModifyDate err", string(jsonData))
	}
}
