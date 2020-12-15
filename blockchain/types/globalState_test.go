package types

import (
	"github.com/DhunterAO/goAuthChain/crypto/types"
	"testing"
)

func TestGenerateTestAuthorization(t *testing.T) {
	auth := GenerateTestAuthorization()
	addr, _ := GetSenderAddrFromAuth(auth)
	gmList := make(map[types.Address]string)
	gmList[*addr] = "WX"
	var preState []*AccountState
	gs := NewGlobalState(gmList, preState)
	if res, err := gs.CheckAuthorization(auth); res == false || err != nil {
		if err == nil {
			t.Error("check authorization fail")
		} else {
			t.Error(err)
		}
	}
	if res, err := gs.AcceptAuthorization(auth); res == false || err != nil {
		if err == nil {
			t.Error("accept authorization fail")
		} else {
			t.Error(err)
		}
	}
}
