# ImgPro
[![Build Status](https://circleci.com/gh/echosoar/imgpro.svg?style=shield)](https://circleci.com/gh/echosoar/imgpro)

多功能图像信息识别与图像处理库，支持多种图片格式

### 如何使用
#### 安装
```shell
$ go get github.com/echosoar/imgpro
```
#### 代码中使用
```go
import (
  "github.com/echosoar/imgpro"
)

func main() {
  result := imgpro.Run("./test/imgs/go.png", []string{"size", "type"})
  
  if result["size"].Int != 60746 {
    panic("size error")
  }
  
  if result["type"].String != "png" {
    panic("type error")
  }
}
```

### 功能列表

| 功能 | 属性/方法名 | PNG | JPG | GIF | BMP | WebP | APNG | AVIF |
| --- | --- | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| 图像大小 | size | ✅ | ✅ | ✅ | ✅ | ✅ |
| 类型检测 | type | ✅ | ✅ | ✅ | ✅ | ✅ |
| 宽度/高度| wh | ✅ | ✅ | ✅ | ✅ | ✅ |
| 帧数| frame | ✅ | ✅ |
| 颜色数据 | rgba | ✅ | ✅ |
| 色调（颜色占比） | hue | ✅ | ✅ |
| 色板（颜色归类） | palette | 
| Exif | exif |  | ✅ |
| 时间 | time | |✅ | | | |
| 地点 | position | | ✅| | | |
| 设备 | device | | | | | |
| 二维码 | qrcode |
| 文字 | ocr |
| 人脸 | face |
| 人体信息 | person |
| 动物信息 | animal |
| 物体信息 | object |

---

© MIT by echosoar
