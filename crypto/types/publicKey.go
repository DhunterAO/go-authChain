package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
)

const (
	CompressedPubkeyLen = 33
)

var errInvalidPubkey = errors.New("invalid secp256k1 public key")

type Pubkey [CompressedPubkeyLen]byte

func (pk Pubkey) Bytes() []byte {
	return pk[:]
}

func BytesToPubkey(b []byte) *Pubkey {
	var pk Pubkey
	pk.SetBytes(b)
	return &pk
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (pk *Pubkey) SetBytes(b []byte) {
	if len(b) > len(pk) {
		b = b[len(b)-CompressedPubkeyLen:]
	}
	copy(pk[CompressedPubkeyLen-len(b):], b)
}

func RecoverEcdsaFromSig(hash *Hash, sig *Signature) (*ecdsa.PublicKey, error) {
	uncompressedPk, err := Ecrecover(hash.Bytes(), sig.Bytes())
	if err != nil {
		return nil, err
	}
	pk, err := UnmarshalPubkey(uncompressedPk)
	if err != nil {
		return nil, err
	}
	return pk, err
}

func RecoverPkFromSig(hash *Hash, sig *Signature) (*Pubkey, error) {
	pk, err := RecoverEcdsaFromSig(hash, sig)
	if err != nil {
		return &Pubkey{}, err
	}
	return BytesToPubkey(CompressPubkey(pk)), err
}

func (pk *Pubkey)Address() (*Address, error) {
	ecdsaPk, err := DecompressPubkey(pk.Bytes())
	if err != nil {
		return nil, err
	}
	addr := EcdsaPubkeyToAddress(*ecdsaPk)
	return addr, nil
}

// UnmarshalPubkey converts bytes to a secp256k1 public key.
func UnmarshalPubkey(pub []byte) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(S256(), pub)
	if x == nil {
		return nil, errInvalidPubkey
	}
	return &ecdsa.PublicKey{Curve: S256(), X: x, Y: y}, nil
}

func FromECDSAPub(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(S256(), pub.X, pub.Y)
}