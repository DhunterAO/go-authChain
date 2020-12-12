package types

import (
	"fmt"
	"github.com/DhunterAO/goAuthChain/common/fileutil"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg := GenerateConfig()
	fmt.Println(cfg)
	err := fileutil.DumpJson("genesis.json", cfg)
	fmt.Println(err)
}
