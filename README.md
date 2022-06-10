# Imgpro

[![CircleCI](https://circleci.com/gh/echosoar/imgpro/tree/main.svg?style=svg&circle-token=355449f4d49bf63a561c68c57221688dadc48691)]((https://circleci.com/gh/echosoar/imgpro/tree/main))

Imgpro is a multifunctional image information recognition library, supporting a variety of image formats. And it can be run in the browser through WebAssembly(wasm).


Online Demo: [imgpro](https://echosoar.github.io/imgpro/)
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
  // run by file path
  result := imgpro.Run("./test/imgs/go.png", []string{"size", "type"})
  // you can also run by file binary data
  // result := imgpro.RunBinary(binary, attributes)
  
  if result["size"].Int != 60746 {
    panic("size error")
  }
  
  if result["type"].String != "png" {
    panic("type error")
  }
}
```

### Method
#### Run

> Get image information by local file path

|Param Index|Param Name|Type|Examples|
| --- | --- | --- |--- |
| 0 | filePath | string | "./test/imgs/go.png" |
| 1 | attributes | []string | []string{"size", "type", "rgba", "hue", "qrcode"} |

```go
import (
  "github.com/echosoar/imgpro"
)

func main() {
  result := imgpro.Run("./test/imgs/go.png", []string{"size", "type", "rgba", "hue", "qrcode"})
}
```

#### RunBinary
> Get image information by file binary data

|Param Index|Param Name|Type|Examples|
| --- | --- | --- |--- |
| 0 | fileBinary | []byte | reader.Read(binary) |
| 1 | attributes | []string | []string{"size", "type", "rgba", "hue", "qrcode"} |

```go
import (
  "bufio"
  "os"

  "github.com/echosoar/imgpro"
)

func main() {
  fileHandler, err := os.Open(filePath)
	if err != nil {
		panic("open error")
	}
	defer fileHandler.Close()
	fileBytes := make([]byte, size)
	reader := bufio.NewReader(fileHandler)
	_, readErr := reader.Read(fileBytes)
	if readErr != nil {
		panic("file read error")
	}
  result := imgpro.RunBinary(fileBytes, []string{"size", "type", "rgba", "hue", "qrcode"})
}
```

### Features

| Features | Attribute | PNG | JPG | GIF | BMP | WebP |
| --- | --- | :---: | :---: | :---: | :---: | :---: | 
| File Size | size | ✅ | ✅ | ✅ | ✅ | ✅ |
| Format Detect | type | ✅ | ✅ | ✅ | ✅ | ✅ |
| Width/Height| wh | ✅ | ✅ | ✅ | ✅ | ✅ |
| Frames | frame | ✅ | ✅ |
| Color data | rgba | ✅ | ✅ |
| Color proportion | hue | ✅ | ✅ |
| Exif | exif |  | ✅ |
| Create Time | time | |✅ | | | |
| Position(GPS) Info | position | |✅  | | | |
| Device Info | device | | ✅| | | |
| QR Code | qrcode | ✅ | ✅| | | |

---

© MIT by echosoar
