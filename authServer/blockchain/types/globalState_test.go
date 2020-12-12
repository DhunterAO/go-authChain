package types

import (
	"github.com/DhunterAO/goAuthChain/crypto/types"
	"testing"
)

func TestNewGlobalState(t *testing.T) {
	auth := GenerateTestAuthorization()
	t.Log(auth.ToString())
	addr, _ := GetSenderAddrFromAuth(auth)
	gmList := make(map[types.Address]string)
	gmList[*addr] = "WX"
	var preState []*AccountState
	gs := NewGlobalState(gmList, preState)
	t.Log(gs.ToString())
	t.Log(gs.CheckAuthorization(auth))

	flag, _ := gs.AcceptAuthorization(auth)
	t.Log(flag)
	t.Log("--------------------")
	t.Log(gs.ToString())
}
