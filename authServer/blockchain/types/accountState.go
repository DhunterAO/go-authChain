package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DhunterAO/goAuthChain/common"
	commonType "github.com/DhunterAO/goAuthChain/common/types"
	cryptoType "github.com/DhunterAO/goAuthChain/crypto/types"
	types2 "github.com/DhunterAO/goAuthChain/crypto/types"
	"math/rand"
	"sync"
)

var (
	NoSuchAttribute = errors.New("no such attribute name")
	DurExtendLimit  = errors.New("authorization duration exceeds limit")
	ExistAttr       = errors.New("attribute added is existed")
	//TimeExtendLimit = errors.New("operation time exceeds limit")
)

type AccountState struct {
	Address    *cryptoType.Address
	Attributes map[string]*commonType.Attribute
	Nonce      uint64

	mu sync.RWMutex
}

func NewAccountState(address *cryptoType.Address, attributes []*commonType.Attribute, nonce uint64) *AccountState {
	newAs := &AccountState{
		Address:    address,
		Attributes: make(map[string]*commonType.Attribute),
		Nonce:      nonce,
	}
	if len(attributes) != 0 {
		for i := range attributes {
			name := attributes[i].Name
			newAs.Attributes[name] = attributes[i]
		}
	}
	return newAs
}

func (as *AccountState) ToString() string {
	return string(common.PrettyPrintJson(as.ToJson()))
}

func (as *AccountState) ToJson() []byte {
	asJson, err := json.Marshal(as.ToMap())
	if err != nil {
		fmt.Println(err)
	}
	return asJson
}

func (as *AccountState) ToMap() map[string]interface{} {
	asMap := make(map[string]interface{})
	asMap["address"] = as.Address.Hex()
	asMap["nonce"] = as.Nonce
	attributesMap := make([]interface{}, 0)
	for _, attribute := range as.Attributes {
		attributesMap = append(attributesMap, attribute.ToMap())
	}
	asMap["attributes"] = attributesMap
	return asMap
}

func JsonToAccountState(asJson []byte) *AccountState {
	var as AccountState
	err := json.Unmarshal(asJson, &as)
	if err != nil {
		fmt.Println(err)
	}
	return &as
}

func (as *AccountState) AddNonce() bool {
	if as.Nonce < common.MaxUint64 {
		as.Nonce += 1
		return true
	}
	return false
}

//func (as *AccountState) CheckAuthorizations(authorizations []*commonType.Authorization) bool {
//	for i := range authorizations {
//		if !as.CheckAuthorization(authorizations[i]) {
//			return false
//		}
//	}
//	return true
//}

func (as *AccountState) CheckAuthorization(auth *Authorization) (bool, error) {
	// check if the nonce in authorization matches the nonce in accountState
	if auth.Nonce != as.Nonce {
		return false, MisMatchNonce
	}

	// check if the attributes in authorization under permissions in accountState
	if flag, err := as.CheckAttributes(auth.Attributes); !flag {
		return false, err
	}
	return true, nil
}

func (as *AccountState) AddAttributes(attributes []*commonType.Attribute) bool {
	for i := range attributes {
		if !as.AddAttribute(attributes[i]) {
			return false
		}
	}
	return true
}

func (as *AccountState) AddAttribute(attribute *commonType.Attribute) bool {
	name := attribute.Name
	if _, ok := as.Attributes[name]; ok {
		return false
	}
	as.Attributes[name] = attribute
	return true
}

func (as *AccountState) CheckAttributes(attributes []*commonType.Attribute) (bool, error) {
	for i := range attributes {
		if flag, err := as.CheckAttribute(attributes[i]); !flag {
			return false, err
		}
	}
	return true, nil
}

func (as *AccountState) CheckAttribute(attribute *commonType.Attribute) (bool, error) {
	as.mu.RLock()
	defer as.mu.RUnlock()

	name := attribute.Name
	attribute, exist := as.Attributes[name]
	if !exist {
		return false, NoSuchAttribute
	}
	if attribute.Dur.Contain(attribute.Dur) {
		return true, nil
	}
	return false, DurExtendLimit
}

//func (as *AccountState) CheckOperations(operations []*operation.Operation) bool {
//	for i := range operations {
//		if !as.CheckOperation(operations[i]) {
//			return false
//		}
//	}
//	return true
//}
//
//func (as *AccountState) CheckOperation(operation *operation.Operation) bool {
//	as.mu.RLock()
//	defer as.mu.RUnlock()
//
//	if operation.Nonce != as.Nonce {
//		util.LogError(MisMatchNonce)
//		return false
//	}
//	return true
//}

func (as *AccountState) DeepCopy() *AccountState {
	as.mu.RLock()
	defer as.mu.RUnlock()

	newAs := &AccountState{
		Address:    as.Address,
		Attributes: make(map[string]*commonType.Attribute),
		Nonce:      as.Nonce,
	}
	if len(as.Attributes) != 0 {
		for name, attr := range as.Attributes {
			newAs.Attributes[name] = attr.DeepCopy()
		}
	}
	return newAs
}

func GenerateTestAccountState() *AccountState {
	as := &AccountState{
		Address:    types2.RandAddress(),
		Attributes: make(map[string]*commonType.Attribute),
		Nonce:      rand.Uint64(),
	}
	if !as.AddAttribute(commonType.GenerateTestAttribute()) {
		fmt.Println(ExistAttr)
	}
	return as
}
