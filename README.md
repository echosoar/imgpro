# ImgPro
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

| 功能 | 属性/方法名 | PNG | JPG | GIF | BMP | WebP | 
| --- | --- | --- | :---: | :---: | :---: | :---: | :---: |
| 图像大小 | size | ✅ | ✅ | ✅ | ✅ | ✅ | 
| 类型检测 | type | ✅ | ✅ | ✅ | ✅ | ✅ | 
| 宽度/高度| wh | ✅ | ✅ | ✅ |  |  | 
| 帧数 | frame | ✅ | ✅ |  | ✅ | |
| 颜色数据 | rgba |
| 色板 | hue | 
| EXIF 信息 | exif |
| 二维码识别 | qrcode |
| 文字识别 | ocr |
| 人体信息 | person |
| 物体信息 | object |
| 人脸信息 | face |
| 类型转换 | Conversion (方法) |
| 压缩 | Compress (方法) |

---

© MIT by echosoar
