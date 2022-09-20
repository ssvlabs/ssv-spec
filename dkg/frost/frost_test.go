package frost

import (
	"fmt"
	"testing"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"golang.org/x/sync/errgroup"
)

type ProcessMsgFnType func(msg *dkg.SignedMessage) (bool, *dkg.KeyGenOutput, error)

type TestNetwork struct {
	ProcessMsgFn map[uint32]ProcessMsgFnType
	Outputs      map[uint32]*dkg.KeyGenOutput
}

func (tn *TestNetwork) StreamDKGOutput(output map[types.OperatorID]*dkg.SignedOutput) error {
	return nil
}

func (tn *TestNetwork) BroadcastDKGMessage(msg *dkg.SignedMessage) error {
	g := errgroup.Group{}
	for operatorID, f := range tn.ProcessMsgFn {
		fn := f
		operatorID := operatorID

		g.Go(func() error {
			finished, o, err := fn(msg)
			if finished {
				tn.Outputs[operatorID] = o
			}
			return err
		})
	}
	return g.Wait()
}

func TestFrost2_4(t *testing.T) {
	testNetwork := TestNetwork{
		ProcessMsgFn: make(map[uint32]ProcessMsgFnType),
		Outputs:      make(map[uint32]*dkg.KeyGenOutput),
	}

	requestID := dkg.RequestID{}
	for i, _ := range requestID {
		requestID[i] = 1
	}

	operators := []types.OperatorID{
		1, 2, 3, 4,
	}

	kgps := make(map[uint32]dkg.KeyGenProtocol)

	for _, operatorID := range operators {
		p := New(requestID, &testNetwork, uint32(operatorID))
		kgps[uint32(operatorID)] = p

		testNetwork.ProcessMsgFn[uint32(operatorID)] = p.ProcessMsg
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
			err := kgps[uint32(operatorID)].Start(initMsg)
			if err != nil {
				return fmt.Errorf("id %d err %s", operatorID, err.Error())
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}

	for _, operatorID := range operators {
		output := testNetwork.Outputs[uint32(operatorID)]
		t.Logf("printing generated keys for id %d\n", operatorID)
		t.Logf("sk %x", output.Share.Serialize())
		t.Logf("vk %x", output.ValidatorPK)
		for opID, publicKey := range output.OperatorPubKeys {
			t.Logf("id %d pk %x", opID, publicKey.Serialize())
		}
	}
}
