package types

import (
	"encoding/json"
	"fmt"
	"github.com/DhunterAO/goAuthChain/common"
	"github.com/DhunterAO/goAuthChain/common/types"
	"github.com/DhunterAO/goAuthChain/crypto"
	cryptoType "github.com/DhunterAO/goAuthChain/crypto/types"
)

const AttributesLimit = 500

type Authorization struct {
	Recipient  *cryptoType.Address
	Attributes []*types.Attribute
	Nonce      uint64
	Signature  *cryptoType.Signature
}

func NewAuthorization(recipient *cryptoType.Address, attributes []*types.Attribute, nonce uint64) Authorization {
	auth := Authorization{
		Recipient: recipient,
		Nonce:     nonce,
	}
	if len(attributes) != 0 {
		auth.Attributes = make([]*types.Attribute, len(attributes))
		for i, attr := range attributes {
			auth.Attributes[i] = attr.DeepCopy()
		}
	}
	return auth
}

func NewAuthorizationFromJson(authJson []byte) *Authorization {
	auth := new(Authorization)
	err := json.Unmarshal(authJson, auth)
	if err != nil {
		panic(err)
	}
	return auth
}

func (auth *Authorization) AddAttribute(attr *types.Attribute) bool {
	if len(auth.Attributes) > AttributesLimit {
		return false
	}
	auth.Attributes = append(auth.Attributes, attr)
	return true
}

func GetSenderAddrFromAuth(auth *Authorization) (*cryptoType.Address, error) {
	hash := auth.CalcHash()
	sig := auth.Signature
	pk, err := crypto.RecoverPkFromSig(hash, sig)
	if err != nil {
		return &cryptoType.Address{}, err
	}
	return crypto.PubkeyToAddress(pk)
}

func (auth *Authorization) Bytes() []byte {
	authBytes := auth.Recipient[:]
	for _, attr := range auth.Attributes {
		authBytes = append(authBytes, attr.ToBytes()...)
	}
	authBytes = append(authBytes, common.Uint64ToBytes(auth.Nonce)...)
	return authBytes
}

func (auth *Authorization) ToMap() map[string]interface{} {
	authMap := make(map[string]interface{})
	authMap["to"] = auth.Recipient.Hex()
	authMap["nonce"] = auth.Nonce
	authMap["signature"] = auth.Signature.Hex()

	attributes := make([]map[string]interface{}, len(auth.Attributes))
	for i := range auth.Attributes {
		attributes[i] = auth.Attributes[i].ToMap()
	}
	authMap["attributes"] = attributes
	return authMap
}

func (auth *Authorization) ToJson() []byte {
	authJson, err := json.Marshal(auth)
	if err != nil {
		fmt.Println(err)
	}
	return authJson
}

func (auth *Authorization) ToPrettyJson() []byte {
	authJson, err := json.Marshal(auth.ToMap())
	if err != nil {
		fmt.Println(err)
	}
	return authJson
}

func (auth *Authorization) ToString() string {
	return string(common.PrettyPrintJson(auth.ToPrettyJson()))
}

func (auth *Authorization) CalcHash() *cryptoType.Hash {
	return crypto.CalcHash(auth.Bytes())
}

func (auth *Authorization) Sign(signature *cryptoType.Signature) {
	auth.Signature = signature
}

func (auth *Authorization) CheckSignature() bool {
	hash := auth.CalcHash()
	pub, err := crypto.RecoverPkFromSig(hash, auth.Signature)
	if err != nil {
		return false
	}
	return crypto.VerifySignature(pub, hash, auth.Signature)
}

func (auth *Authorization) DeepCopy() *Authorization {
	newAuth := &Authorization{
		Recipient: auth.Recipient,
		Nonce:     auth.Nonce,
		Signature: auth.Signature.DeepCopy(),
	}
	if len(auth.Attributes) != 0 {
		for name, attr := range auth.Attributes {
			newAuth.Attributes[name] = attr.DeepCopy()
		}
	}
	return newAuth
}

// generate fake authorization used for test
func GenerateTestAuthorization() *Authorization {
	sk, _ := cryptoType.GenerateKey()
	pk := sk.PublicKey

	//from := crypto.PubkeyToAddress(pk)
	to := cryptoType.EcdsaPubkeyToAddress(pk)

	attr := types.GenerateTestAttribute()
	attrs := []*types.Attribute{attr}
	nonce := uint64(0)
	auth := NewAuthorization(to, attrs, nonce)

	sig, _ := crypto.Sign(auth.CalcHash(), sk)
	auth.Sign(sig)
	return &auth
}
