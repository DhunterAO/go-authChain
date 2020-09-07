package types

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/DhunterAO/goAuthChain/common/hexutil"
	"reflect"
	"testing"
)

func TestEcdsaPubkeyToAddress(t *testing.T) {
	sk, _ := GenerateKey()
	pk := sk.PublicKey
	fmt.Println(pk)
	addr := EcdsaPubkeyToAddress(pk)
	fmt.Println(addr.Hex())
	pkCom := CompressPubkey(&pk)
	fmt.Println(pkCom)
	pubkey := BytesToPubkey(pkCom)
	fmt.Println(pubkey)
	pk2, _ := DecompressPubkey(pubkey.Bytes())
	fmt.Println(addr.Hex() == EcdsaPubkeyToAddress(*pk2).Hex())
}

func TestUnmarshalPubkey(t *testing.T) {
	key, err := UnmarshalPubkey(nil)
	if err != errInvalidPubkey || key != nil {
		t.Fatalf("expected error, got %v, %v", err, key)
	}
	key, err = UnmarshalPubkey([]byte{1, 2, 3})
	if err != errInvalidPubkey || key != nil {
		t.Fatalf("expected error, got %v, %v", err, key)
	}

	var (
		enc, _ = hex.DecodeString("04760c4460e5336ac9bbd87952a3c7ec4363fc0a97bd31c86430806e287b437fd1b01abc6e1db640cf3106b520344af1d58b00b57823db3e1407cbc433e1b6d04d")
		dec    = &ecdsa.PublicKey{
			Curve: S256(),
			X:     hexutil.MustDecodeBig("0x760c4460e5336ac9bbd87952a3c7ec4363fc0a97bd31c86430806e287b437fd1"),
			Y:     hexutil.MustDecodeBig("0xb01abc6e1db640cf3106b520344af1d58b00b57823db3e1407cbc433e1b6d04d"),
		}
	)
	key, err = UnmarshalPubkey(enc)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(key, dec) {
		t.Fatal("wrong result")
	}
}

func TestNewContractAddress(t *testing.T) {
	key, _ := HexToECDSA(testPrivHex)
	addr := HexToAddress(testAddrHex)
	genAddr := EcdsaPubkeyToAddress(key.PublicKey)
	// sanity check before using addr to create contract address
	checkAddr(t, *genAddr, *addr)

	caddr0 := CreateAddress(*addr, 0)
	caddr1 := CreateAddress(*addr, 1)
	caddr2 := CreateAddress(*addr, 2)
	checkAddr(t, *HexToAddress("333c3310824b7c685133f2bedb2ca4b8b4df633d"), *caddr0)
	checkAddr(t, *HexToAddress("8bda78331c916a08481428e4b07c96d3e916d165"), *caddr1)
	checkAddr(t, *HexToAddress("c9ddedf451bc62ce88bf9292afb13df35b670699"), *caddr2)
}

func checkAddr(t *testing.T, addr0, addr1 Address) {
	if addr0 != addr1 {
		t.Fatalf("address mismatch: want: %x have: %x", addr0, addr1)
	}
}
