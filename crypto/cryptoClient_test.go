package crypto

import (
	"fmt"
	"goauth/crypto/types"
	"testing"
)

func TestSignAndVerify(t *testing.T) {
	sk, _ := types.GenerateKey()
	msg := []byte("hello world")
	hash := CalcHash(msg)
	sig, _ := Sign(hash, sk)
	pub, _ := types.Ecrecover(hash.Bytes(), sig.Bytes())
	pubkey, _ := types.UnmarshalPubkey(pub)
	pk := types.CompressPubkey(pubkey)
	t.Log(len(pk), pk)
	pub2 := types.CompressPubkey(&sk.PublicKey)
	t.Log(len(pub2), pub2)
	t.Log(types.VerifySignature(pk, hash.Bytes(), sig.Bytes()[:len(sig.Bytes())-1]))
	t.Log(VerifySignature(types.BytesToPubkey(pk), hash, sig))
	t.Log(types.VerifySignature(pub2, hash.Bytes(), sig.Bytes()[:len(sig.Bytes())-1]))
	t.Log(VerifySignature(types.BytesToPubkey(pub2), hash, sig))

	pub3, _ := RecoverPkFromSig(hash, sig)
	t.Log(VerifySignature(pub3, hash, sig))
}

func TestGenerateKey(t *testing.T) {
	sk, _ := types.LoadECDSA("../data/keys/student.key")
	addr, _ := PubkeyToAddress(types.BytesToPubkey(types.CompressPubkey(&sk.PublicKey)))
	fmt.Println(addr.Hex())
}
