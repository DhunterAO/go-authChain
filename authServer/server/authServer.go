package server

import (
	"fmt"
	"github.com/DhunterAO/goAuthChain/authServer/blockchain"
	"github.com/DhunterAO/goAuthChain/authServer/consensus/pbft"
	"github.com/DhunterAO/goAuthChain/log"
	"github.com/kataras/iris"
)

type AuthServer struct {
	Bc   *blockchain.Blockchain
	Pbft *pbft.Pbft
	App  *iris.Application
	log  *log.Logger
}

var authServer AuthServer

func InitAuthServer(cfgPath string, bc *blockchain.Blockchain, logger *log.Logger) (*AuthServer, error) {
	fmt.Println("start init server")
	// initialize pbft config
	pt, err := pbft.InitPbftConfig(cfgPath, bc, logger)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	authServer.Pbft = pt
	authServer.App = iris.New()
	authServer.Bc = bc
	// add urls into application of authServer
	addUrls(authServer.App)
	return &authServer, nil
}

func (au *AuthServer) Start() {
	au.Pbft.Start()
	err := au.App.Run(iris.Addr(":" + au.Pbft.Port))
	if err != nil {
		fmt.Println("app start err")
	}
}
