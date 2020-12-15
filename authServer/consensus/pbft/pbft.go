package pbft

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/DhunterAO/goAuthChain/authServer/consensus/pbft/pbftTypes"
	"github.com/DhunterAO/goAuthChain/blockchain"
	"github.com/DhunterAO/goAuthChain/blockchain/types"
	"github.com/DhunterAO/goAuthChain/common/fileutil"
	"github.com/DhunterAO/goAuthChain/common/netutil"
	commonType "github.com/DhunterAO/goAuthChain/common/types"
	"github.com/DhunterAO/goAuthChain/crypto"
	cryptoType "github.com/DhunterAO/goAuthChain/crypto/types"
	"github.com/DhunterAO/goAuthChain/log"
	"path"
	"runtime"
	"sync"
)

const (
	BlockInterval      = 5
	NullRequestTime    = 2
	ViewChangeInterval = 60
)

var (
	CandidateLack = errors.New("the parameter candidate must be needed")
	NotCandidate  = errors.New("the address of server is not a candidate")
)

type NodeInfo struct {
	Ip   string
	Addr cryptoType.Address
}

type Pbft struct {
	// config of current code
	NodeID        uint64            `json:"node_id"`
	CandidateInfo []*NodeInfo       `json:"candidate_info"`
	Sk            *ecdsa.PrivateKey `json:"-"`
	Host          string            `json:"host"`
	Port          string            `json:"port"`

	// pbft consensus status
	N          uint64 `json:"n"`           // the number of total nodes
	F          uint64 `json:"f"`           // the limit number of malicious nodes
	ViewID     uint64 `json:"view_id"`     // view id
	NewViewID  uint64 `json:"new_view_id"` // next view id
	ResendTime uint64 `json:"resend_time"` // resend time

	// pools used to store received messages in pbft
	ViewChangePool map[uint64]*commonType.Uint64Set          `json:"-"`
	ResponsePool   map[cryptoType.Hash]*commonType.Uint64Set `json:"-"`
	CommitPool     map[cryptoType.Hash]*commonType.Uint64Set `json:"-"`
	BlockPool      map[cryptoType.Hash]*types.Block          `json:"-"`

	// timers used for control process of pbft
	NullRequestTimer      *commonType.MyTimer `json:"-"` // change view when null request received in an interval
	NewBlockTimer         *commonType.MyTimer `json:"-"` // create/wait for new block every 1 min
	ViewChangeTimer       *commonType.MyTimer `json:"-"` // view change timer used to change the primary each minute
	ViewChangeResendTimer *commonType.MyTimer `json:"-"` // resend view change
	BlockPassTimer        *commonType.MyTimer `json:"-"` // change view when time out for a new block

	// blockchain used for check messages in pbft consensus
	bc *blockchain.Blockchain `json:"-"`

	logger *log.Logger `json:"-"`
	ViewMu sync.Mutex  `json:"-"` // mutex for view id
}

var pbft Pbft

