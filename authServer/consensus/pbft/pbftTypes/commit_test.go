package pbftTypes

import (
	"fmt"
	"github.com/DhunterAO/goAuthChain/crypto/types"
	"testing"
)

func TestNewCommit(t *testing.T) {
	cmt := NewCommit(1, 0, &types.Hash{})
	fmt.Println(cmt)

	fmt.Println(cmt.CheckSignature())
}
