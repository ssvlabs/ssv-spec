package frost

import (
	crand "crypto/rand"
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/coinbase/kryptology/pkg/sharing"
	ecies "github.com/ecies/go/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestFrost2_4(t *testing.T) {
	requestID := dkg.RequestID{}
	for i := range requestID {
		requestID[i] = 1
	}

	operators := []types.OperatorID{
		1, 2, 3, 4,
	}

	nodes := make(map[types.OperatorID]*testingutils.DKGNetworkNode)

	dkgsigner := testingutils.NewTestingKeyManager()
	storage := testingutils.NewTestingStorage()

	for _, operator := range operators {
		operatorID := types.OperatorID(operator)

		node := testingutils.NewDKGNetworkNode(operatorID)
		nodes[operatorID] = node
	}

	network := testingutils.NewDKGNetwork(nodes)
	dkgNetwork := network.(*testingutils.DKGNetwork)

	kgps := make(map[uint32]dkg.KeyGenProtocol)

	for _, operatorID := range operators {
		p := New(network, operatorID, requestID, dkgsigner, storage)
		kgps[uint32(operatorID)] = p

		dkgNetwork.GetNetworkNode(operatorID).SetProcessMsgFn(p.ProcessMsg)
	}

	for _, node := range dkgNetwork.Nodes {
		go node.Run()
		defer node.Exit()
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
		for _, node := range dkgNetwork.Nodes {
			node.Mu.Lock()
			if node.Output == nil {
				finished = false
			}
			node.Mu.Unlock()
		}

		if finished {
			break
		}
	}

	for _, operatorID := range operators {
		output := dkgNetwork.Nodes[operatorID].Output
		t.Logf("printing generated keys for id %d\n", operatorID)
		t.Logf("sk %x", output.Share.Serialize())
		t.Logf("vk %x", output.ValidatorPK)
		for opID, publicKey := range output.OperatorPubKeys {
			t.Logf("id %d pk %x", opID, publicKey.Serialize())
		}
	}
}

func getRandRequestID() dkg.RequestID {
	requestID := dkg.RequestID{}
	for i := range requestID {
		requestID[i] = byte(rand.Intn(256))
	}
	return requestID
}

func getSignedMessage(requestID dkg.RequestID, operatorID types.OperatorID) *dkg.SignedMessage {
	storage := testingutils.NewTestingStorage()
	signer := testingutils.NewTestingKeyManager()

	signedMessage := &dkg.SignedMessage{
		Message: &dkg.Message{
			MsgType:    dkg.ProtocolMsgType,
			Identifier: requestID,
			Data:       []byte{1, 1, 1, 1, 1},
		},
		Signer:    operatorID,
		Signature: nil,
	}

	_, op, _ := storage.GetDKGOperator(operatorID)
	sig, _ := signer.SignDKGOutput(signedMessage, op.ETHAddress)
	signedMessage.Signature = sig
	return signedMessage
}

func TestProcessBlameTypeInconsistentMessage(t *testing.T) {
	reqID := getRandRequestID()

	data := getSignedMessage(reqID, 1)
	dataBytes, _ := data.Encode()

	validData := getSignedMessage(reqID, 1)
	validDataBytes, _ := validData.Encode()

	// tamperedData := getSignedMessage(reqID, 1)
	// tamperedData.Message.Data = []byte{2, 2, 2, 2, 2}
	// tamperedDataBytes, _ := tamperedData.Encode()

	tests := map[string]struct {
		blameMessage *BlameMessage
		expected     bool
	}{
		"blame_req_is_valid": {
			blameMessage: &BlameMessage{
				Type:      InconsistentMessage,
				BlameData: [][]byte{dataBytes, validDataBytes},
			},
			expected: true,
		},
		/*
			TODO: Uncomment this section once signed message's validate
			function is implemented
		*/
		// "blame_req_is_invalid": {
		// 	blameMessage: &BlameMessage{
		// 		Type:      InconsistentMessage,
		// 		BlameData: [][]byte{dataBytes, tamperedDataBytes},
		// 	},
		// 	expected: false,
		// },
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			fr := &FROST{}
			got, err := fr.processBlameTypeInconsistentMessage(1, test.blameMessage)
			if err != nil {
				t.Error(err)
			}

			if got != test.expected {
				t.Fatalf("expected %t got %t", test.expected, got)
			}
		})
	}
}

func TestProcessBlameTypeInvalidShare(t *testing.T) {
	// Test with valid share
	fldmn, _ := sharing.NewFeldman(2, 4, thisCurve)
	verifiers, shares, _ := fldmn.Split(thisCurve.Scalar.Random(crand.Reader), crand.Reader)

	commitments := make([][]byte, 0)
	for _, commitment := range verifiers.Commitments {
		commitments = append(commitments, commitment.ToAffineCompressed())
	}

	sessionSK, _ := ecies.GenerateKey()
	operatorShare := shares[0] // share for operatorID 1
	eShare, _ := ecies.Encrypt(sessionSK.PublicKey, operatorShare.Value)

	round1Message := &Round1Message{
		Commitment: commitments,
		Shares: map[uint32][]byte{
			1: eShare,
		},
	}
	round1Bytes, _ := json.Marshal(round1Message)

	blameData := make([][]byte, 0)
	blameData = append(blameData, round1Bytes)

	blameMessage := &BlameMessage{
		Type:             InvalidShare,
		TargetOperatorID: 1,
		BlameData:        blameData,
		BlamerSessionSk:  sessionSK.Bytes(),
	}

	network := testingutils.NewTestingNetwork()
	dkgsigner := testingutils.NewTestingKeyManager()
	storage := testingutils.NewTestingStorage()

	kgp := New(network, 2, getRandRequestID(), dkgsigner, storage)
	fr := kgp.(*FROST)

	valid, err := fr.processBlameTypeInvalidShare(1, blameMessage)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, valid, true)

	// Test with invalid share
	invalidShare := shares[2].Value
	eInvalidShare, _ := ecies.Encrypt(sessionSK.PublicKey, invalidShare)
	round1Message.Shares[1] = eInvalidShare

	round1Bytes, _ = json.Marshal(round1Message)
	blameData[0] = round1Bytes
	blameMessage.BlameData = blameData

	valid, err = fr.processBlameTypeInvalidShare(1, blameMessage)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, valid, false)

}
