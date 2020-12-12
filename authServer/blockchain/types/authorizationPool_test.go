package types

import (
	"fmt"
	"testing"
)

func TestNewAuthorizationPool(t *testing.T) {
	//fmt.Println(auth.ToString())
	var auths []*Authorization
	for i := 0; i < 3; i += 1 {
		auth := GenerateTestAuthorization()
		auths = append(auths, auth)
	}
	pool, _ := NewAuthPool(auths)
	fmt.Println(pool.ToString())
}
