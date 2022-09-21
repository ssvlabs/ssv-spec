package frost

import (
	"testing"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func TestFrost2_4(t *testing.T) {
	requestID := dkg.RequestID{}
	for i, _ := range requestID {
		requestID[i] = 1
	}

	operators := []types.OperatorID{
		1, 2, 3, 4,
	}

	mockNetwork := MockNetwork{
		nodes: make(map[types.OperatorID]*Node),
	}

	for _, operator := range operators {
		operatorID := types.OperatorID(operator)

		node := newNode(operatorID)
		mockNetwork.nodes[operatorID] = node
	}

	kgps := make(map[uint32]dkg.KeyGenProtocol)

	for _, operatorID := range operators {
		p := New(requestID, &mockNetwork, uint32(operatorID))
		kgps[uint32(operatorID)] = p

		mockNetwork.nodes[operatorID].SetProcessMsgFn(p.ProcessMsg)
	}

	for _, node := range mockNetwork.nodes {
		go node.Run()
		defer node.cancel()
	}

	threshold := 2

	g := errgroup.Group{}
	for _, operatorID := range operators {
		operatorID := operatorID

		initMsg := &dkg.Init{
			OperatorIDs: operators,
			Threshold:   uint16(threshold),
		}

		g.Go(func() error {
			if err := kgps[uint32(operatorID)].Start(initMsg); err != nil {
				return errors.Wrapf(err, "failed to start operator %d", operatorID)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}

	for {
		finished := true
		for _, node := range mockNetwork.nodes {
			node.mu.Lock()
			if node.Output == nil {
				finished = false
			}
			node.mu.Unlock()
		}

		if finished {
			break
		}
	}

	for _, operatorID := range operators {
		output := mockNetwork.nodes[operatorID].Output
		t.Logf("printing generated keys for id %d\n", operatorID)
		t.Logf("sk %x", output.Share.Serialize())
		t.Logf("vk %x", output.ValidatorPK)
		for opID, publicKey := range output.OperatorPubKeys {
			t.Logf("id %d pk %x", opID, publicKey.Serialize())
		}
	}
}
