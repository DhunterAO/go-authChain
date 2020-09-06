package types

import (
	util2 "goauth/util"
	"math/big"
)

const (
	// SignatureLength is the expected length of the signature
	SignatureLength = 65
)

type Signature [SignatureLength]byte

var EmptySignature = Signature{}

func (sig Signature) ToECDSA() (*big.Int, *big.Int) {
	halfLen := SignatureLength / 2
	r := new(big.Int).SetBytes(sig[:halfLen])
	s := new(big.Int).SetBytes(sig[halfLen:])
	return r, s
}

func ECDSAtoSignature(r, s *big.Int) *Signature {
	var sig Signature
	rLen := len(r.Bytes())
	sLen := len(s.Bytes())
	copy(sig[SignatureLength/2-rLen:SignatureLength/2], r.Bytes())
	copy(sig[SignatureLength-sLen:], s.Bytes())
	return &sig
}

func (sig *Signature) Bytes() []byte {
	return sig[:]
}

func BytesToSignature(b []byte) *Signature {
	var sig Signature
	if len(b) > SignatureLength {
		b = b[:SignatureLength]
	}
	copy(sig[SignatureLength-len(b):], b)
	return &sig
}

func (sig *Signature) Hex() string {
	return string(util2.Expand(sig[:]))
}

func HexToSignature(hex string) *Signature {
	return BytesToSignature([]byte(hex))
}

func (sig *Signature) DeepCopy() *Signature {
	var newSig Signature
	copy(newSig[:], sig[:])
	return &newSig
}

func (sig *Signature) Pure() []byte {
	return sig[:SignatureLength-1]
}


// ValidateSignatureValues verifies whether the signature values are valid with
// the given chain rules. The v value is assumed to be either 0 or 1.
func ValidateSignatureValues(v byte, r, s *big.Int, homestead bool) bool {
	if r.Cmp(util2.Big1) < 0 || s.Cmp(util2.Big1) < 0 {
		return false
	}
	// reject upper range of s values (ECDSA malleability)
	// see discussion in secp256k1/libsecp256k1/include/secp256k1.h
	if homestead && s.Cmp(secp256k1halfN) > 0 {
		return false
	}
	// Frontier: allow s to be in full N range
	return r.Cmp(secp256k1N) < 0 && s.Cmp(secp256k1N) < 0 && (v == 0 || v == 1)
}
