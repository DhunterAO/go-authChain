package pbftTypes

import (
	"fmt"
	"github.com/DhunterAO/goAuthChain/crypto"
	"github.com/DhunterAO/goAuthChain/crypto/types"
	"testing"
)

func TestNewViewChange(t *testing.T) {
	vc := NewViewChange(0, 1, 2)

	sk, _ := types.LoadECDSA("../../../data/keys/test.key")
	fmt.Println(sk)
	sig, _ := crypto.Sign(vc.CalcHash(), sk)
	vc.Signature = sig
	fmt.Println(vc.CheckSignature())
}
