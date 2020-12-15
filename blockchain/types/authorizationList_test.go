package types

import (
	"testing"
)

func TestAuthList_Add(t *testing.T) {
	auth := GenerateTestAuthorization()
	authList := NewAuthList()
	if authList.Add(auth) == false {
		t.Error("add authorization failed")
	}
}
