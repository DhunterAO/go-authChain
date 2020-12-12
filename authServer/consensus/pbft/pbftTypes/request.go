package pbftTypes

import (
	"encoding/json"
	"fmt"
	"github.com/DhunterAO/goAuthChain/common"
	"github.com/DhunterAO/goAuthChain/crypto"
	"github.com/DhunterAO/goAuthChain/crypto/types"
)

type Request struct {
	InfoID    uint8
	NodeID    uint64
	ViewID    uint64
	BlockJson []byte
	Signature *types.Signature
}

func NewRequest(nodeID uint64, viewID uint64, blockJson []byte) *Request {
	newRequest := &Request{
		InfoID:    RequestInfo,
		NodeID:    nodeID,
		ViewID:    viewID,
		BlockJson: blockJson,
		Signature: &types.Signature{},
	}
	return newRequest
}

func BytesToRequest(b []byte) *Request {
	newReq := new(Request)
	err := json.Unmarshal(b, newReq)
	if err != nil {
		fmt.Println(err)
	}
	return newReq
}

func (req *Request) Sign(signature *types.Signature) {
	req.Signature = signature
}

func (req *Request) Bytes() []byte {
	respBytes := make([]byte, 0)
	respBytes = append(respBytes, req.InfoID)
	respBytes = append(respBytes, common.Uint64ToBytes(req.NodeID)...)
	respBytes = append(respBytes, common.Uint64ToBytes(req.ViewID)...)
	respBytes = append(respBytes, req.BlockJson...)
	return respBytes
}

func (req *Request) ToJson() []byte {
	reqJson, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	return reqJson
}

func (req *Request) CalcHash() *types.Hash {
	return crypto.CalcHash(req.Bytes())
}

func (req *Request) CheckSignature() (bool, error) {
	pk, err := crypto.RecoverPkFromSig(req.CalcHash(), req.Signature)
	if err != nil {
		return false, err
	}
	return crypto.VerifySignature(pk, crypto.CalcHash(req.Bytes()), req.Signature), nil
}

//func (req *Request) ToMap() map[string]interface{} {
//	reqMap := make(map[string]interface{})
//	reqMap["InfoID"] = req.InfoID
//	reqMap["NodeID"] = req.NodeID
//	reqMap["ViewID"] = req.ViewID
//	b := blockchainTypes.NewBlockFromJson(req.BlockJson)
//	reqMap["BlockJson"] = b.ToMap()
//	reqMap["Signature"] = req.Signature.ToHex()
//	return reqMap
//}
//
//func (req *Request) ToPrettyJson() []byte {
//	reqPreJson, err := json.Marshal(req.ToMap())
//	util.LogError(err)
//	return util.PrettyPrintJson(reqPreJson)
//}
//
//func (req *Request) ToString() string {
//	return string(req.ToPrettyJson())
//}
