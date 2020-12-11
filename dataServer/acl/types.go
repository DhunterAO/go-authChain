package acl

import (
	"goauth/common/authorization"
	"goauth/common/operation"
	"goauth/util/types"
	"strings"
)

type Policy struct {
	SubAttr  string
	ObjAttr  string
	OptAttrs operation.OpCode
	EnvAttr  authorization.Duration
}

func (policy *Policy) CheckOperation(gm bool, subAttrs []string, objAttrs []string, op operation.OpCode, timestamp types.Timestamp) bool {
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

func (pl *PolicyList) CheckOperation(gm bool, subAttrs []string, objAttrs []string, op operation.OpCode, timestamp types.Timestamp) bool {
	for _, policy := range pl.Policies {
		if policy.CheckOperation(gm, subAttrs, objAttrs, op, timestamp) {
			return true
		}
	}
	return false
}
