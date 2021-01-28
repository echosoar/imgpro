# imgpro
image processor

### Usage
```go
import (
  "github.com/echosoar/imgpro"
)

func main() {
	result := imgpro.Run("./test/imgs/go.png", []string{"size"})
	if result["size"].Int != 60746 {
		panic("size error")
	}
}
```

### Features
- [x] size
- [ ] width/height
- [ ] image type 
- [ ] image rgba data
- [ ] exif info
- [ ] qrcode
- [ ] color
- [ ] person/face detection
- [ ] object detection
- [ ] ocr
