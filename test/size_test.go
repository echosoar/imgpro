package test

import (
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestSize(t *testing.T) {
	result := Run("./imgs/go.png", []string{"size"})
	if result["size"].Int != 60746 {
		t.Fatal("size error")
	}
}
