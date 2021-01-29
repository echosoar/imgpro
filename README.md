# imgpro
image processor

### Usage
#### Install
```shell
$ go get github.com/echosoar/imgpro
```
#### Use
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

### Features

| Feature | PNG | JPG | GIF | BMP | WebP |
| --- | :---: | :---: | :---: | :---: | :---: |
| size | ✅ | ✅ | ✅ | ✅ | ✅ | 
| type detect | ✅ | ✅ | ✅ | ✅ | ✅ | 
| width/height | ✅ | ✅ | ✅ |  |  | 
| rgba data |
| exif info |
| qrcode recognition |
| hue |
| person detect |
| object detect |
| ocr |
| face detect |
| compress |
| type conversion |

