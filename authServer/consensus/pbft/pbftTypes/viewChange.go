package pbftTypes

import (
	"github.com/DhunterAO/goAuthChain/common"
	"github.com/DhunterAO/goAuthChain/crypto"
	"github.com/DhunterAO/goAuthChain/crypto/types"
)

type ViewChange struct {
	InfoID    uint8
	NodeID    uint64
	ViewID    uint64
	NewViewID uint64
	Signature *types.Signature
}

func NewViewChange(nodeID uint64, viewID uint64, newViewID uint64) *ViewChange {
	newViewChange := &ViewChange{
		InfoID:    ViewChangeInfo,
		NodeID:    nodeID,
		ViewID:    viewID,
		NewViewID: newViewID,
		Signature: &types.Signature{},
	}
	return newViewChange
}

//func (vc *ViewChange) ToMap() map[string]interface{} {
//	vcMap := make(map[string]interface{})
//	vcMap["InfoID"] = vc.InfoID
//	vcMap["NodeID"] = vc.NodeID
//	vcMap["ViewID"] = vc.ViewID
//	vcMap["NewViewID"] = vc.NewViewID
//	vcMap["Signature"] = vc.Signature.Hex()
//	return vcMap
//}
//
//func (vc *ViewChange) ToPrettyJson() []byte {
//	vcPrettyJson, err := json.Marshal(vc.ToMap())
//	util.LogError(err)
//	return util.PrettyPrintJson(vcPrettyJson)
//}
//
//func NewViewChangeFromJson(vcJson []byte) *ViewChange {
//	vc := new(ViewChange)
//	err := json.Unmarshal(vcJson, vc)
//	if err != nil {
//		fmt.Println(err)
//	}
//	return vc
//}

func (vc *ViewChange) Bytes() []byte {
	vcBytes := make([]byte, 0)
	vcBytes = append(vcBytes, vc.InfoID)
	vcBytes = append(vcBytes, common.Uint64ToBytes(vc.NodeID)...)
	vcBytes = append(vcBytes, common.Uint64ToBytes(vc.ViewID)...)
	vcBytes = append(vcBytes, common.Uint64ToBytes(vc.NewViewID)...)
	return vcBytes
}

func (vc *ViewChange) CalcHash() *types.Hash {
	return crypto.CalcHash(vc.Bytes())
}

//func (vc *ViewChange) ToString() string {
//	return string(vc.ToPrettyJson())
//}

func (vc *ViewChange) CheckSignature() (bool, error) {
	pk, err := crypto.RecoverPkFromSig(vc.CalcHash(), vc.Signature)
	if err != nil {
		return false, err
	}
	return crypto.VerifySignature(pk, crypto.CalcHash(vc.Bytes()), vc.Signature), nil
}
