package test

import (
	"encoding/json"
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestQRCodeNumeric(t *testing.T) {
	result := Run("./imgs/qrcode/1234567.png", []string{"qrcode"})

	if result["qrcode"].Frames[0].List[0].Values["value"].Int != 1234567 {
		jsonData, _ := json.Marshal(result)
		t.Fatal("qrcode numeric err", string(jsonData))
	}
}

func TestQRCodeByte(t *testing.T) {
	result := Run("./imgs/qrcode/baidu.png", []string{"qrcode"})

	if result["qrcode"].Frames[0].List[0].Values["value"].String != "https://www.baidu.com/" {
		jsonData, _ := json.Marshal(result)
		t.Fatal("qrcode byte err", string(jsonData))
	}
}

func TestQRCodeCNEN(t *testing.T) {
	result := Run("./imgs/qrcode/cnen.png", []string{"qrcode"})

	if result["qrcode"].Frames[0].List[0].Values["value"].String != "好用的imgpro" {
		jsonData, _ := json.Marshal(result)
		t.Fatal("qrcode cnen err", string(jsonData))
	}
}

func TestQRCodeCN(t *testing.T) {
	result := Run("./imgs/qrcode/cn.png", []string{"qrcode"})
	if result["qrcode"].Frames[0].List[0].Values["value"].String != "阿萨德" {
		jsonData, _ := json.Marshal(result)
		t.Fatal("qrcode cnen err", string(jsonData))
	}
}

// func TestQRCodeTest(t *testing.T) {
// 	result := Run("./imgs/qrcode/003.jpg", []string{"qrcode"})
// 	if result["qrcode"].Frames[0].List[0].Values["value"].String != "http://qm.qq.com/cgi-bin/qm/qr?k=LXqWJrE69ShewYXMOyls0HbEWpzaWoee" {
// 		jsonData, _ := json.Marshal(result)
// 		t.Fatal("qrcode cnen err", string(jsonData))
// 	}
// }
