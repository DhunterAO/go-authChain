package types

import (
	"fmt"
	"testing"
)

func TestGenerateTestHeader(t *testing.T) {
	h := GenerateTestHeader()
	h.PrettyPrint()
}

func TestGenerateTestBlock(t *testing.T) {
	b := GenerateTestBlock()
	fmt.Println(b.ToPrettyJson())
}
