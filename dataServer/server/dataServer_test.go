package server

//
//import (
//	"fmt"
//	"github.com/DhunterAO/goAuthChain/authServer/blockchain"
//	"github.com/DhunterAO/goAuthChain/common/operation"
//	"github.com/DhunterAO/goAuthChain/crypto"
//	"github.com/DhunterAO/goAuthChain/crypto/types"
//	db2 "github.com/DhunterAO/goAuthChain/dataServer/db"
//	"github.com/DhunterAO/goAuthChain/log"
//	"os"
//	"testing"
//)
//
//func TestInitDataServer(t *testing.T) {
//	logPath := "../../data/logs/"
//	confPath := "../../data/conf/"
//
//	blockchainLogFile, err  := os.Create(logPath + "blockchain.log")
//	if err != nil {
//		log.Fatalln("open blockchain log file error !")
//	}
//	defer blockchainLogFile.Close()
//	blockchainLogger := log.New(blockchainLogFile, "[LOG]", log.LstdFlags)
//
//	dataServerLogFile, err  := os.Create(logPath + "authServer.log")
//	if err != nil {
//		log.Fatalln("open auth server log file error !")
//	}
//	defer dataServerLogFile.Close()
//	dataServerLogger := log.New(dataServerLogFile, "[LOG]", log.LstdFlags)
//
//	bc, _ := blockchain.InitBlockchain(confPath + "blockchain.conf", blockchainLogger)
//	fmt.Println(string(bc.ToJson()))
//	dataS, _ := InitDataServer(confPath + "dataServer.conf", bc, dataServerLogger)
//	fmt.Println("init authServer finish")
//	dataS.Start()
//}
//
//func TestLoadAclList(t *testing.T) {
//	LoadAcList("../../data/acl/abac.list")
//}
//
//func TestDataServer_ProcessOperation(t *testing.T) {
//	sk, _ := types.LoadECDSA("../../data/keys/test.key")
//	addr, _ := crypto.PubkeyToAddress(types.BytesToPubkey(types.CompressPubkey(&sk.PublicKey)))
//	fmt.Println(addr.Hex())
//	op := operation.NewOperation(addr, operation.OP_ADD, []byte("exam_paper_2019_"), []byte("[\"A_exam\",\"A_grade\"]"), nil)
//	sig, _ := crypto.Sign(op.CalcHash(), sk)
//	op.Sign(sig)
//	fmt.Println("operation", op.ToString())
//	flag, _ := op.CheckSignature()
//	fmt.Println("check signature of operation: ", flag)
//
//	logger := log.NewTestLogger()
//	bc, _ := blockchain.InitBlockchain("../../data/conf/blockchain.conf", logger)
//	_, _ = InitDataServer("../../data/conf/dataServer.conf", bc, logger)
//	fmt.Println("init data server")
//	fmt.Println(dataServer.dataDB)
//	fmt.Println(string(processOperation(op)))
//}
//
//func TestDataServer_ProcessOperation2(t *testing.T) {
//	db, _ := db2.OpenDB("../../data/dataDB/")
//	_ = db2.Del(db, []byte("exam_paper_2019"))
//	val, _ := db2.Get(db, []byte("exam_paper_2019_"))
//	fmt.Println(string(val))
//}
