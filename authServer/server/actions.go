package server

import (
	"encoding/json"
	"fmt"
	types2 "github.com/DhunterAO/goAuthChain/authServer/blockchain/types"
	"github.com/DhunterAO/goAuthChain/authServer/consensus/pbft/pbftTypes"
	"github.com/DhunterAO/goAuthChain/common/authorization"
	"github.com/DhunterAO/goAuthChain/crypto/types"
	util2 "github.com/DhunterAO/goAuthChain/util"
	"github.com/DhunterAO/goAuthChain/util/netutil"
	"github.com/kataras/iris"
	"os"
	"time"
)

func getNonceFromAddress(ctx iris.Context) {
	address := ctx.Params().Get("address")
	nonce := authServer.Bc.State.GetAccountNonce(types.HexToAddress(address))
	_, err := ctx.Write([]byte(string(nonce)))
	if err != nil {
		fmt.Println(err)
	}
}

func getAccountStateFromAddress(ctx iris.Context) {
	address := ctx.Params().Get("address")
	state := authServer.Bc.State.GetAccountState(types.HexToAddress(address))

	_, err := ctx.Write(state.ToJson())
	if err != nil {
		fmt.Println(err)
	}
}

func getGmList(ctx iris.Context) {
	gmList := authServer.Bc.State.GmList
	gmListBytes, err := json.Marshal(gmList)
	if err != nil {
		fmt.Println(err)
	}
	_, err = ctx.Write(gmListBytes)
	if err != nil {
		fmt.Println(err)
	}
}

func getBlockchainIndex(ctx iris.Context) {
	//fmt.Println("get blockchain")
	_, err := ctx.Write(authServer.Bc.ToIndexJson())
	if err != nil {
		fmt.Println(err)
	}
}

func getBlockchain(ctx iris.Context) {
	fmt.Println("get blockchain")
	_, err := ctx.Write(authServer.Bc.ToPrettyJson())
	if err != nil {
		fmt.Println(err)
	}
}

func getAuthorizationPool(ctx iris.Context) {
	_, err := ctx.Write(authServer.Bc.AuthPool.ToJson())
	if err != nil {
		fmt.Println(err)
	}
}

func getGlobalState(ctx iris.Context) {
	fmt.Println("get global state")
	_, err := ctx.Write(util2.PrettyPrintJson(authServer.Bc.State.ToJson()))
	if err != nil {
		fmt.Println(err)
	}
}

func postAuthorization(ctx iris.Context) {
	fmt.Println("receive authorization")

	auth := new(authorization.Authorization)
	if err := ctx.ReadJSON(auth); err != nil {
		fmt.Println(err)
		_, _ = ctx.Writef("json format of authorization received")
		return
	}
	fmt.Println(auth.ToString())

	if !auth.CheckSignature() {
		fmt.Println("authServer.go line 108: authorization received with invalid signature")
		return
	}

	if flag, _ := authServer.Bc.AuthPool.AddAuthorization(auth); !flag {
		fmt.Println("invalid authorization")
		return
	}
}

func postBlock(ctx iris.Context) {
	block := new(types2.Block)
	err := ctx.ReadJSON(block)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func postViewChange(ctx iris.Context) {
	vc := new(pbftTypes.ViewChange)
	err := ctx.ReadJSON(vc)
	if err != nil {
		fmt.Println(err)
	}

	if !authServer.Pbft.CheckViewChange(vc) {
		_, _ = ctx.Writef("invalid view change")
		return
	}
	authServer.Pbft.ReceiveViewChange(vc)
}

func postRequest(ctx iris.Context) {
	// parse the post data into request
	req := new(pbftTypes.Request)
	err := ctx.ReadJSON(req)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("receive request: ", string(req.ToJson()))

	// the pbft processes request
	authServer.Pbft.ReceiveRequest(req)
}

func postResponse(ctx iris.Context) {
	resp := new(pbftTypes.Response)
	err := ctx.ReadJSON(resp)
	if err != nil {
		fmt.Println(err)
	}

	if !authServer.Pbft.CheckResponse(resp) {
		fmt.Println("invalid response")
	}

	authServer.Pbft.ReceiveResponse(resp)
}

func postCommit(ctx iris.Context) {
	commit := new(pbftTypes.Commit)
	err := ctx.ReadJSON(commit)
	if err != nil {
		fmt.Println(err)
	}

	if !authServer.Pbft.CheckCommit(commit) {
		_, _ = ctx.Writef("invalid commit")
	}

	authServer.Pbft.ReceiveCommit(commit)
}

func startAll(ctx iris.Context) {
	var ipList []string
	for _, candidate := range authServer.Pbft.CandidateInfo {
		ipList = append(ipList, candidate.Ip)
	}
	netutil.BroadCastGet(ipList, "/start")
}

func start(ctx iris.Context) {
	authServer.Pbft.Start()
}

func endAll(ctx iris.Context) {
	var ipList []string
	for _, candidate := range authServer.Pbft.CandidateInfo {
		ipList = append(ipList, candidate.Ip)
	}
	netutil.BroadCastGet(ipList, "/end")
}

func end(ctx iris.Context) {
	authServer.Pbft.StopAllTimer()
	time.Sleep(4 * time.Second)
	os.Exit(0)
}
