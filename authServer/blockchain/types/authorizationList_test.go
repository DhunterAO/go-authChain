package types

import (
	"fmt"
	"testing"
)

func TestAuthList_Add(t *testing.T) {
	auth := GenerateTestAuthorization()
	authList := NewAuthList()
	authList.Add(auth)
	fmt.Println(authList.ToString())
}
