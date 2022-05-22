package test

import (
	"encoding/json"
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestPositionJpg(t *testing.T) {
	result := Run("./imgs/exif-gps.jpeg", []string{"position"})

	if result["position"].String != "42 45 3985/1000 N,84 29 28508/1000 W" {
		jsonData, _ := json.Marshal(result)
		t.Fatal("position err", string(jsonData))
	}
}