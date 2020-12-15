package types

import (
	"testing"
)

func TestGenerateTestAccountState(t *testing.T) {
	as := GenerateTestAccountState()
	for attr := range as.Attributes {
		if res, err := as.CheckAttribute(as.Attributes[attr]); res == false || err != nil {
			if err == nil {
				t.Error("The attribute in account should be valid but not", attr)
			} else {
				t.Error(err)
			}
		}
	}
}
