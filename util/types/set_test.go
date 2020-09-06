package types

import (
	"fmt"
	"testing"
)

func TestNewIntSet(t *testing.T) {
	s := NewIntSet()
	s.Add(1)
	fmt.Println(s.Has(1))
	fmt.Println(s.Has(0))
	s.Remove(1)
	fmt.Println("---------")
	fmt.Println(s.Has(1))
	fmt.Println(s.Has(0))

}
