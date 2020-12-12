package blockchain

import (
	blockType "github.com/DhunterAO/goAuthChain/authServer/blockchain/types"
	commonType "github.com/DhunterAO/goAuthChain/common/types"
	"github.com/DhunterAO/goAuthChain/crypto"
	"github.com/DhunterAO/goAuthChain/crypto/types"
	"github.com/DhunterAO/goAuthChain/log"
	"os"
	"testing"
)

func TestEmptyBlockchain(t *testing.T) {
	// generate blockchain using config
	var logger log.Logger
	bc, err := InitBlockchain("../../data/conf/blockchain.conf", &logger)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log("blockchain:\n", string(bc.ToPrettyJson()))

	// generate new block from the blockchain
	b := bc.GenerateNewBlock()
	t.Log("generate new block\n", string(b.ToPrettyJson()))

	// validate the new block generated
	t.Log("check new block ", bc.CheckNewBlock(b))

	// add new block into the blockchain
	if res, err := bc.AddNewBlock(b); err != nil {
		t.Error(err, res)
	}
	t.Log("new blockchain:\n", string(bc.ToPrettyJson()))
}

func TestBlockchainWithAuth(t *testing.T) {
	t.Log("1--------initial new blockchain-------------")
	var logger log.Logger
	configPath := "../../data/conf/blockchain.conf"
	if _, err := os.Stat(configPath); err != nil {
		os.Create(configPath)
	}
	bc, err := InitBlockchain(configPath, &logger)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(string(bc.ToPrettyJson()))

	t.Log("2----------generate test authorization from gm-----------")
	gmSkHex := "0000000000000000000000000000000000000000000000000000000000001979"
	privKey, _ := types.HexToECDSA(gmSkHex)
	auth := &blockType.Authorization{
		Recipient:  types.RandAddress(),
		Attributes: []*commonType.Attribute{},
		Nonce:      0,
		Signature:  nil,
	}
	sig, _ := crypto.Sign(auth.CalcHash(), privKey)
	auth.Sign(sig)
	t.Log("print new authorization\n", auth.ToString())

	t.Log("4----------add authorization into AuthorizationPool-----------")
	flag, err := bc.AuthPool.AddAuthorization(auth)
	t.Log("add result: ", flag, err)
	t.Log("authPool after insertion\n", bc.AuthPool.ToString())

	t.Log("5----------generate new block-----------")
	b := bc.GenerateNewBlock()
	t.Log("new block", string(b.ToJson()))

	t.Log("6----------add new block into blockchain-----------")
	if res, err := bc.AddNewBlock(b); err != nil {
		t.Error(err, res)
	}
	t.Log("print new blockchain\n", bc.ToString())

	t.Log("7----------generate and add another new block into blockchain-----------")
	b = bc.GenerateNewBlock()
	if res, err := bc.AddNewBlock(b); err != nil {
		t.Error(err, res)
	}
	t.Log("new blockchain\n", string(bc.ToJson()))
}
