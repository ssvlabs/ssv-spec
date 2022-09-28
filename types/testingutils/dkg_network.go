package testingutils

import (
	"context"
	"fmt"
	"sync"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
)

type DKGNetwork struct {
	Nodes map[types.OperatorID]*DKGNetworkNode
}

func NewDKGNetwork(nodes map[types.OperatorID]*DKGNetworkNode) dkg.Network {
	return &DKGNetwork{
		Nodes: nodes,
	}
}

func (network *DKGNetwork) GetNetworkNode(operatorID types.OperatorID) *DKGNetworkNode {
	return network.Nodes[operatorID]
}

func (network *DKGNetwork) StreamDKGOutput(output map[types.OperatorID]*dkg.SignedOutput) error {
	return nil
}

func (network *DKGNetwork) BroadcastDKGMessage(msg *dkg.SignedMessage) error {
	wg := &sync.WaitGroup{}
	for _, node := range network.Nodes {
		if node.ID == msg.Signer {
			continue
		}

		node := node
		wg.Add(1)
		go func() {
			node.Add(msg)
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

type ProcessMsgFnType func(msg *dkg.SignedMessage) (bool, *dkg.KeyGenOutput, error)

type DKGNetworkNode struct {
	ID           types.OperatorID
	ProcessMsgFn ProcessMsgFnType

	Mu     *sync.Mutex
	Output *dkg.KeyGenOutput

	queue  chan *dkg.SignedMessage
	ctx    context.Context
	cancel context.CancelFunc
}

func NewDKGNetworkNode(id types.OperatorID) *DKGNetworkNode {
	ctx, cancel := context.WithCancel(context.Background())
	return &DKGNetworkNode{
		ID: id,
		Mu: &sync.Mutex{},

		queue:  make(chan *dkg.SignedMessage, 100),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (n *DKGNetworkNode) SetProcessMsgFn(fn ProcessMsgFnType) {
	n.ProcessMsgFn = fn
}

func (n *DKGNetworkNode) Add(msg *dkg.SignedMessage) {
	n.queue <- msg
}

func (n *DKGNetworkNode) Run() error {
	for {
		select {
		case <-n.ctx.Done():
			return nil
		case msg := <-n.queue:
			err := n.processMsg(msg)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}
}

func (n *DKGNetworkNode) Exit() {
	n.cancel()
	close(n.queue)
}

func (n *DKGNetworkNode) processMsg(msg *dkg.SignedMessage) error {
	finished, o, err := n.ProcessMsgFn(msg)
	if finished {
		n.Mu.Lock()
		n.Output = o
		n.Mu.Unlock()
	}
	return err
}
