package acl

import (
	commonType "github.com/DhunterAO/goAuthChain/common/types"
	serverType "github.com/DhunterAO/goAuthChain/dataServer/server/types"
	"strings"
)

type Policy struct {
	SubAttr  string
	ObjAttr  string
	OptAttrs serverType.OpCode
	EnvAttr  commonType.Duration
}

func (policy *Policy) CheckOperation(gm bool, subAttrs []string, objAttrs []string, op serverType.OpCode, timestamp commonType.Timestamp) bool {
	if !policy.EnvAttr.ContainTime(timestamp) {
		return false
	}
	flag := false
	if op&policy.OptAttrs != 0 {
		flag = true
	}
	if !flag {
		return false
	}
	flag = false

	if gm {
		for _, subAttr := range subAttrs {
			if strings.HasPrefix(policy.SubAttr, subAttr) {
				flag = true
				break
			}
		}
	} else {
		for _, subAttr := range subAttrs {
			if policy.SubAttr == subAttr {
				flag = true
				break
			}
		}
	}
	if !flag {
		return false
	}
	flag = false
	for _, objAttr := range objAttrs {
		if objAttr == policy.ObjAttr {
			flag = true
			break
		}
	}
	if !flag {
		return false
	}
	return true
}

type PolicyList struct {
	Policies []*Policy
}

func (pl *PolicyList) CheckOperation(gm bool, subAttrs []string, objAttrs []string, op serverType.OpCode, timestamp commonType.Timestamp) bool {
	for _, policy := range pl.Policies {
		if policy.CheckOperation(gm, subAttrs, objAttrs, op, timestamp) {
			return true
		}
	}
	return false
}
