package types

import (
	util2 "goauth/util"
	"testing"
)

func TestDuration_DeepCopy(t *testing.T) {
	d := GenerateTestDuration()
	d2 := d.DeepCopy()
	//fmt.Println(d.ToString(), d2.ToString())
	d.Start += 1
	d.End += 1
	//fmt.Println(d.ToString(), d2.ToString())
	util2.PrettyPrintJson(d.ToJson())
	util2.PrettyPrintJson(d2.ToJson())
}