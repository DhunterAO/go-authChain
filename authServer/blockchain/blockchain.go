package blockchain

import (
	"encoding/json"
	"errors"
	blockType "github.com/DhunterAO/goAuthChain/authServer/blockchain/types"
	"github.com/DhunterAO/goAuthChain/common"
	"github.com/DhunterAO/goAuthChain/common/fileutil"
	commonType "github.com/DhunterAO/goAuthChain/common/types"
	cryptoTypes "github.com/DhunterAO/goAuthChain/crypto/types"
	"github.com/DhunterAO/goAuthChain/log"
	"github.com/syndtr/goleveldb/leveldb"
	"sync"
)

const (
	AuthorizationLimit = 500
)

var (
	HashMisMatch  = errors.New("mismatch hash")
	InvalidAuths  = errors.New("invalid authorizations")
	ChainPathLack = errors.New("the parameter chainPath must be needed")
)

type Blockchain struct {
	// blocks in the blockchain
	Blocks []*blockType.Block
	// global state used for authentication
	State *blockType.GlobalState
	// pending authorizations and operations stored in pool
	AuthPool *blockType.AuthPool
	// chainDB used for blockchain data storage
	ChainDB *leveldb.DB
	// log used for logging
	Log log.Logger
	// global mutex for chain lock
	mu sync.RWMutex
	// quit channel used for terminal
	quit chan struct{}
}

func InitBlockchain(chainConfPath string, logger *log.Logger) (*Blockchain, error) {
	// load all parameters into cfg from file at cfgPath
	var cfg map[string]interface{}
	err := fileutil.LoadJson(chainConfPath, &cfg)
	if err != nil {
		logger.Error("load blockchain configuration failed !")
		return nil, err
	}
	//cfgJson, _ := json.Marshal(cfg)
	//fmt.Println(string(cfgJson))

	// check parameters in cfg including dataPath, chainPath, gmList, preAuths
	chainPath, ok := cfg["chainPath"].(string)
	if !ok {
		return nil, ChainPathLack
	}
	//fmt.Println(chainPath)
	chainDB, err := leveldb.OpenFile(chainPath, nil)

	gmMap, ok := cfg["gmList"].(map[string]interface{})
	if !ok {
		logger.Info("gmList not exist")
		gmMap = make(map[string]interface{})
	}
	gmList := make(map[cryptoTypes.Address]string)
	for k, v := range gmMap {
		gmList[*cryptoTypes.HexToAddress(k)] = v.(string)
	}
	//fmt.Println(gmList)

	preState := &blockType.GlobalState{
		GmList: gmList,
		States: map[cryptoTypes.Address]*blockType.AccountState{},
	}
	for addr := range gmList {
		preState.States[addr] = blockType.NewAccountState(&addr, []*commonType.Attribute{}, 0)
	}
	//fmt.Println(preState)

	preAuths, ok := cfg["preAuths"].([]*blockType.Authorization)
	if !ok {
		preAuths = []*blockType.Authorization{}
	}

	header := &blockType.Header{
		ParentHash: &cryptoTypes.Hash{},
		AuthRoot:   nil,
		//StateRoot:  &types.Hash{},
		Timestamp: 0,
	}
	genesisBlock := &blockType.Block{
		Header:         header,
		Authorizations: []*blockType.Authorization{},
		//State: 			preState,
	}
	genesisBlock.SetRoot()

	authPool, err := blockType.NewAuthPool(preAuths)
	if err != nil {
		return nil, err
	}

	bc := &Blockchain{
		Blocks:   []*blockType.Block{genesisBlock},
		State:    preState,
		AuthPool: authPool,
		ChainDB:  chainDB,
		quit:     make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-bc.quit:
				err = chainDB.Close()
				if err != nil {
					bc.Log.Error(err.Error())
				}
			default:
			}
		}
	}()
	logger.Info("init blockchain success")
	logger.Info(string(bc.ToJson()))
	return bc, nil
}

func (bc *Blockchain) GetHeight() int {
	return len(bc.Blocks)
}

func (bc *Blockchain) GetLastBlock() *blockType.Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *Blockchain) GetLastHash() *cryptoTypes.Hash {
	lastBlock := bc.GetLastBlock()
	return lastBlock.CalcHash()
}

func (bc *Blockchain) GenerateNewBlock() *blockType.Block {
	chosenAuthorizations := bc.AuthPool.ChooseAuthorizations(AuthorizationLimit, bc.State)
	newHeader := &blockType.Header{
		ParentHash: bc.GetLastHash(),
		AuthRoot:   nil,
		//StateRoot:  nil,
		Timestamp: commonType.CurrentTime(),
	}
	newBlock := blockType.Block{
		Header:         newHeader,
		Authorizations: chosenAuthorizations,
		//State:          nil,
	}
	newBlock.SetRoot()
	return &newBlock
}

func (bc *Blockchain) CheckNewBlock(b *blockType.Block) bool {
	if *b.Header.ParentHash != *bc.GetLastHash() {
		bc.Log.Error(HashMisMatch.Error())
		return false
	}

	vState := bc.State.DeepCopy()
	if !vState.CheckAuthorizations(b.Authorizations) {
		bc.Log.Error(InvalidAuths.Error())
		return false
	}
	return true
}

func (bc *Blockchain) AddNewBlock(b *blockType.Block) (bool, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if !bc.State.AcceptAuthorizations(b.Authorizations) {
		return false, InvalidAuths
	}
	bc.AuthPool.DeleteAuths(b.Authorizations)
	bc.Blocks = append(bc.Blocks, b)
	return true, nil
}

func (bc *Blockchain) ToMap() *map[string]interface{} {
	bcMap := make(map[string]interface{})
	//bcMap["genesis"] = bc.Genesis.ToMap()
	blocks := make([]*map[string]interface{}, len(bc.Blocks))
	for i, block := range bc.Blocks {
		blocks[i] = block.ToMap()
	}
	bcMap["blocks"] = blocks
	return &bcMap
}

func (bc *Blockchain) ToIndexMap() *map[string]interface{} {
	bcMap := make(map[string]interface{})
	blocks := make([]*map[string]interface{}, len(bc.Blocks))
	for i, block := range bc.Blocks {
		blocks[i] = block.ToIndexMap()
	}
	bcMap["blocks"] = blocks
	return &bcMap
}

func (bc *Blockchain) ToIndexJson() []byte {
	bcJson, err := json.Marshal(bc.ToIndexMap())
	if err != nil {
		bc.Log.Error(err.Error())
	}
	return common.PrettyPrintJson(bcJson)
}

func (bc *Blockchain) ToPrettyJson() []byte {
	bcJson, err := json.Marshal(bc)
	if err != nil {
		bc.Log.Error(err.Error())
	}
	return common.PrettyPrintJson(bcJson)
}

func (bc *Blockchain) ToJson() []byte {
	bcJson, err := json.Marshal(bc)
	if err != nil {
		bc.Log.Error(err.Error())
	}
	return bcJson
}

func (bc *Blockchain) ToString() string {
	return string(bc.ToJson())
}

func (bc *Blockchain) Quit() {
	close(bc.quit)
}