func InitPbftConfig(cfgPath string, bc *blockchain.Blockchain, logger *log.Logger) (*Pbft, error) {
	var cfg map[string]interface{}
	err := fileutil.LoadJson(cfgPath, &cfg)
	if err != nil {
		logger.Error("load pbft configuration failed !")
		return nil, err
	}

	// check parameters in cfg including nodeID, candidateInfo, keyPath, host, port
	candidateList, ok := cfg["candidateInfo"].([]interface{})
	if !ok {
		return nil, CandidateLack
	}

	//keyPath, ok := cfg["keyPath"].(string)
	//if !ok {
	//	keyPath = "data/keys/test.key"
	//}
	const dataFile = "../../data/keys/test.key"
	_, filename, _, _ := runtime.Caller(1)
	keyPath := path.Join(path.Dir(filename), dataFile)
	//fmt.Println(keyPath)
	//fmt.Println(keyPath)

	host, ok := cfg["host"].(string)
	if !ok {
		host = "127.0.0.1"
	}

	port, ok := cfg["port"].(string)
	if !ok {
		port = "8000"
	}

	// process parameters loaded
	sk, err := cryptoType.LoadECDSA(keyPath)
	if err != nil {
		return nil, err
	}
	addr := cryptoType.EcdsaPubkeyToAddress(sk.PublicKey)

	nodeID := -1
	var candidateInfos []*NodeInfo
	for idx, candidate := range candidateList {
		info := candidate.(map[string]interface{})
		nodeInfo := &NodeInfo{
			Ip:   info["Ip"].(string),
			Addr: *cryptoType.HexToAddress(info["Addr"].(string)),
		}
		candidateInfos = append(candidateInfos, nodeInfo)
		if nodeInfo.Addr == *addr {
			nodeID = idx
		}
	}
	n := len(candidateInfos)
	f := (n - 1) / 3

	if nodeID == -1 {
		return nil, NotCandidate
	}

	// initialize pbft with configuration
	pbft = Pbft{
		NodeID:                uint64(nodeID),
		CandidateInfo:         candidateInfos,
		Sk:                    sk,
		Host:                  host,
		Port:                  port,
		N:                     uint64(n),
		F:                     uint64(f),
		ViewID:                0,
		NewViewID:             0,
		ResendTime:            1,
		ViewChangePool:        make(map[uint64]*commonType.Uint64Set),
		ResponsePool:          make(map[cryptoType.Hash]*commonType.Uint64Set),
		CommitPool:            make(map[cryptoType.Hash]*commonType.Uint64Set),
		BlockPool:             make(map[cryptoType.Hash]*types.Block),
		NullRequestTimer:      commonType.NewMyTimer("null request timer ", NullRequestTime, ChangeView, logger),
		NewBlockTimer:         commonType.NewMyTimer("new block timer ", BlockInterval, NewBlock, logger),
		ViewChangeTimer:       commonType.NewMyTimer("view change timer ", ViewChangeInterval, ChangeView, logger),
		ViewChangeResendTimer: commonType.NewMyTimer("view change resend timer ", pbft.ResendTime, ChangeView, logger),
		BlockPassTimer:        commonType.NewMyTimer("block time out timer ", BlockInterval/2, ChangeView, logger),
		bc:                    bc,
		logger:                logger,
		ViewMu:                sync.Mutex{},
	}
	return &pbft, nil
}

func (pt *Pbft) Start() {
	fmt.Println("start pbft")
	pt.NewBlockTimer.Start(pt.NewBlockTimer.Delay)
	//pt.ViewChangeTimer.Start(pt.ViewChangeTimer.Delay)
}

func BroadcastPostInfo(message []byte, infos []*NodeInfo, url string, self bool) {
	var ipList []string
	for idx, info := range infos {
		if uint64(idx) == pbft.NodeID {
			if self {
				ipList = append(ipList, info.Ip)
			}
		} else {
			ipList = append(ipList, info.Ip)
		}
	}
	//fmt.Println("post info: ",ipList, url, message)
	netutil.BroadCastPost(ipList, url, message)
}

func (pt *Pbft) ReceiveRequest(req *pbftTypes.Request) {
	// check the request received first
	if !pt.CheckRequest(req) {
		fmt.Println("request check error")
		return
	}
	// stop null request timer
	pt.NullRequestTimer.Stop()
	// add block in request into pool
	block := types.JsonToBlock(req.BlockJson)
	hash := *block.CalcHash()

	if _, ok := pt.BlockPool[hash]; ok {
		fmt.Println("block in received request existed")
		return
	}

	pt.BlockPool[hash] = block

	// broadcast response if accept the request
	response := pbftTypes.NewResponse(pbft.NodeID, pbft.ViewID, hash)
	response.Signature, _ = crypto.Sign(response.CalcHash(), pt.Sk)
	//check, _ := response.CheckSignature()
	//fmt.Println("check response generated: ", check)
	BroadcastPostInfo(response.ToJson(), pt.CandidateInfo, pbftTypes.ResponseUrl, true)
}

