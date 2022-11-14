package frost

import (
	"encoding/hex"
	"sort"
	"testing"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/frost"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type MessagesForNodes map[uint32][]*dkg.SignedMessage

type FrostSpecTest struct {
	Name   string
	Keyset *testingutils.TestKeySet

	// Keygen Options
	RequestID dkg.RequestID
	Threshold uint64
	Operators []types.OperatorID

	// Resharing Options
	IsResharing       bool
	OperatorsOld      []types.OperatorID
	OldKeygenOutcomes testingutils.TestOutcome

	// Expected
	ExpectedOutcome testingutils.TestOutcome
	ExpectedError   string

	InputMessages map[int]MessagesForNodes
}

func (test *FrostSpecTest) TestName() string {
	return test.Name
}

func (test *FrostSpecTest) Run(t *testing.T) {

	outcomes, blame, err := test.TestingFrost()

	if len(test.ExpectedError) > 0 {
		require.EqualError(t, err, test.ExpectedError)
	} else {
		require.NoError(t, err)
	}

	if blame != nil {
		require.Equal(t, test.ExpectedOutcome.BlameOutcome.Valid, blame.Valid)
		return
	}

	for _, operatorID := range test.Operators {

		outcome := outcomes[uint32(operatorID)]
		if outcome.ProtocolOutput != nil {
			vk := hex.EncodeToString(outcome.ProtocolOutput.ValidatorPK)
			sk := outcome.ProtocolOutput.Share.SerializeToHexStr()
			pk := outcome.ProtocolOutput.OperatorPubKeys[operatorID].SerializeToHexStr()

			t.Logf("printing outcome keys for operatorID %d\n", operatorID)
			t.Logf("vk %s\n", vk)
			t.Logf("sk %s\n", sk)
			t.Logf("pk %s\n", pk)

			require.Equal(t, test.ExpectedOutcome.KeygenOutcome.ValidatorPK, vk)
			require.Equal(t, test.ExpectedOutcome.KeygenOutcome.Share[uint32(operatorID)], sk)
			require.Equal(t, test.ExpectedOutcome.KeygenOutcome.OperatorPubKeys[uint32(operatorID)], pk)
		}
	}
}

func (test *FrostSpecTest) TestingFrost() (map[uint32]*dkg.ProtocolOutcome, *dkg.BlameOutput, error) {

	testingutils.ResetRandSeed()
	dkgsigner := testingutils.NewTestingKeyManager()
	storage := testingutils.NewTestingStorage()
	network := testingutils.NewTestingNetwork()

	nodes := test.TestingFrostNodes(
		test.RequestID,
		network,
		storage,
		dkgsigner,
	)

	alloperators := test.Operators
	if test.IsResharing {
		alloperators = append(alloperators, test.OperatorsOld...)
	}

	initMessages, exists := test.InputMessages[0]
	if !exists {
		return nil, nil, errors.New("init messages not found in spec")
	}

	for operatorID, messages := range initMessages {
		for _, message := range messages {

			messageBytes, _ := message.Encode()
			startMessage := &types.SSVMessage{
				MsgType: types.DKGMsgType,
				Data:    messageBytes,
			}
			if err := nodes[types.OperatorID(operatorID)].ProcessMessage(startMessage); err != nil {
				return nil, nil, errors.Wrapf(err, "failed to start dkg protocol for operator %d", operatorID)
			}
		}
	}

	for round := 1; round <= 5; round++ {

		messages := network.BroadcastedMsgs
		network.BroadcastedMsgs = make([]*types.SSVMessage, 0)
		for _, msg := range messages {

			dkgMsg := &dkg.SignedMessage{}
			if err := dkgMsg.Decode(msg.Data); err != nil {
				return nil, nil, err
			}

			msgsToBroadcast := []*types.SSVMessage{}
			if testMessages, ok := test.InputMessages[round][uint32(dkgMsg.Signer)]; ok {
				for _, testMessage := range testMessages {
					testMessageBytes, _ := testMessage.Encode()
					msgsToBroadcast = append(msgsToBroadcast, &types.SSVMessage{
						MsgType: msg.MsgType,
						Data:    testMessageBytes,
					})
				}
			} else {
				msgsToBroadcast = append(msgsToBroadcast, msg)
			}

			operatorList := alloperators
			if test.IsResharing && round > 2 {
				operatorList = test.Operators
			}

			sort.SliceStable(operatorList, func(i, j int) bool {
				return operatorList[i] < operatorList[j]
			})

			for _, operatorID := range operatorList {

				if operatorID == dkgMsg.Signer {
					continue
				}

				for _, msgToBroadcast := range msgsToBroadcast {
					if err := nodes[operatorID].ProcessMessage(msgToBroadcast); err != nil {
						return nil, nil, err
					}
				}
			}
		}

	}

	ret := make(map[uint32]*dkg.ProtocolOutcome)

	outputs := network.DKGOutputs
	blame := network.BlameOutput
	if blame != nil {
		return nil, blame, nil
	}

	for operatorID, output := range outputs {
		if output.BlameData != nil {
			signedMsg := &dkg.SignedMessage{}
			_ = signedMsg.Decode(output.BlameData.BlameMessage)
			ret[uint32(operatorID)] = &dkg.ProtocolOutcome{
				BlameOutput: &dkg.BlameOutput{
					Valid:        output.BlameData.Valid,
					BlameMessage: signedMsg,
				},
			}
			continue
		}

		pk := &bls.PublicKey{}
		_ = pk.Deserialize(output.Data.SharePubKey)

		share, _ := dkgsigner.Decrypt(test.Keyset.DKGOperators[operatorID].EncryptionKey, output.Data.EncryptedShare)
		sk := &bls.SecretKey{}
		_ = sk.Deserialize(share)

		ret[uint32(operatorID)] = &dkg.ProtocolOutcome{
			ProtocolOutput: &dkg.KeyGenOutput{
				ValidatorPK: output.Data.ValidatorPubKey,
				Share:       sk,
				OperatorPubKeys: map[types.OperatorID]*bls.PublicKey{
					operatorID: pk,
				},
				Threshold: test.Threshold,
			},
		}
	}

	return ret, nil, nil
}

func (test *FrostSpecTest) TestingFrostNodes(
	requestID dkg.RequestID,
	network dkg.Network,
	storage dkg.Storage,
	dkgsigner types.DKGSigner,
) map[types.OperatorID]*dkg.Node {

	nodes := make(map[types.OperatorID]*dkg.Node)
	for _, operatorID := range test.Operators {
		_, operator, _ := storage.GetDKGOperator(operatorID)
		node := dkg.NewNode(
			operator,
			&dkg.Config{
				KeygenProtocol: frost.New,
				Network:        network,
				Storage:        storage,
				Signer:         dkgsigner,
			})
		nodes[operatorID] = node
	}

	if test.IsResharing {
		operatorsOldList := types.OperatorList(test.OperatorsOld).ToUint32List()
		keygenOutcomeOld := test.OldKeygenOutcomes.KeygenOutcome.ToKeygenOutcomeMap(test.Threshold, operatorsOldList)

		for _, operatorID := range test.OperatorsOld {
			storage := testingutils.NewTestingStorage()
			_ = storage.SaveKeyGenOutput(keygenOutcomeOld[uint32(operatorID)])

			_, operator, _ := storage.GetDKGOperator(operatorID)
			node := dkg.NewResharingNode(
				operator,
				test.OperatorsOld,
				&dkg.Config{
					ReshareProtocol: frost.NewResharing,
					Network:         network,
					Storage:         storage,
					Signer:          dkgsigner,
				})
			nodes[operatorID] = node
		}

		for _, operatorID := range test.Operators {
			storage := testingutils.NewTestingStorage()
			_ = storage.SaveKeyGenOutput(keygenOutcomeOld[operatorsOldList[0]])

			_, operator, _ := storage.GetDKGOperator(operatorID)
			node := dkg.NewResharingNode(
				operator,
				test.OperatorsOld,
				&dkg.Config{
					ReshareProtocol: frost.NewResharing,
					Network:         network,
					Storage:         storage,
					Signer:          dkgsigner,
				})
			nodes[operatorID] = node
		}
	}
	return nodes
}
