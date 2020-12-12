package types

import (
	"fmt"
	"testing"
)

func TestAccountState_ToJson(t *testing.T) {
	as := GenerateTestAccountState()
	fmt.Println(as.ToString())
}
