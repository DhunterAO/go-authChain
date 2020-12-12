package types

import (
	"testing"
)

func TestAuthorization_DeepCopy(t *testing.T) {
	auth := GenerateTestAuthorization()
	t.Log(auth.ToString())
	t.Log(auth.CheckSignature())
	t.Log(GetSenderAddrFromAuth(auth))
}
