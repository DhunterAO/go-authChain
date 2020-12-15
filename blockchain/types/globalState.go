package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DhunterAO/goAuthChain/common"
	commonType "github.com/DhunterAO/goAuthChain/common/types"
	cryptoType "github.com/DhunterAO/goAuthChain/crypto/types"
	"reflect"
	"strings"
	"sync"
)

var (
	UnKnownAddress = errors.New("unknown address in global state")
	MisMatchNonce  = errors.New("nonce mismatch between authorization/operation and account state")
	//NoPermission      = errors.New("authorization/operation without permission")
	NoPermissionForGm = errors.New("no permission of attribute with prefix of gm")
)

type GlobalState struct {
	GmList map[cryptoType.Address]string
	States map[cryptoType.Address]*AccountState
	mu     sync.RWMutex
}

func NewGlobalState(gmList map[cryptoType.Address]string, preStates []*AccountState) *GlobalState {
	state := &GlobalState{
		GmList: gmList,
		States: make(map[cryptoType.Address]*AccountState),
	}
	for _, aState := range preStates {
		state.AddState(aState)
	}
	return state
}

func (gState *GlobalState) AddState(aState *AccountState) bool {
	gState.mu.Lock()
	defer gState.mu.Unlock()

	_, exists := gState.States[*aState.Address]
	if !exists {
		gState.States[*aState.Address] = aState
		return true
	}
	return false
}

func (gState *GlobalState) GetAccountNonce(addr *cryptoType.Address) uint64 {
	gState.mu.RLock()
	defer gState.mu.RUnlock()

	state, exists := gState.States[*addr]
	if exists {
		return state.Nonce
	}
	return common.MaxUint64
}

func (gState *GlobalState) GetAccountState(address *cryptoType.Address) *AccountState {
	gState.mu.RLock()
	defer gState.mu.RUnlock()

	if state, ok := gState.States[*address]; ok {
		return state
	}
	return nil
}

func (gState *GlobalState) GetAttrsFromAddr(address *cryptoType.Address) (bool, []string) {
	if attr, ok := gState.GmList[*address]; ok {
		fmt.Println("attr in gmlist")
		return true, []string{attr}
	}
	as := gState.GetAccountState(address)
	var attrs []string
	for _, attr := range as.Attributes {
		if attr.Dur.ContainTime(commonType.CurrentTime()) {
			attrs = append(attrs, attr.Name)
		}
	}
	return false, attrs
}

func (gState *GlobalState) CheckAuthorizations(auths []*Authorization) bool {
	virGs := gState.DeepCopy()
	for i := range auths {
		if flag, _ := virGs.CheckAuthorization(auths[i]); !flag {
			return false
		} else {
			if res, err := virGs.AcceptAuthorization(auths[i]); err != nil || res != true {
				return false
			}
		}
	}
	return true
}

// check authorization from current global state, return answer and error
func (gState *GlobalState) CheckAuthorization(auth *Authorization) (bool, error) {
	fmt.Println(gState)
	// recover address from signature of authorization
	addr, err := GetSenderAddrFromAuth(auth)
	if err != nil {
		return false, nil
	}

	// verify whether the address is a gm address
	// if the address is a gm, check if the attributes in the auth has the prefix of the gm
	gmList := gState.GmList
	fmt.Println("gmlist", gmList)
	if prefix, ok := gmList[*addr]; ok {
		fmt.Println("gm", *addr, prefix)
		for _, attr := range auth.Attributes {
			fmt.Println(reflect.TypeOf(attr), attr)
			if !strings.HasPrefix(attr.Name, prefix) {
				return false, NoPermissionForGm
			}
		}
		return true, nil
	} else {
		fmt.Println("not gm")
		// if the address is not a gm, check the authorization with the accountState
		accountState := gState.GetAccountState(addr)
		// if the address not exist in global state, return UnKnownAddress
		if accountState == nil {
			return false, UnKnownAddress
		}
		// return the check from account state
		return accountState.CheckAuthorization(auth)
	}
}

func (gState *GlobalState) AcceptAuthorizations(auths []*Authorization) bool {
	for i := range auths {
		if flag, _ := gState.AcceptAuthorization(auths[i]); !flag {
			return false
		}
	}
	return true
}

func (gState *GlobalState) AcceptAuthorization(auth *Authorization) (bool, error) {
	from, err := GetSenderAddrFromAuth(auth)
	if err != nil {
		return false, err
	}
	if !gState.States[*from].AddNonce() {
		return false, err
	}

	to := auth.Recipient
	if _, ok := gState.States[*to]; ok {
		return gState.States[*to].AddAttributes(auth.Attributes), nil
	} else {
		gState.States[*to] = NewAccountState(to, auth.Attributes, 0)
	}
	return true, nil
}

func (gState *GlobalState) CalcHash() *cryptoType.Hash {
	hash := cryptoType.Hash{}

	return &hash
}

func (gState *GlobalState) DeepCopy() *GlobalState {
	newGs := &GlobalState{
		GmList: make(map[cryptoType.Address]string),
		States: make(map[cryptoType.Address]*AccountState),
		mu:     sync.RWMutex{},
	}
	if len(gState.States) != 0 {
		for address, as := range gState.States {
			newGs.States[address] = as.DeepCopy()
		}
	}
	if len(gState.GmList) != 0 {
		for address, attr := range gState.GmList {
			newGs.GmList[address] = attr
		}
	}
	return newGs
}

func (gState *GlobalState) ToMap() map[string]interface{} {
	gsMap := make(map[string]interface{})
	states := make(map[string]interface{})
	for address, as := range gState.States {
		states[address.Hex()] = as.ToMap()
	}
	gsMap["states"] = states
	return gsMap
}

func (gState *GlobalState) ToJson() []byte {
	gsJson, err := json.Marshal(gState.ToMap())
	if err != nil {
		fmt.Println(err)
	}
	return gsJson
}

func (gState *GlobalState) ToString() string {
	return string(common.PrettyPrintJson(gState.ToJson()))
}
