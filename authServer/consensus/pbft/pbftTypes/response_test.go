package pbftTypes

import (
	"fmt"
	"github.com/DhunterAO/goAuthChain/crypto/types"
	"testing"
)

func TestNewResponse(t *testing.T) {
	resp := NewResponse(2, 0, *types.RandHash())
	fmt.Println(string(resp.ToJson()))
	fmt.Println(resp.CheckSignature())
}
