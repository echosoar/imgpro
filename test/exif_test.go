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

func TestExifJpgGPS(t *testing.T) {
	result := Run("./imgs/exif-gps.jpeg", []string{"exif"})

	if result["exif"].Values["GPSLatitude"].String != "42 45 3985/1000" {
		jsonData, _ := json.Marshal(result)
		t.Fatal("exif GPSLatitude err", string(jsonData))
	}

	if result["exif"].Values["GPSLongitude"].String != "84 29 28508/1000" {
		jsonData, _ := json.Marshal(result)
		t.Fatal("exif GPSLatitude err", string(jsonData))
	}
}
