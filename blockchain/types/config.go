package types

import (
	"github.com/DhunterAO/goAuthChain/crypto/types"
)

type Config struct {
	GmList   map[types.Address]string
	PreAuths []*Authorization
	// the paths of db storing data and blockchain
	DataDB  string
	ChainDB string
}

func GenerateConfig() Config {
	var gmList = map[types.Address]string{}
	gmAddr := types.HexToAddress("0x2266Cf3af337dd95f42352160490D9ad259204c2")
	gmList[*gmAddr] = "B"
	cfg := Config{
		GmList:   gmList,
		PreAuths: []*Authorization{},
		DataDB:   "testDataDB/",
		ChainDB:  "testChainDB/",
	}
	return cfg
}
