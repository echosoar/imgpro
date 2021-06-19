package processor

import (
	"bufio"
	"fmt"
	"os"

	img "github.com/echosoar/imgpro/core"
)

// https://exiftool.org/TagNames/EXIF.html
var exifTypeName = map[int]string{
	0x010f: "Make",
	0x0110: "Model",
}

// ExifProcessor get image type
func ExifProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:         []string{"exif"},
		Precondition: []string{"size"},
		Runner:       exifRunner,
	})
}

func findByteIndex(find []byte, from []byte) int {
	for i, bt := range from {
		if bt == find[0] {
			isMatch := true
			for fi, fbt := range find {
				if fbt != from[i+fi] {
					isMatch = false
					break
				}
			}
			if isMatch {
				return i
			}
		}
	}
	return -1
}

func byteToInt(bytes []byte) int {
	byteLen := len(bytes)
	res := 0
	for index, bt := range bytes {
		res += int(bt) << ((byteLen - index - 1) * 8)
	}
	return res
}

func exifRunner(core *img.Core) map[string]img.Value {
	f, err := os.Open(core.FilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	size := core.Result["size"].Int
	fileBytes := make([]byte, size)
	reader := bufio.NewReader(f)
	_, readErr := reader.Read(fileBytes)
	if readErr != nil {
		panic(readErr)
	}

	value := make(map[string]img.Value)
	app1BytesIndex := findByteIndex([]byte("\xFF\xE1"), fileBytes)
	exifBytesIndex := findByteIndex([]byte("\x45\x78\x69\x66"), fileBytes)
	if exifBytesIndex-app1BytesIndex != 4 {
		// error
	}
	// exifSizeByte := fileBytes[app1BytesIndex+2 : exifBytesIndex]
	// exifSize := int(exifSizeByte[0])<<8 + int(exifSizeByte[1])

	// exifBytes := fileBytes[app1BytesIndex:exifSize];
	// 0-1 app1
	// 2-3 size
	// 4-7 "exif"
	// 8-9 00
	// 10-11 MM/II MM 大端 ; II 小端
	// 12-13 00 2A
	// 14-17 00 00 00 08
	// 18-19 tag size 12
	// 20-21 01 0F tag
	// 22-23 00 02 格式为2
	// 24-27 00 00 00 06 count 6
	// 28-31 00 00 00 86 偏移量 + 1D = a3

	tagSize := byteToInt(fileBytes[app1BytesIndex+18 : app1BytesIndex+20])
	tagStartIndex := app1BytesIndex + 20

	for tagIndex := 0; tagIndex < tagSize; tagIndex++ {
		curTagStartIndex := tagStartIndex + tagIndex*12
		tagName := byteToInt(fileBytes[curTagStartIndex : curTagStartIndex+2])

		tagNameString, exists := exifTypeName[tagName]
		if !exists {
			fmt.Println("tag", fileBytes[curTagStartIndex:curTagStartIndex+2])
			continue
		}
		fmt.Println("tagNameString", tagNameString)
		// 2 string
		// 3 int
		tagType := byteToInt(fileBytes[curTagStartIndex+2 : curTagStartIndex+4])
		tagValueCount := byteToInt(fileBytes[curTagStartIndex+4 : curTagStartIndex+8])
		tagValueOffset := exifBytesIndex + 6 + byteToInt(fileBytes[curTagStartIndex+8:curTagStartIndex+12])
		tagValue := fileBytes[tagValueOffset : tagValueOffset+tagValueCount]

		if tagType == 2 {
			value[tagNameString] = img.Value{
				Type:   img.ValueTypeString,
				String: string(tagValue),
			}
		} else {

		}
	}
	fmt.Println("value", value, "<---")
	return map[string]img.Value{
		"exif": {
			Type:   img.ValueTypeMap,
			Values: value,
		},
	}
}
