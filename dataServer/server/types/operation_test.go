package types

import (
	"fmt"
	"testing"
)

func TestGenerateTestOperation(t *testing.T) {
	op := GenerateTestOperation()
	fmt.Println(op.ToString())
}