// check the request received
func (pt *Pbft) CheckRequest(req *pbftTypes.Request) bool {
	// check the info ID of req
	if req.InfoID != pbftTypes.RequestInfo {
		fmt.Println("info id in request is invalid")
		return false
	}

	// check the view ID of req
	if req.ViewID != pbft.ViewID {
		fmt.Println("view id in request is invalid")
		return false
	}

	// check the node ID of req
	if req.NodeID != pbft.ViewID%pbft.N {
		fmt.Println("node id in request is not current primary")
		return false
	}

	// check the signature of req
	if check, _ := req.CheckSignature(); !check {
		fmt.Println("signature in request is invalid")
		return false
	}

	// check the block in request
	b := types.JsonToBlock(req.BlockJson)
	if !pt.bc.CheckNewBlock(b) {
		pt.logger.Warn("invalid new block")
		return false
	}
	return true
}

// process response received
func (pt *Pbft) ReceiveResponse(resp *pbftTypes.Response) {
	// check the response received first
	if !pt.CheckResponse(resp) {
		return
	}
	// add the nodeID into the response pool according to the block hash
	nodeID := resp.NodeID
	hash := resp.BlockHash
	if _, ok := pt.ResponsePool[hash]; !ok {
		pt.ResponsePool[hash] = commonType.NewUint64Set()
	}
	pt.ResponsePool[hash].Add(nodeID)

	// if the number of set equal to the threshold, the node becomes prepared state and broadcasts commit
	if pt.ResponsePool[hash].Len() >= pbft.N-pbft.F {
		pt.ResponsePool[hash].Clear()
		// generate
		commit := pbftTypes.NewCommit(pbft.NodeID, pbft.ViewID, &hash)
		commit.Signature, _ = crypto.Sign(commit.CalcHash(), pt.Sk)

		BroadcastPostInfo(commit.ToJson(), pt.CandidateInfo, pbftTypes.CommitUrl, true)
	}
}

// check response response received
func (pt *Pbft) CheckResponse(resp *pbftTypes.Response) bool {
	// check info ID of response
	if resp.InfoID != pbftTypes.ResponseInfo {
		fmt.Println("info id in response is invalid")
		return false
	}

	// check view ID of response
	if resp.ViewID != pt.ViewID {
		fmt.Println("view id in response is invalid")
		return false
	}

	// check signature of response
	if flag, _ := resp.CheckSignature(); !flag {
		fmt.Println("signature in response is invalid")
		return false
	}
	return true
}

func (pt *Pbft) ReceiveCommit(commit *pbftTypes.Commit) {
	if !pt.CheckCommit(commit) {
		return
	}

	nodeID := commit.NodeID
	hash := *commit.BlockHash

	if _, ok := pt.CommitPool[hash]; !ok {
		pt.CommitPool[hash] = commonType.NewUint64Set()
	}
	pt.CommitPool[hash].Add(nodeID)

	// if the number of received commits achieve the threshold, then add the block into blockchain
	if pt.CommitPool[hash].Len() >= pt.N-pt.F {
		pt.CommitPool[hash].Clear()
		if block, ok := pt.BlockPool[hash]; ok {
			if pt.bc.CheckNewBlock(block) {
				pt.logger.Info("chain add new block" + string(pt.bc.GetHeight()))
				if _, err := pt.bc.AddNewBlock(block); err != nil {
					pt.logger.Error(err.Error())
				}
				//fmt.Println(string(pt.bc.ToJson()))
				pt.BlockPassTimer.Stop()
			}
		}
	}
}

func (pt *Pbft) CheckCommit(commit *pbftTypes.Commit) bool {
	// check info id of commit
	if commit.InfoID != pbftTypes.CommitInfo {
		fmt.Println("info id in commit is invalid")
		return false
	}

	// the view id of commit should not be checked

	// check node id of commit with the public key from signature
	nodeID := commit.NodeID
	pk, err := cryptoType.RecoverPkFromSig(commit.CalcHash(), commit.Signature)
	if err != nil {
		pt.logger.Warn(err.Error())
	}
	addr, err := crypto.PubkeyToAddress(pk)
	if err != nil {
		pt.logger.Warn(err.Error())
	}
	if addr == nil || *addr != pt.CandidateInfo[nodeID].Addr {
		return false
	}

	// check signature of commit
	if !commit.CheckSignature() {
		fmt.Println("signature in commit is invalid")
		return false
	}
	return true
}

