package types

import "sync"

type Uint64Set struct {
	m map[uint64]bool
	s sync.RWMutex
}

func NewUint64Set() *Uint64Set {
	return &Uint64Set{
		m: map[uint64]bool{},
	}
}

func (uint64Set *Uint64Set) Add(x uint64) {
	uint64Set.s.Lock()
	defer uint64Set.s.Unlock()

	uint64Set.m[x] = true
}

func (uint64Set *Uint64Set) Remove(x uint64) {
	uint64Set.s.Lock()
	defer uint64Set.s.Unlock()
	delete(uint64Set.m, x)
}

func (uint64Set *Uint64Set) Has(item uint64) bool {
	uint64Set.s.RLock()
	defer uint64Set.s.RUnlock()

	_, ok := uint64Set.m[item]
	return ok
}

func (uint64Set *Uint64Set) Len() uint64 {
	return uint64(len(uint64Set.m))
}

func (uint64Set *Uint64Set) Clear() {
	uint64Set.s.RLock()
	defer uint64Set.s.RUnlock()

	uint64Set.m = make(map[uint64]bool)
}

type IntSet struct {
	m map[int]bool
	s sync.RWMutex
}

func NewIntSet() *IntSet {
	return &IntSet{
		m: map[int]bool{},
	}
}

func (intSet *IntSet) Add(x int) {
	intSet.s.Lock()
	defer intSet.s.Unlock()

	intSet.m[x] = true
}

func (intSet *IntSet) Remove(x int) {
	intSet.s.Lock()
	defer intSet.s.Unlock()

	delete(intSet.m, x)
}

func (intSet *IntSet) Has(item int) bool {
	intSet.s.RLock()
	defer intSet.s.RUnlock()

	_, ok := intSet.m[item]
	return ok
}
