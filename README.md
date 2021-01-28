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
- [x] size
- [x] image type [support: png / jpg / bmp / gif / webp]
- [ ] width/height
- [ ] image rgba data
- [ ] exif info
- [ ] qrcode
- [ ] color
- [ ] person/face detection
- [ ] object detection
- [ ] ocr
