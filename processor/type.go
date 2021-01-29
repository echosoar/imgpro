package processor

import (
	"bufio"
	"bytes"
	"os"

	img "github.com/echosoar/imgpro/core"
)

// TypeProcessor get image type
func TypeProcessor(imgCore *img.Core) {
	imgCore.Bind(&img.Processor{
		Keys:         []string{"type"},
		Precondition: []string{"size"},
		Runner:       typeRunner,
	})
}

func typeRunner(core *img.Core) map[string]img.Value {
	f, err := os.Open(core.FilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	size := core.Result["size"].Int
	readSize := 14 // for webp
	if size < readSize {
		readSize = size
	}

	fileBytes := make([]byte, readSize)
	reader := bufio.NewReader(f)
	count, readErr := reader.Read(fileBytes)
	if readErr != nil {
		panic(readErr)
	}
	fileBytes = fileBytes[:count]
	// refs: https://golang.org/src/net/http/sniff.go
	imgType := "unknow"
	if bytes.HasPrefix(fileBytes, []byte("BM")) {
		imgType = "bmp"
	} else if bytes.HasPrefix(fileBytes, []byte("GIF87a")) || bytes.HasPrefix(fileBytes, []byte("GIF89a")) {
		imgType = "gif"
	} else if bytes.HasPrefix(fileBytes, []byte("\x89PNG\x0D\x0A\x1A\x0A")) {
		imgType = "png"
	} else if bytes.HasPrefix(fileBytes, []byte("\xFF\xD8\xFF")) {
		imgType = "jpg"
	} else if len(fileBytes) >= 14 {
		isWebp := func() bool {
			mask := []byte("\xFF\xFF\xFF\xFF\x00\x00\x00\x00\xFF\xFF\xFF\xFF\xFF\xFF")
			pat := []byte("RIFF\x00\x00\x00\x00WEBPVP")
			for i, pb := range pat {
				maskedData := fileBytes[i] & mask[i]
				if maskedData != pb {
					return false
				}
			}
			return true
		}()
		if isWebp {
			imgType = "webp"
		}
	}
	return map[string]img.Value{
		"type": {
			Type:   img.ValueTypeString,
			String: imgType,
		},
	}
}
