package types

import (
	"encoding/json"
	"fmt"
)

type AuthList map[uint64]*Authorization

func NewAuthList() *AuthList {
	return &AuthList{}
}

func (al *AuthList) Add(auth *Authorization) bool {
	nonce := auth.Nonce
	auList := *al
	if _, ok := auList[nonce]; ok {
		return false
	}
	auList[nonce] = auth
	return true
}

func (al *AuthList) Remove(nonce uint64) bool {
	auList := *al
	if _, ok := auList[nonce]; ok {
		delete(auList, nonce)
		return true
	}
	return false
}

func (al *AuthList) ToMap() map[string]interface{} {
	alMap := make(map[string]interface{})
	for nonce, auth := range *al {
		alMap[string(nonce)] = auth.ToMap()
	}
	return alMap
}

func (al *AuthList) ToJson() []byte {
	alJson, err := json.Marshal(al.ToMap())
	if err != nil {
		fmt.Println(err)
	}
	return alJson
}

func (al *AuthList) ToString() string {
	return string(al.ToJson())
}

//type authList struct {
//	strict bool         // Whether nonces are strictly continuous or not
//	auths    *authSortedMap // Heap indexed sorted hash map of the transactions
//	lowNonce uint64
//}
//
//// newAuthList create a new transaction list for maintaining nonce-indexable fast,
//// gapped, sortable transaction lists.
//func NewAuthList(strict bool) *authList {
//	return &authList{
//		strict:  strict,
//		auths:     newAuthSortedMap(),
//	}
//}
//
//// nonceHeap is a heap.Interface implementation over 64bit unsigned integers for
//// retrieving sorted transactions from the possibly gapped future queue.
//type nonceHeap []uint64
//
//func (h nonceHeap) Len() int           { return len(h) }
//func (h nonceHeap) Less(i, j int) bool { return h[i] < h[j] }
//func (h nonceHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
//
//func (h *nonceHeap) Push(x interface{}) {
//	*h = append(*h, x.(uint64))
//}
//
//func (h *nonceHeap) Pop() interface{} {
//	old := *h
//	n := len(old)
//	x := old[n-1]
//	*h = old[0 : n-1]
//	return x
//}
//
//// authSortedMap is a nonce->Authorization hash map with a heap based index to allow
//// iterating over the contents in a nonce-incrementing way.
//type authSortedMap struct {
//	items map[uint64]*Authorization // Hash map storing the transaction data
//	index *nonceHeap                    // Heap of nonces of all the stored transactions (non-strict mode)
//}
//
//// newTxSortedMap creates a new nonce-sorted transaction map.
//func newAuthSortedMap() *authSortedMap {
//	return &authSortedMap{
//		items: make(map[uint64]*Authorization),
//		index: new(nonceHeap),
//	}
//}
//
//// Get retrieves the current transactions associated with the given nonce.
//func (m *authSortedMap) Get(nonce uint64) *Authorization {
//	return m.items[nonce]
//}
//
//// Put inserts a new transaction into the map, also updating the map's nonce
//// index. If a transaction already exists with the same nonce, it's overwritten.
//func (m *authSortedMap) Put(auth *Authorization) {
//	nonce := auth.Nonce
//	if m.items[nonce] == nil {
//		heap.Push(m.index, nonce)
//	}
//	m.items[nonce] = auth
//}
//
//// Remove deletes a transaction from the maintained map, returning whether the
//// transaction was found.
//func (m *authSortedMap) Remove(nonce uint64) bool {
//	// Short circuit if no transaction is present
//	_, ok := m.items[nonce]
//	if !ok {
//		return false
//	}
//	// Otherwise delete the transaction and fix the heap index
//	for i := 0; i < m.index.Len(); i++ {
//		if (*m.index)[i] == nonce {
//			heap.Remove(m.index, i)
//			break
//		}
//	}
//	delete(m.items, nonce)
//	return true
//}
//
//// Add tries to insert a new authorization into the list, returning whether the
//// authorization was accepted.
//func (l *authList) Add(auth *Authorization) bool {
//	fmt.Println(auth.Nonce)
//	old := l.auths.Get(auth.Nonce)
//	if old != nil {
//		fmt.Println("auth of the same nonce existed")
//		return false
//	}
//
//	l.auths.Put(auth)
//	if auth.Nonce < l.lowNonce {
//		l.lowNonce = auth.Nonce
//	}
//	return true
//}
//
//// Choose valid authorizations from authList whose nonces start at given nonce and consistent
//func (l *authList) ChooseAuthorizations(accountState *AccountState, limit uint64) []*Authorization {
//	// at most choose authorization of limit
//	chosenAuths := make([]*Authorization, limit)
//	nonce := accountState.Nonce
//	for i := nonce; i-nonce < limit; i += 1 {
//		auth := l.auths.Get(i)
//		if auth == nil {
//			return chosenAuths[:i-nonce]
//		}
//
//		ok, err := accountState.CheckAuthorization(auth)
//		if ! ok {
//			fmt.Println(err)
//			return chosenAuths[:i-nonce]
//		}
//		chosenAuths[i-nonce] = auth
//	}
//	return chosenAuths
//}
//
//// Remove deletes a transaction from the maintained list, returning whether the
//// transaction was found, and also returning any transaction invalidated due to
//// the deletion (strict mode only).
//func (l *authList) Remove(nonce uint64) bool {
//	if removed := l.auths.Remove(nonce); !removed {
//		return false
//	}
//	return true
//}
