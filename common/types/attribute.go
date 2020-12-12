package types

import (
	"encoding/json"
	"fmt"
	"github.com/DhunterAO/goAuthChain/common"
)

type Attribute struct {
	Name string
	Dur  *Duration
}

func NewAttribute(name string, dur *Duration) *Attribute {
	attr := &Attribute{
		Name: name,
		Dur:  dur,
	}
	return attr
}

//func JsonToAttr(attrJson []byte) Attribute {
//	var attr Attribute
//	err := json.Unmarshal(attrJson, &attr)
//	if err != nil {
//		fmt.Println(err)
//	}
//	return attr
//}

func (a *Attribute) ToBytes() []byte {
	attrBytes := append([]byte(a.Name), a.Dur.ToBytes()...)
	return attrBytes
}

func (a *Attribute) DeepCopy() *Attribute {
	newAttr := &Attribute{
		Name: a.Name,
		Dur:  a.Dur.DeepCopy(),
	}
	return newAttr
}

func (a *Attribute) ToMap() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap["name"] = a.Name
	attrMap["duration"] = a.Dur.ToMap()
	return attrMap
}

func (a *Attribute) ToJson() []byte {
	attrJson, err := json.Marshal(a.ToMap())
	if err != nil {
		fmt.Println(err)
	}
	return attrJson
}

func GenerateTestAttribute() *Attribute {
	name := common.RandString(8)
	duration := GenerateTestDuration()
	attr := NewAttribute(name, duration)
	return attr
}
