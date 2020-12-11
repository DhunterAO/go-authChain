package server

import (
	"fmt"
	"github.com/kataras/iris"
	"goauth/common/operation"
	"goauth/crypto"
	"goauth/util/types"
)

func postOperation(ctx iris.Context) {
	op := new(operation.Operation)
	err := ctx.ReadJSON(op)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(op.ToJson())
	ctx.Write(processOperation(op))
}

func processOperation(op *operation.Operation) []byte {
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
	envAttr := types.CurrentTime()
	//fmt.Println("envAttrs: ", envAttr)
	if (op.OpCode == operation.OP_ADD && len(objAttrs) == 0) || dataServer.Acl.CheckOperation(gm, subAttrs, objAttrs, optAttr, envAttr) {
		return dataServer.ProcessOperation(op)
	}
	return []byte("no Permission")
}
