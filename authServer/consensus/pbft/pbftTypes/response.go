package pbftTypes

import (
	"encoding/json"
	"fmt"
	"github.com/DhunterAO/goAuthChain/common"
	"github.com/DhunterAO/goAuthChain/crypto"
	"github.com/DhunterAO/goAuthChain/crypto/types"
)

type Response struct {
	InfoID    uint8
	NodeID    uint64
	ViewID    uint64
	BlockHash types.Hash
	Signature *types.Signature
}

func NewResponse(nodeID uint64, viewID uint64, blockHash types.Hash) *Response {
	newResponse := &Response{
		InfoID:    ResponseInfo,
		NodeID:    nodeID,
		ViewID:    viewID,
		BlockHash: blockHash,
		Signature: &types.Signature{},
	}
	return newResponse
}

func (resp *Response) ToJson() []byte {
	respJson, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
	}
	return respJson
}

//func (resp *Response) ToMap() map[string]interface{} {
//	respMap := make(map[string]interface{})
//	respMap["InfoID"] = resp.InfoID
//	respMap["NodeID"] = resp.NodeID
//	respMap["ViewID"] = resp.ViewID
//	respMap["BlockHash"] = resp.BlockHash.Hex()
//	respMap["Signature"] = resp.Signature.Hex()
//	return respMap
//}
//
//func (resp *Response) ToPrettyJson() []byte {
//	respPreJson, err := json.Marshal(resp.ToMap())
//	util.LogError(err)
//	return util.PrettyPrintJson(respPreJson)
//}
//
//func (resp *Response) ToString() string {
//	return string(resp.ToPrettyJson())
//}

//func JsonToResponse(respJson []byte) *Response {
//	resp := new(Response)
//	err := json.Unmarshal(respJson, resp)
//	if err != nil {
//		fmt.Println(err)
//	}
//	return resp
//}

func (resp *Response) Bytes() []byte {
	respBytes := make([]byte, 0)
	respBytes = append(respBytes, resp.InfoID)
	respBytes = append(respBytes, common.Uint64ToBytes(resp.NodeID)...)
	respBytes = append(respBytes, common.Uint64ToBytes(resp.ViewID)...)
	respBytes = append(respBytes, resp.BlockHash.Bytes()...)
	return respBytes
}

func (resp *Response) CalcHash() *types.Hash {
	return crypto.CalcHash(resp.Bytes())
}

func (resp *Response) CheckSignature() (bool, error) {
	pk, err := crypto.RecoverPkFromSig(resp.CalcHash(), resp.Signature)
	if err != nil {
		return false, err
	}
	return crypto.VerifySignature(pk, crypto.CalcHash(resp.Bytes()), resp.Signature), nil
}
