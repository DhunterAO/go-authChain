package server

import (
	"fmt"
	"github.com/DhunterAO/goAuthChain/authServer/blockchain"
	"github.com/DhunterAO/goAuthChain/log"
	"github.com/kataras/iris"
	"os"
	"testing"
)

func TestInitAuthServer(t *testing.T) {
	logPath := "../../data/logs/"
	confPath := "../../data/conf/"

	blockchainLogFile, err := os.Create(logPath + "blockchain.log")
	if err != nil {
		log.Fatalln("open blockchain log file error !")
	}
	defer blockchainLogFile.Close()
	blockchainLogger := log.New(blockchainLogFile, "[LOG]", log.LstdFlags)

	authServerLogFile, err := os.Create(logPath + "authServer.log")
	if err != nil {
		log.Fatalln("open auth server log file error !")
	}
	defer authServerLogFile.Close()
	authServerLogger := log.New(authServerLogFile, "[LOG]", log.LstdFlags)

	bc, _ := blockchain.InitBlockchain(confPath+"blockchain.conf", blockchainLogger)
	fmt.Println(string(bc.ToJson()))
	au, _ := InitAuthServer(confPath+"consensus.conf", bc, authServerLogger)
	fmt.Println("init authServer finish")
	authServer.Pbft.Start()
	err = au.App.Run(iris.Addr(":" + au.Pbft.Port))
	if err != nil {
		fmt.Println("app start err")
	}
}
