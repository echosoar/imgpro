# Imgpro

![CircleCI](https://circleci.com/gh/echosoar/imgpro/tree/main.svg?style=svg)

Multifunctional image information recognition library, supporting a variety of image formats. 

### Usage
#### Initial
```shell
$ go get github.com/echosoar/imgpro
```
#### Use in code
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

| Features | Attribute | PNG | JPG | GIF | BMP | WebP | APNG | AVIF |
| --- | --- | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| File Size | size | ✅ | ✅ | ✅ | ✅ | ✅ |
| Format Detect | type | ✅ | ✅ | ✅ | ✅ | ✅ |
| Width/Height| wh | ✅ | ✅ | ✅ | ✅ | ✅ |
| Frames | frame | ✅ | ✅ |
| Color data | rgba | ✅ | ✅ |
| Color proportion） | hue | ✅ | ✅ |
| Exif | exif |  | ✅ |
| Create Time | time | |✅ | | | |
| Position(GPS) Info | position | | | | | |
| Device Info | device | | ✅| | | |

---

© MIT by echosoar
