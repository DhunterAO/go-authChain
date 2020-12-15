package server

import (
	"fmt"
	commonType "github.com/DhunterAO/goAuthChain/common/types"
	"github.com/DhunterAO/goAuthChain/crypto"
	types2 "github.com/DhunterAO/goAuthChain/dataServer/server/types"
	"github.com/kataras/iris"
)

func postOperation(ctx iris.Context) {
	op := new(types2.Operation)
	err := ctx.ReadJSON(op)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(op.ToJson())
	ctx.Write(processOperation(op))
}

func processOperation(op *types2.Operation) []byte {
	//fmt.Println(op.OpCode)
	pk, err := crypto.RecoverPkFromSig(op.CalcHash(), op.Signature)
	if err != nil {
		fmt.Println("signature error")
		return []byte("signature error")
	}
	addr, err := crypto.PubkeyToAddress(pk)
	if err != nil {
		return []byte("signature error")
	}
	gm, subAttrs := dataServer.Bc.State.GetAttrsFromAddr(addr)
	//fmt.Println("subAttrs: ", subAttrs, " gm: ", gm)
	objAttrs := dataServer.GetAttrsForKey(op.Key)
	//fmt.Println("objAttrs: ", objAttrs)
	optAttr := op.OpCode
	//fmt.Println("optAttr: ", optAttr)
	envAttr := commonType.CurrentTime()
	//fmt.Println("envAttrs: ", envAttr)
	if (op.OpCode == types2.OP_ADD && len(objAttrs) == 0) || dataServer.Acl.CheckOperation(gm, subAttrs, objAttrs, optAttr, envAttr) {
		return dataServer.ProcessOperation(op)
	}
	return []byte("no Permission")
}
