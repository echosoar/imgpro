package test

import (
	"encoding/json"
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestQRCodeJpg(t *testing.T) {
	result := Run("./imgs/qrcode/1.png", []string{"qrcode"})

	if result["qrcode"].Int < 100 {
		jsonData, _ := json.Marshal(result)
		t.Fatal("qrcode err", string(jsonData))
	}
}
