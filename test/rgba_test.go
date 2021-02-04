package test

import (
	"testing"

	. "github.com/echosoar/imgpro"
)

func TestRGBAPng(t *testing.T) {
	result := Run("./imgs/go.png", []string{"rgba"})
	t.Fatal(result)
}
