package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DhunterAO/goAuthChain/crypto"
	cryptoTypes "github.com/DhunterAO/goAuthChain/crypto/types"
	"github.com/DhunterAO/goAuthChain/log"
	"sync"
)

var (
	nonceExist = errors.New("the nonce of auth is existed in auth list")
)

type AuthPool struct {
	Pending map[cryptoTypes.Address]*AuthList `json:"pending"`

	mu  sync.RWMutex `json:"-"`
	log *log.Logger  `json:"-"`
}

func NewAuthPool(pendingAuths []*Authorization) (*AuthPool, error) {
	authPool := new(AuthPool)
	authPool.Pending = make(map[cryptoTypes.Address]*AuthList)
	for _, auth := range pendingAuths {
		_, err := authPool.AddAuthorization(auth)
		if err != nil {
			return &AuthPool{}, err
		}
	}
	return authPool, nil
}

// add authorization into authPool, exactly into the authList in the pool
func (authPool *AuthPool) AddAuthorization(auth *Authorization) (bool, error) {
	authPool.mu.Lock()
	defer authPool.mu.Unlock()

	pk, err := crypto.RecoverPkFromSig(auth.CalcHash(), auth.Signature)
	if err != nil {
		return false, err
	}
	from, err := crypto.PubkeyToAddress(pk)
	if err != nil {
		return false, err
	}
	if authL, ok := authPool.Pending[*from]; ok {
		//fmt.Println("exist")
		return authL.Add(auth), nil
	} else {
		//fmt.Println("not exist")
		authPool.Pending[*from] = NewAuthList()
		if authPool.Pending[*from].Add(auth) {
			//s, _ := json.Marshal(authPool)
			//fmt.Println("added ", authPool, s)
			return true, nil
		}
		return false, nonceExist
	}
}

func (authPool *AuthPool) DeleteAuths(auths []*Authorization) {
	for i := range auths {
		if err := authPool.DeleteAuth(auths[i]); err != nil {
			authPool.log.Error(err.Error())
		}
	}
	return
}

func (authPool *AuthPool) DeleteAuth(auth *Authorization) error {
	from, err := GetSenderAddrFromAuth(auth)
	if err != nil {
		return err
	}
	if _, ok := authPool.Pending[*from]; !ok {
		return errors.New("the sender address is not in the pool")
	}
	authList := authPool.Pending[*from]
	nonce := auth.Nonce
	if authList.Remove(nonce) == false {
		return errors.New("the nonce of authorization is not in the pool")
	}
	return nil
}

func (authPool *AuthPool) ToJson() []byte {
	authPoolJson, err := json.Marshal(authPool)
	if err != nil {
		authPool.log.Error(err.Error())
	}
	return authPoolJson
}

func (authPool *AuthPool) ToString() string {
	return string(authPool.ToJson())
}

func JsonToAuthPool(b []byte) *AuthPool {
	authPool := new(AuthPool)
	err := json.Unmarshal(b, authPool)
	if err != nil {
		fmt.Println(err)
	}
	return authPool
}

func (authPool *AuthPool) ChooseAuthorizations(limit uint64, state *GlobalState) []*Authorization {
	authPool.mu.RLock()
	defer authPool.mu.RUnlock()

	chosenAuth := make([]*Authorization, 0)
	virState := state.DeepCopy()
	for address, al := range authPool.Pending {
		nonce := virState.GetAccountNonce(&address)
		authList := *al
		for ; limit > 0 && authList[nonce] != nil; nonce += 1 {
			if flag, err := virState.CheckAuthorization(authList[nonce]); flag {
				chosenAuth = append(chosenAuth, authList[nonce])
				limit -= 1
				if _, err := virState.AcceptAuthorization(authList[nonce]); err != nil {
					authPool.log.Error(err.Error())
				}
			} else {
				if err != nil {
					fmt.Println(err.Error())
				}
				break
			}
		}
		if limit == 0 {
			return chosenAuth
		}
	}
	fmt.Println(chosenAuth)
	return chosenAuth
}
