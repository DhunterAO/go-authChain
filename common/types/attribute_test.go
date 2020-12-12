package types

import (
	"fmt"
	"testing"
)

func TestJsonConvert(t *testing.T) {
	a := GenerateTestAttribute()
	//fmt.Println(a)
	aJson := a.ToJson()
	fmt.Println(string(aJson))

	b := GenerateTestAttribute()
	fmt.Println(string(b.ToJson()))

	//b := JsonToAttr(aJson)
	//fmt.Println(b)
	//if a != b {
	//	fmt.Println("convert fail")
	//}
}
