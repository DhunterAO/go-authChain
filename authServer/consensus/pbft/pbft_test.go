package pbft

import (
	"encoding/json"
	"fmt"
	"github.com/DhunterAO/goAuthChain/authServer/blockchain"
	"github.com/DhunterAO/goAuthChain/log"
	"os"
	"testing"
)

func initBlockchain(logger *log.Logger) *blockchain.Blockchain {
	bc, _ := blockchain.InitBlockchain("../../data/conf/blockchain.conf", logger)
	return bc
}

func TestInitPbftConfig(t *testing.T) {
	testLogFile, err := os.Create("../../data/logs/test.log")
	if err != nil {
		log.Fatalln("open file error !")
	}
	defer testLogFile.Close()
	logger := log.New(testLogFile, "[PBFT]", log.LstdFlags)
	logger.Info("start logging")
	bc := initBlockchain(logger)
	t.Log(string(bc.ToJson()))

	pt, err := InitPbftConfig("../../data/conf/consensus.conf", bc, logger)
	if err != nil {
		fmt.Println(err.Error())
		t.Fatal(err)
	}
	fmt.Println("pbft")
	j, _ := json.Marshal(pt)
	t.Log(string(j))
}

func TestStart(t *testing.T) {
	testLogFile, err := os.Create("../../data/logs/test.log")
	if err != nil {
		log.Fatalln("open file error !")
	}
	defer testLogFile.Close()
	logger := log.New(testLogFile, "[PBFT]", log.LstdFlags)
	logger.Info("start logging")

	bc := initBlockchain(logger)

	pt, err := InitPbftConfig("../../data/conf/consensus.conf", bc, logger)
	if err != nil {
		fmt.Println(err.Error())
		t.Fatal(err)
	}
	bs, _ := json.Marshal(pt)
	fmt.Println(string(bs))
	pt.Start()
}
