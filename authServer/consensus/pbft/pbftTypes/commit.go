package pbftTypes

import (
	"encoding/json"
	"fmt"
	"github.com/DhunterAO/goAuthChain/common"
	"github.com/DhunterAO/goAuthChain/crypto"
	"github.com/DhunterAO/goAuthChain/crypto/types"
)

type Commit struct {
	InfoID    uint8
	NodeID    uint64
	ViewID    uint64
	BlockHash *types.Hash
	Signature *types.Signature
}

func NewCommit(nodeID uint64, viewID uint64, blockHash *types.Hash) *Commit {
	newCommit := &Commit{
		InfoID:    CommitInfo,
		NodeID:    nodeID,
		ViewID:    viewID,
		BlockHash: blockHash,
		Signature: &types.Signature{},
	}
	return newCommit
}

func (cmt *Commit) ToJson() []byte {
	cmtJson, err := json.Marshal(cmt)
	if err != nil {
		fmt.Println(err)
	}
	return cmtJson
}

func JsonToCommit(cmtJson []byte) *Commit {
	cmt := new(Commit)
	err := json.Unmarshal(cmtJson, cmt)
	if err != nil {
		fmt.Println(err)
	}
	return cmt
}

func (cmt *Commit) Bytes() []byte {
	cmtBytes := make([]byte, 0)
	cmtBytes = append(cmtBytes, cmt.InfoID)
	cmtBytes = append(cmtBytes, common.Uint64ToBytes(cmt.NodeID)...)
	cmtBytes = append(cmtBytes, common.Uint64ToBytes(cmt.ViewID)...)
	cmtBytes = append(cmtBytes, cmt.BlockHash.Bytes()...)
	return cmtBytes
}

func (cmt *Commit) CalcHash() *types.Hash {
	return crypto.CalcHash(cmt.Bytes())
}

func (cmt *Commit) CheckSignature() bool {
	pk, err := crypto.RecoverPkFromSig(cmt.CalcHash(), cmt.Signature)
	if err != nil {
		return false
	}
	return crypto.VerifySignature(pk, crypto.CalcHash(cmt.Bytes()), cmt.Signature)
}
