package types

import (
	"github.com/DhunterAO/goAuthChain/common"
	"testing"
)

func TestDuration_DeepCopy(t *testing.T) {
	d := GenerateTestDuration()
	d2 := d.DeepCopy()
	//fmt.Println(d.ToString(), d2.ToString())
	d.Start += 1
	d.End += 1
	//fmt.Println(d.ToString(), d2.ToString())
	common.PrettyPrintJson(d.ToJson())
	common.PrettyPrintJson(d2.ToJson())
}
