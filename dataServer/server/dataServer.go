package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DhunterAO/goAuthChain/blockchain"
	"github.com/DhunterAO/goAuthChain/common/fileutil"
	commonType "github.com/DhunterAO/goAuthChain/common/types"
	"github.com/DhunterAO/goAuthChain/dataServer/acl"
	"github.com/DhunterAO/goAuthChain/dataServer/db"
	dataServerType "github.com/DhunterAO/goAuthChain/dataServer/server/types"
	"github.com/DhunterAO/goAuthChain/log"
	"github.com/kataras/iris"
	"github.com/syndtr/goleveldb/leveldb"
	"math"
)

var DataPathLack = errors.New("the parameter dataPath must be needed")

type DataServer struct {
	dataDB *leveldb.DB
	Acl    *acl.PolicyList
	Bc     *blockchain.Blockchain
	App    *iris.Application
	Port   string

	logger *log.Logger
}

var dataServer DataServer

func InitDataServer(dataServerConfigPath string, bc *blockchain.Blockchain, logger *log.Logger) (*DataServer, error) {
	fmt.Println(dataServerConfigPath, bc, logger)
	// load all parameters into cfg from file at cfgPath
	var cfg map[string]interface{}
	err := fileutil.LoadJson(dataServerConfigPath, &cfg)
	if err != nil {
		logger.Error("load data server configuration failed !")
		return nil, err
	}
	// check parameters in cfg including nodeID, candidateInfo, keyPath, host, port
	port, ok := cfg["port"].(string)
	if !ok {
		port = "8001"
	}

	dbPath, ok := cfg["dbPath"].(string)
	if !ok {
		dbPath = "../data/dataDB/"
	}

	aclPath, ok := cfg["aclPath"].(string)
	if !ok {
		aclPath = "../data/acl/abac.list"
	}
	acList, err := LoadAcList(aclPath)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	dataServer.Acl = acList
	dataServer.App = iris.New()
	dataServer.Bc = bc
	dataServer.Port = port
	dataDB, err := db.OpenDB(dbPath)
	if err != nil {
		return nil, err
	}
	dataServer.dataDB = dataDB
	addUrls(dataServer.App)
	return &dataServer, nil
}

func (dt *DataServer) Start() {
	err := dt.App.Run(iris.Addr(":" + dt.Port))
	if err != nil {
		fmt.Println("app start err")
	}
}

func LoadAcList(filePath string) (*acl.PolicyList, error) {
	var cfg []interface{}
	err := fileutil.LoadJson(filePath, &cfg)
	if err != nil {
		return nil, err
	}
	fmt.Println(cfg)
	var acList = &acl.PolicyList{}
	for _, policy := range cfg {
		pol := policy.(map[string]interface{})
		subAttr := pol["subAttr"].(string)
		objAttr := pol["objAttr"].(string)
		optAttr := math.Floor(pol["optAttr"].(float64))
		fmt.Println(optAttr)
		envAttr := pol["envAttr"].(map[string]interface{})
		start := math.Floor(envAttr["start"].(float64))
		end := math.Floor(envAttr["end"].(float64))
		acList.Policies = append(acList.Policies, &acl.Policy{
			SubAttr:  subAttr,
			ObjAttr:  objAttr,
			OptAttrs: dataServerType.OpCode(optAttr),
			EnvAttr: commonType.Duration{
				Start: commonType.Timestamp(start),
				End:   commonType.Timestamp(end),
			},
		})
	}
	return acList, nil
}

func (dt *DataServer) GetAttrsForKey(key []byte) []string {
	key = append(key, '_')
	value, err := db.Get(dt.dataDB, key)
	if err != nil {
		fmt.Println(err.Error())
		//dt.logger.Error(err.Error())
		return []string{}
	}
	fmt.Println("value: ", string(value))

	attrs := []string{}
	err = json.Unmarshal(value, attrs)
	if err != nil {
		fmt.Println(err.Error())
		dt.logger.Error(err.Error())
		return []string{}
	}
	return attrs
}

func (dt *DataServer) ProcessOperation(op *dataServerType.Operation) []byte {
	switch op.OpCode {
	case dataServerType.OP_QRY:
		key := op.Key
		value, err := db.Get(dt.dataDB, key)
		if err != nil {
			dt.logger.Error(err.Error())
			return []byte{}
		}
		return value
	case dataServerType.OP_ADD:
		key := op.Key
		value := op.Value
		err := db.Put(dt.dataDB, key, value)
		if err != nil {
			dt.logger.Error(err.Error())
			return []byte{}
		}
		key = append(key, '_')
		err = db.Put(dt.dataDB, key, op.Ext)
		if err != nil {
			dt.logger.Error(err.Error())
			return []byte{}
		}
		return []byte("success")
	default:
		dt.logger.Error("no such opCode")
		return []byte{}
	}
}
