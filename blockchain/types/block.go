package types

import (
	"encoding/json"
	"fmt"
	"github.com/DhunterAO/goAuthChain/common"
	commonType "github.com/DhunterAO/goAuthChain/common/types"
	"github.com/DhunterAO/goAuthChain/crypto"
	cryptoType "github.com/DhunterAO/goAuthChain/crypto/types"
	"math/rand"
)

type Header struct {
	ParentHash *cryptoType.Hash
	AuthRoot   *cryptoType.Hash
	//StateRoot     *types.Hash
	Timestamp commonType.Timestamp
}

// create new header
func NewHeader(parentHash *cryptoType.Hash, authRoot *cryptoType.Hash, timestamp commonType.Timestamp) *Header {
	h := &Header{
		ParentHash: parentHash,
		AuthRoot:   authRoot,
		Timestamp:  timestamp,
	}
	return h
}

type Block struct {
	Header         *Header
	Authorizations []*Authorization
	//State          *GlobalState
}

func NewBlock(h *Header, auths []*Authorization) *Block {
	b := &Block{
		Header:         h,
		Authorizations: auths,
	}
	return b
}

func JsonToBlock(blockJson []byte) *Block {
	b := new(Block)
	err := json.Unmarshal(blockJson, b)
	if err != nil {
		fmt.Println(err)
	}
	return b
}

func (b *Block) CalcHash() *cryptoType.Hash {
	message := b.Header.ToBytes()
	hash := crypto.CalcHash(message)
	return hash
}

func (b *Block) CalcAuthRoot() *cryptoType.Hash {
	var auths []byte
	for i := range b.Authorizations {
		auths = append(auths, b.Authorizations[i].Bytes()...)
	}
	return crypto.CalcHash(auths)
}

//func (b *Block) CalcStateRoot() *types.Hash {
//	return b.State.CalcHash()
//}

func (h *Header) SetAuthRoot(authRoot *cryptoType.Hash) {
	h.AuthRoot = authRoot
}

func (b *Block) SetRoot() {
	authRoot := b.CalcAuthRoot()
	b.Header.SetAuthRoot(authRoot)
}

func (h *Header) ToBytes() []byte {
	var hBytes []byte
	hBytes = append(h.ParentHash[:], h.AuthRoot[:]...)
	hBytes = append(hBytes, h.Timestamp.Bytes()...)
	return hBytes
}

func (h *Header) DeepCopy() *Header {
	newHeader := &Header{
		ParentHash: h.ParentHash.DeepCopy(),
		AuthRoot:   h.AuthRoot.DeepCopy(),
		Timestamp:  h.Timestamp,
	}
	return newHeader
}

func (b *Block) DeepCopy() *Block {
	newBlock := &Block{
		Header: b.Header.DeepCopy(),
	}
	if len(b.Authorizations) != 0 {
		newBlock.Authorizations = make([]*Authorization, len(b.Authorizations))
		for i := range b.Authorizations {
			newBlock.Authorizations[i] = b.Authorizations[i].DeepCopy()
		}
	}
	return newBlock
}

func (h *Header) ToMap() *map[string]interface{} {
	hMap := make(map[string]interface{})
	hMap["prevHash"] = string(common.Expand(h.ParentHash[:]))
	hMap["authRoot"] = string(common.Expand(h.AuthRoot[:]))
	hMap["time"] = h.Timestamp
	return &hMap
}

func (b *Block) ToIndexMap() *map[string]interface{} {
	bMap := make(map[string]interface{})
	hMap := b.Header.ToMap()
	for k, v := range *hMap {
		bMap[k] = v
	}
	bMap["authNum"] = len(b.Authorizations)
	return &bMap
}

func (b *Block) ToMap() *map[string]interface{} {
	bMap := make(map[string]interface{})
	hMap := b.Header.ToMap()
	for k, v := range *hMap {
		bMap[k] = v
	}
	var authList []interface{}
	for _, auth := range b.Authorizations {
		authList = append(authList, auth.ToMap())
	}
	bMap["authList"] = authList
	bMap["authNum"] = len(authList)
	return &bMap
}

func (b *Block) ToPrettyJson() []byte {
	blockJson, err := json.Marshal(b)
	if err != nil {
		fmt.Println(err)
	}
	return common.PrettyPrintJson(blockJson)
}

func (b *Block) ToJson() []byte {
	blockJson, err := json.Marshal(b)
	if err != nil {
		fmt.Println(err)
	}
	return blockJson
}

func (h *Header) ToPrettyJson() []byte {
	blockJson, err := json.Marshal(h.ToMap())
	if err != nil {
		fmt.Println(err)
	}
	return blockJson
}

func (h *Header) ToJson() []byte {
	blockJson, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	return blockJson
}

func (b *Block) ToString() string {
	return string(b.ToPrettyJson())
}

func (h *Header) PrettyPrint() {
	fmt.Println(string(common.PrettyPrintJson(h.ToJson())))
}

func (h *Header) ToString() string {
	return string(h.ToJson())
}

func GenerateTestHeader() *Header {
	prev := cryptoType.RandHash()
	authRoot := cryptoType.RandHash()
	time := commonType.CurrentTime()
	h := NewHeader(prev, authRoot, time)
	return h
}

// generate random block for testing
func GenerateTestBlock() *Block {
	h := GenerateTestHeader()
	var authList []*Authorization
	authNum := rand.Int() % 5
	for i := 0; i < authNum; i += 1 {
		auth := GenerateTestAuthorization()
		authList = append(authList, auth)
	}
	b := NewBlock(h, authList)
	return b
}
