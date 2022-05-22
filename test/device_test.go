package test

import (
	"encoding/json"
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestDeviceJpg(t *testing.T) {
	result := Run("./imgs/exif.jpg", []string{"device"})

	if result["device"].String != "iPhone 8 Plus" {
		jsonData, _ := json.Marshal(result)
		t.Fatal("device err", string(jsonData))
	}
}