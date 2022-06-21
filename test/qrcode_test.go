package test

import (
	"encoding/json"
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestQRCodeNumeric(t *testing.T) {
	result := Run("./imgs/qrcode/1234567.png", []string{"qrcode"})

	if result["qrcode"].Frames[0].List[0].Values["value"].String != "1234567" {
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

func TestQRCodeTest001(t *testing.T) {
	result := Run("./imgs/qrcode/001.jpg", []string{"qrcode"})
	if result["qrcode"].Frames[0].List[0].Values["value"].String != "imgpro" {
		jsonData, _ := json.Marshal(result)
		t.Fatal("qrcode cnen err", string(jsonData))
	}
}

func TestQRCodeTest002(t *testing.T) {
	result := Run("./imgs/qrcode/002.jpg", []string{"qrcode"})
	if result["qrcode"].Frames[0].List[0].Values["value"].String != "imgpro imgpro imgpro imgpro imgpro imgpro imgpro imgpro imgpro imgpro imgpro imgpro imgpro imgpro imgpro imgpro imgpro imgpro imgpro imgpro imgpro" {
		jsonData, _ := json.Marshal(result)
		t.Fatal("qrcode cnen err", string(jsonData))
	}
}

func TestQRCodeTest003(t *testing.T) {
	result := Run("./imgs/qrcode/003.png", []string{"qrcode"})
	if result["qrcode"].Frames[0].List[0].Values["value"].String != "http://qm.qq.com/cgi-bin/qm/qr?k=LXqWJrE69ShewYXMOyls0HbEWpzaWoee" {
		jsonData, _ := json.Marshal(result)
		t.Fatal("qrcode cnen err", string(jsonData))
	}
}

func TestQRCodeTestAlipay(t *testing.T) {
	result := Run("./imgs/qrcode/alipay.jpeg", []string{"qrcode"})
	if result["qrcode"].Frames[0].List[0].Values["value"].String != "https://qr.alipay.com/fkx1204145jqmapfxwzbzfa" {
		jsonData, _ := json.Marshal(result)
		t.Fatal("qrcode cnen err", string(jsonData))
	}
}

func TestQRCodeTestECI(t *testing.T) {
	result := Run("./imgs/qrcode/eci.png", []string{"qrcode"})
	if result["qrcode"].Frames[0].List[0].Values["value"].String != "https://qm.qq.com/cgi-bin/qm/qr?k=ELaNID3csLNBJSUXW91-Sbv8Bad22pGq&authKey=CmpRJqghJxNEv/Sb7F9Z3SJHMMZpshhOevlJr+nkP7X9QPIAqm4dCQNJvUxMbJ0U&noverify=0" {
		jsonData, _ := json.Marshal(result)
		t.Fatal("qrcode cnen err", string(jsonData))
	}
}

func TestQRCodeTestNX(t *testing.T) {
	result := Run("./imgs/qrcode/nx.jpeg", []string{"qrcode"})
	if result["qrcode"].Frames[0].List[0].Values["value"].String != "https://tm-web.pin-dao.cn/nx-xp?oid=668361335772295168&sc=26074125&pid=LY260741252022061800160" {
		jsonData, _ := json.Marshal(result)
		t.Fatal("qrcode cnen err", string(jsonData))
	}
}

func BenchmarkTestNX(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Run("./imgs/qrcode/nx.jpeg", []string{"qrcode"})
	}
}
