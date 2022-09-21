package frost

import (
	"context"
	"fmt"
	"sync"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
)

type ProcessMsgFnType func(msg *dkg.SignedMessage) (bool, *dkg.KeyGenOutput, error)

type Node struct {
	ID           types.OperatorID
	ProcessMsgFn ProcessMsgFnType

	mu     *sync.Mutex
	Output *dkg.KeyGenOutput

	queue  chan *dkg.SignedMessage
	ctx    context.Context
	cancel context.CancelFunc
}

func newNode(id types.OperatorID) *Node {
	ctx, cancel := context.WithCancel(context.Background())
	return &Node{
		ID: id,
		mu: &sync.Mutex{},

		queue:  make(chan *dkg.SignedMessage, 100),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (n *Node) SetProcessMsgFn(fn ProcessMsgFnType) {
	n.ProcessMsgFn = fn
}

func (n *Node) Add(msg *dkg.SignedMessage) {
	n.queue <- msg
}

func (n *Node) Run() error {
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

func (n *Node) processMsg(msg *dkg.SignedMessage) error {
	finished, o, err := n.ProcessMsgFn(msg)
	if finished {
		n.mu.Lock()
		n.Output = o
		n.mu.Unlock()
	}
	return err
}

type MockNetwork struct {
	nodes map[types.OperatorID]*Node
}

func (ntwrk *MockNetwork) StreamDKGOutput(output map[types.OperatorID]*dkg.SignedOutput) error {
	return nil
}

func (ntwrk *MockNetwork) BroadcastDKGMessage(msg *dkg.SignedMessage) error {
	wg := &sync.WaitGroup{}
	for _, node := range ntwrk.nodes {
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
