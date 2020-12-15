package pbftTypes

import (
	"fmt"
	types2 "github.com/DhunterAO/goAuthChain/blockchain/types"
	"github.com/DhunterAO/goAuthChain/crypto"
	"github.com/DhunterAO/goAuthChain/crypto/types"
	"testing"
)

func TestNewRequest(t *testing.T) {
	block := types2.GenerateTestBlock()
	request := NewRequest(uint64(0), 0, block.ToJson())
	fmt.Println(string(request.Bytes()))
	sk, _ := types.GenerateKey()
	sig, _ := crypto.Sign(request.CalcHash(), sk)
	request.Sign(sig)
	fmt.Println(request.CheckSignature())
	fmt.Println("--------------------")

	req2 := BytesToRequest(request.Bytes())
	fmt.Println(string(req2.ToJson()))
}
