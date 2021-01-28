package imgpro

import (
	"testing"
)

func TestSize(t *testing.T) {
	result := Run("./test/imgs/go.png", []string{"size"})
	if result["size"].Int != 60746 {
		t.Fatal("size error")
	}
}