func (pt *Pbft) CheckViewChange(vc *pbftTypes.ViewChange) bool {
	// check info ID of commit
	if vc.InfoID != pbftTypes.ViewChangeInfo {
		pt.logger.Warn("info id in view change is invalid")
		return false
	}

	// check view ID of commit
	if vc.NewViewID <= pbft.ViewID {
		pt.logger.Warn("new view id is too old")
		return false
	}
	return true
}

func (pt *Pbft) ReceiveViewChange(vc *pbftTypes.ViewChange) {
	//fmt.Println("receive view change from ", vc.NodeID)
	nodeID := vc.NodeID
	newViewID := vc.NewViewID
	if _, ok := pt.ViewChangePool[newViewID]; !ok {
		pt.ViewChangePool[newViewID] = commonType.NewUint64Set()
	}
	pt.ViewChangePool[newViewID].Add(nodeID)

	// if the number of received view change achieves the threshold, then convert into the new view ID
	if pt.ViewChangePool[newViewID].Len() >= pt.N-pt.F {
		pt.ViewMu.Lock()
		defer pt.ViewMu.Unlock()
		//fmt.Println(pt.ViewChangePool[newViewID])

		pt.ViewChangePool[newViewID].Clear()

		//fmt.Println("change viewID from ", ViewID, " to ", NewViewID)
		pt.ViewID = newViewID
		pt.ViewChangeTimer.Reset(ViewChangeInterval)
		pt.ViewChangeResendTimer.Stop()
		pt.NewBlockTimer.Reset(BlockInterval)
	}
}

func NewBlock() {
	pbft.BlockPassTimer.Reset(pbft.BlockPassTimer.Delay)
	pbft.NewBlockTimer.Reset(pbft.NewBlockTimer.Delay)
	//fmt.Println("new block")
	// check if this node is the primary now
	if pbft.ViewID%pbft.N == pbft.NodeID {
		// if primary, generate a new block from bc, and broadcast request to backups
		//fmt.Println("i'm the primary")
		newBlock := pbft.bc.GenerateNewBlock()
		//fmt.Println(newBlock.ToString())
		pbft.logger.Info(newBlock.ToString())
		request := pbftTypes.NewRequest(pbft.NodeID, pbft.ViewID, newBlock.ToJson())
		sig, err := crypto.Sign(request.CalcHash(), pbft.Sk)
		if err != nil {
			pbft.logger.Error(err.Error())
		}
		request.Sign(sig)
		//fmt.Println("generate signature: ", sig, err)
		//check, err := request.CheckSignature()
		//fmt.Println("check signature: ", check, err)
		BroadcastPostInfo(request.ToJson(), pbft.CandidateInfo, pbftTypes.RequestUrl, true)
	} else {
		// if backup, start null request timer, waiting for request from the primary
		//fmt.Println("i'm a backup")
		pbft.NullRequestTimer.Reset(NullRequestTime)
	}
}

// some error occurs from current primary, try to change view to next candidate
func ChangeView() {
	pbft.logger.Info("start change view")
	// stop all timer during change view
	pbft.StopAllTimer()

	pbft.ViewMu.Lock()
	defer pbft.ViewMu.Unlock()

	// convert newViewID into next one and update resentTime
	pbft.NewViewID += 1
	pbft.ResendTime *= 2
	if pbft.ResendTime > 60 {
		pbft.ResendTime = 60
	}
	// generate viewChange and broadcast it to all nodes
	vc := pbftTypes.NewViewChange(pbft.NodeID, pbft.ViewID, pbft.NewViewID)
	BroadcastPostInfo(vc.Bytes(), pbft.CandidateInfo, pbftTypes.ViewChangeUrl, true)
	// reset the resend timer, if the view change failed, start a new one
	pbft.ViewChangeResendTimer.Reset(pbft.ResendTime)
}

func (pt *Pbft) StopAllTimer() {
	pt.ViewChangeTimer.Stop()
	pt.NullRequestTimer.Stop()
	pt.NewBlockTimer.Stop()
	pt.ViewChangeResendTimer.Stop()
	pt.BlockPassTimer.Stop()
}
