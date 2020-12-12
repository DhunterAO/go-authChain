package types

import (
	"encoding/json"
	"fmt"
	"github.com/DhunterAO/goAuthChain/common"
	"github.com/DhunterAO/goAuthChain/crypto"
	"github.com/DhunterAO/goAuthChain/crypto/types"
	"math/rand"
)

type OpCode uint8

const (
	OP_ADD OpCode = iota
	OP_DEL
	OP_UPT
	OP_QRY
)

type Operation struct {
	Address   *types.Address
	Nonce     uint64
	OpCode    OpCode
	Key       []byte
	Value     []byte
	Ext       []byte
	Signature *types.Signature
}

func NewOperation(address *types.Address, opCode OpCode, key []byte, value []byte, signature *types.Signature) *Operation {
	op := &Operation{
		Address: address,
		//Nonce:   nonce, // no need for nonce
		OpCode: opCode,
		Key:    key,
		Value:  value,
	}
	if signature != nil {
		op.Signature = signature
	} else {
		op.Signature = &types.EmptySignature
	}
	return op
}

func (op *Operation) ToString() string {
	return string(op.ToJson())
}

func (op *Operation) ToJson() []byte {
	opJson, err := json.Marshal(op.ToMap())
	if err != nil {
		fmt.Println(err)
	}
	return opJson
}

func (op *Operation) ToMap() map[string]interface{} {
	opMap := make(map[string]interface{})
	opMap["address"] = op.Address.Hex()
	opMap["opCode"] = op.OpCode
	opMap["key"] = op.Key
	opMap["value"] = op.Value
	//opMap["nonce"] = op.Nonce
	opMap["signature"] = op.Signature.Hex()
	return opMap
}

func JsonToOperation(opJson []byte) *Operation {
	op := new(Operation)
	err := json.Unmarshal(opJson, op)
	if err != nil {
		fmt.Println(err)
	}
	return op
}

func (op *Operation) ToBytes() []byte {
	opBytes := append(op.Address[:], common.Uint64ToBytes(op.Nonce)...)
	opBytes = append(opBytes, byte(op.OpCode))
	opBytes = append(opBytes, op.Key[:]...)
	opBytes = append(opBytes, op.Value[:]...)
	return opBytes
}

func (op *Operation) Sign(sig *types.Signature) {
	op.Signature = sig
}

func (op *Operation) CheckSignature() (bool, error) {
	hash := op.CalcHash()
	signature := op.Signature
	pk, err := crypto.RecoverPkFromSig(hash, signature)
	if err != nil {
		return false, err
	}
	return crypto.VerifySignature(pk, hash, signature), nil
}

func (op *Operation) CalcHash() *types.Hash {
	mess := op.ToBytes()
	return crypto.CalcHash(mess)
}

func (op *Operation) DeepCopy() *Operation {
	newOp := &Operation{
		Address:   op.Address,
		Nonce:     op.Nonce,
		Key:       op.Key,
		Value:     op.Value,
		Signature: op.Signature,
	}
	return newOp
}

func GenerateTestOperation() *Operation {
	address := types.RandAddress()
	opCode := OpCode(rand.Uint32())
	key := []byte(common.RandString(4))
	value := []byte(common.RandString(10))
	op := NewOperation(address, opCode, key, value, nil)
	sk, _ := types.GenerateKey()
	sig, _ := crypto.Sign(op.CalcHash(), sk)
	op.Sign(sig)
	return op
}
