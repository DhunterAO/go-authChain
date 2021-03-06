package types

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	types2 "github.com/DhunterAO/goAuthChain/common/types"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadECDSAFile(t *testing.T) {
	keyBytes := types2.FromHex(testPrivHex)
	fileName0 := "test_key0"
	fileName1 := "test_key1"
	checkKey := func(k *ecdsa.PrivateKey) {
		checkAddr(t, *EcdsaPubkeyToAddress(k.PublicKey), *HexToAddress(testAddrHex))
		loadedKeyBytes := FromECDSA(k)
		if !bytes.Equal(loadedKeyBytes, keyBytes) {
			t.Fatalf("private key mismatch: want: %x have: %x", keyBytes, loadedKeyBytes)
		}
	}

	ioutil.WriteFile(fileName0, []byte(testPrivHex), 0600)
	defer os.Remove(fileName0)

	key0, err := LoadECDSA(fileName0)
	if err != nil {
		t.Fatal(err)
	}
	checkKey(key0)

	// again, this time with SaveECDSA instead of manual save:
	err = SaveECDSA(fileName1, key0)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fileName1)

	key1, err := LoadECDSA(fileName1)
	if err != nil {
		t.Fatal(err)
	}
	checkKey(key1)
}

func TestSaveAndLoad(t *testing.T) {
	fileName := "./test.key"
	_, err := os.Create(fileName)
	defer os.Remove(fileName)
	if err != nil {
		fmt.Println(err.Error())
	}

	sk, _ := GenerateKey()
	err = SaveECDSA(fileName, sk)
	if err != nil {
		t.Fatal(err)
	}

	sk2, err := LoadECDSA(fileName)
	if err != nil {
		t.Fatal(err)
	}

	if *BytesToPubkey(CompressPubkey(&sk.PublicKey)) != *BytesToPubkey(CompressPubkey(&sk2.PublicKey)) {
		t.Log(BytesToPubkey(CompressPubkey(&sk.PublicKey)))
		t.Log(BytesToPubkey(CompressPubkey(&sk2.PublicKey)))
	}
}
