package types

import (
	"testing"
)

func TestAuthorization_CheckSignature(t *testing.T) {
	auth := GenerateTestAuthorization()
	if auth.CheckSignature() != true {
		t.Error("Check signature error")
	}
}
