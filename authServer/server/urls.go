package server

import (
	"github.com/DhunterAO/goAuthChain/authServer/consensus/pbft/pbftTypes"
	"github.com/kataras/iris"
)

func addUrls(app *iris.Application) {
	app.Get("/", getBlockchainIndex)
	app.Get("/gmlist", getGmList)
	app.Get("/authPool", getAuthorizationPool)
	app.Get("/state", getGlobalState)
	app.Get("/blockchain", getBlockchain)

	// get personal info
	app.Get("/api/nonce/{address:string}", getNonceFromAddress)
	app.Get("/api/state/{address:string}", getAccountStateFromAddress)

	// post info
	app.Post("/api/authorization", postAuthorization)
	//app.Post("/api/operation", postOperation)
	app.Post("/api/block", postBlock)

	// pbft info
	app.Post(pbftTypes.RequestUrl, postRequest)
	app.Post(pbftTypes.ResponseUrl, postResponse)
	app.Post(pbftTypes.CommitUrl, postCommit)
	app.Post(pbftTypes.ViewChangeUrl, postViewChange)

	// server control
	app.Get("/start", start)
	app.Get("/startall", startAll)
	app.Get("/end", end)
	app.Get("/endall", endAll)
}
