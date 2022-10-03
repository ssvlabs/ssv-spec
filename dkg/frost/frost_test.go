package frost

import (
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"math/big"
	mrand "math/rand"
	"testing"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/coinbase/kryptology/pkg/sharing"
	ecies "github.com/ecies/go/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func resetRandSeed() {
	src := mrand.NewSource(1)
	src.Seed(12345)
	crand.Reader = mrand.New(src)
}

func TestFrostDKG(t *testing.T) {

	resetRandSeed()

	type ExpectedKGOutput struct {
		Share           map[uint32]string
		ValidatorPK     string
		OperatorPubKeys map[uint32]string
	}

	expected := ExpectedKGOutput{
		Share: map[uint32]string{
			1: "285a26f43b026b246ca0c33b34aaf90890c016d943a75456efbe00d4d0bdee01",
			2: "1d3701ab6e7b902bd482ac899ec7bab1852376ae234474bae1a3f83bb41dc48f",
			3: "42afa077e46dd25be4d7bb5be8734e77df5f074e0933f6ef6af8bdbe3e205cd0",
			4: "67c262ae06e14097b7b3e5a1a36526d6640ac899407bf61fd38c3490e43afed4",
		},
		ValidatorPK: "84d633334d8d615d6739d1f011f2c9b194601e38213937999868ed9b016cab8500e16319a2866ed853411ce1628e84b3",
		OperatorPubKeys: map[uint32]string{
			1: "960498d1f66481d570b80c2cb6fbafa40a250f46510412eb51abaf1b62aa17e747c8c40f69d01606cd29dd0466f4a9a2",
			2: "a73f10841b40509f3a727a3311c77ee46ce0ae43ffdbd44aca87f837e392772834f51d1b38eacbe91d21057c0717a51b",
			3: "8982bd51c3a08d8eb0d470eeb57fe3a8a8db4f426026019bf27a5faa745fc13bc75e3e2bea2435f47fa9148313959000",
			4: "af4ce0c5ec16cc0d52acb5419d8b51051bcb271275680ab17c3a445d4de3c661971f19786667ab60216955bf20a13ea7",
		},
	}

	requestID := getRandRequestID()

	operators := []types.OperatorID{
		1, 2, 3, 4,
	}

	dkgsigner := testingutils.NewTestingKeyManager()
	storage := testingutils.NewTestingStorage()
	network := testingutils.NewTestingNetwork()

	kgps := make(map[types.OperatorID]dkg.KeyGenProtocol)
	for _, operatorID := range operators {
		p := New(network, operatorID, requestID, dkgsigner, storage)
		kgps[operatorID] = p
	}

	threshold := 2
	outputs := make(map[uint32]*dkg.KeyGenOutput)

	// preparation round
	initMsg := &dkg.Init{
		OperatorIDs: operators,
		Threshold:   uint16(threshold),
	}

	for _, operatorID := range operators {
		if err := kgps[operatorID].Start(initMsg); err != nil {
			t.Error(errors.Wrapf(err, "failed to start dkg protocol for operator %d", operatorID))
		}
	}

	rounds := []string{"round 1", "round 2", "keygen output"}

	for _, round := range rounds {
		t.Logf("proceeding with %s", round)

		messages := network.BroadcastedMsgs
		network.BroadcastedMsgs = make([]*types.SSVMessage, 0)

		for _, msg := range messages {
			dkgMsg := &dkg.SignedMessage{}
			if err := dkgMsg.Decode(msg.Data); err != nil {
				t.Error(err)
			}

			for _, operatorID := range operators {
				if operatorID == dkgMsg.Signer {
					continue
				}

				finished, output, err := kgps[operatorID].ProcessMsg(dkgMsg)
				if err != nil {
					t.Error(err)
				}

				if finished {
					outputs[uint32(operatorID)] = output
				}
			}
		}

	}

	for _, operatorID := range operators {
		output := outputs[uint32(operatorID)]

		require.Equal(t, expected.ValidatorPK, hex.EncodeToString(output.ValidatorPK))
		require.Equal(t, expected.Share[uint32(operatorID)], output.Share.SerializeToHexStr())
		for opID, publicKey := range output.OperatorPubKeys {
			require.Equal(t, expected.OperatorPubKeys[uint32(opID)], publicKey.SerializeToHexStr())
		}
	}
}

func getRandRequestID() dkg.RequestID {
	requestID := dkg.RequestID{}
	for i := range requestID {
		rndInt, _ := crand.Int(crand.Reader, big.NewInt(255))
		requestID[i] = rndInt.Bytes()[0]
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
