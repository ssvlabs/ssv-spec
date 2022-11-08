package frost

import (
	"encoding/hex"
	"testing"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var (
	expectedFrostOutput = testingutils.TestKeygenOutcome{
		Share: map[uint32]string{
			1: "13b4682b21fe50088beff43530787d1dac1e50c8e0686ec55849c8c9c9c5c044",
			2: "2b2becc7a00babd145cd75772126d9a8a10f3ca975ad4fa862abe06f4c7e8b59",
			3: "5bfa5c11d0a68028a94abf35e43955c2c7f4f4d18ef62fb9eed6b905dfe4d2ef",
			4: "32320eb68a314fc6832df969700e1966cd11d53e2c44b2fafcca528e83f89705",
		},
		ValidatorPK: "ab2caf206286eb161d47124885b05b0e92d5d77ba29ce7aa77d9cd38ea24cfc6198d037d7b2011388b475c24ab40091e",
		OperatorPubKeys: map[uint32]string{
			1: "b9dbb91742532eb1e8641491bd3a2ee149584d4d6c68169daad84addfa848088c38c3c6302abbcb4f648441b0c67c6e4",
			2: "ad68795bfe98239f64eaaea753ad6cb5fbdc51fdecf1b42abcee65906eabe4f376b1fd85dbc11e15bf0b04d28fbda199",
			3: "8ec8fa19ece71538a6435a9784d7565496c57ffbaa1160a020ca14e4c64bdae5d6073bdb43bad401fec264b4cc554295",
			4: "ad67ab94ab4f560a414f3fdc7b15bf0cf091ff72791c37515631a14c6446462a93b0559851238f28550cb82ab0808e22",
		},
	}
)

func TestFrostDKG(t *testing.T) {

	operators := []types.OperatorID{
		1, 2, 3, 4,
	}

	outputs, err := TestingFrost(
		3,
		operators,
		nil,
		false,
		nil,
	)
	if err != nil {
		t.Error(err)
	}

	for _, operatorID := range operators {
		output := outputs[uint32(operatorID)].ProtocolOutput
		require.Equal(t, expectedFrostOutput.ValidatorPK, hex.EncodeToString(output.ValidatorPK))
		require.Equal(t, expectedFrostOutput.Share[uint32(operatorID)], output.Share.SerializeToHexStr())
		for opID, publicKey := range output.OperatorPubKeys {
			require.Equal(t, expectedFrostOutput.OperatorPubKeys[uint32(opID)], publicKey.SerializeToHexStr())
		}
	}
}

func TestResharing(t *testing.T) {

	tests := map[string]struct {
		input, expected testingutils.TestKeygenOutcome
	}{
		"test_1": {
			input: testingutils.TestKeygenOutcome{
				ValidatorPK: "8ae6e5255472e548d039d5333001c66109c9896473d99c56f64f8e27da1bd8f645ec4e6a0c576b78c722896bce372812",
				Share: map[uint32]string{
					1: "5365b83d582c9d1060830fa50a958df9f7e287e9860a70c97faab36a06be2912",
					2: "533959ffa931481f392b2e86e203410fb1245436588db34dde389456dc0251b7",
					3: "442f11f780536f53eda21438cda8c1835eccc54c4473d77b158d006f99044186",
					4: "2646e024dd9312ae7de7c0bacd860f5500dbdb2b49bcdd5125a7f7b43dc3f87f",
				},
				OperatorPubKeys: map[uint32]string{
					1: "add523513d851787ec611256fe759e21ee4e84a684bc33224973a5481b202061bf383fac50319ce1f903207a71a4d8fa",
					2: "8b9dfd049985f0aa84a8c309914df6752f32803c3b5590b279b1c24dba5b83f574ea6dba3038f55275d62a4f25a11cf5",
					3: "b31e1a5da47be70788ebfdc4ec162b9dff1fe2d177af9187af41b472f10ecd0a90f9d9834be6103ce4690a36f25fe051",
					4: "a9697dea52e229d8171a3051514df7a491e1228d8208f0561538e06f138dd37ddd6e0f7e3975cadf159bc2a02819d037",
				},
			},
			expected: testingutils.TestKeygenOutcome{
				ValidatorPK: "8ae6e5255472e548d039d5333001c66109c9896473d99c56f64f8e27da1bd8f645ec4e6a0c576b78c722896bce372812",
				Share: map[uint32]string{
					5: "437a5713f74cbfc67bf2781dd6c8cb74db8bb0a1e598b9bb372e84c56da99dad",
					6: "72475d8d45ee21f23d27a30f839782814114a24e5fe356bc9301af2289127ee5",
					7: "0e7645e7292c9aabd34be54e52f4bd42ae63932c2370ad076fc0e50496ec1460",
					8: "73d0061b1de0a1cbd80cc6f261c603c91eb16f44303bd098cd6c266897365e21",
				},
				OperatorPubKeys: map[uint32]string{
					5: "a7bc4dbb7be0e2ce4b5f6121d7a0fc2902eb9abe4c8b17875cc74a5ce3bd61784eeca8561ce169b2f92e5a157d8d0e49",
					6: "83678e9ffc680a5ba164226eae678b575c87a8d67696c0c997ae0ad283a942cad21e576994db93c5b55ce42852e0fe87",
					7: "8f3dc952da48c29313099f36d11b66b307896ab7f8a8396e0ac7b64908533641818241bf963b04c4113d3c7f01d6e72b",
					8: "8153e8b4f1820e625cc695c7fde34e9f4498e889e9efcc4b86f2ba645d3026df284a6d133ac28f40e6569f658bf4c39c",
				},
			},
		},
		"test_2": {
			input: expectedFrostOutput,
			expected: testingutils.TestKeygenOutcome{
				Share: map[uint32]string{
					5: "1459f89fc085be6d93fafbfaa197635491a6816be7d636de0cfb9931ac2c47c6",
					6: "4326ff190f272099553026ec4e661a60f72f73186220d3df68cec38ec79528fe",
					7: "53438ec61c03169b1e8e413327652d27b83c07f925ac8629458df96fd56ebe7a",
					8: "44afa7a6e719a072f0154acf2c949ba8d4cc400e32794dbba3393ad4d5b9083a",
				},
				ValidatorPK: "ab2caf206286eb161d47124885b05b0e92d5d77ba29ce7aa77d9cd38ea24cfc6198d037d7b2011388b475c24ab40091e",
				OperatorPubKeys: map[uint32]string{
					5: "a841b23b62457fe942389e5289b78848c4c7c709623da1e769ed2d5f7197f831f926bbba945940cb80e27b476e33a003",
					6: "b4a227e0baf87ce87ae7dd773615d04bd0732a14bc3db4db8f3a68c02ad639d1b1ed4420b63927fa37da7aa4fd66eef6",
					7: "84df18d6eef01054a3cb21104570b261b2705ee6ff82ed542f0120bc558e359255e3cf22d6280995983e5f5a216313ea",
					8: "8746cd13d7b9f4c472f7e851a23eb30dca95b5fd02bccdfbe9e27742eed30c1a0a9dda775079016701e4e56362c5049b",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			operatorsOld := []types.OperatorID{
				1, 2, 3, // 4,
			}

			operators := []types.OperatorID{
				5, 6, 7, 8,
			}

			outcomes, err := TestingFrost(
				3,
				operators,
				operatorsOld,
				true,
				&test.input,
			)
			if err != nil {
				t.Fatalf("failed to run frost: %s", err.Error())
			}

			for _, operatorID := range operators {
				outcome := outcomes[uint32(operatorID)].ProtocolOutput

				require.Equal(t, test.expected.ValidatorPK, hex.EncodeToString(outcome.ValidatorPK))
				require.Equal(t, test.expected.Share[uint32(operatorID)], outcome.Share.SerializeToHexStr())
				for opID, publicKey := range outcome.OperatorPubKeys {
					require.Equal(t, test.expected.OperatorPubKeys[uint32(opID)], publicKey.SerializeToHexStr())
				}
			}
		})
	}
}

func TestingFrost(
	threshold uint64,
	operators, operatorsOld []types.OperatorID,
	isResharing bool,
	oldKeygenOutcomes *testingutils.TestKeygenOutcome,
) (map[uint32]*dkg.ProtocolOutcome, error) {

	testingutils.ResetRandSeed()
	requestID := testingutils.GetRandRequestID()
	dkgsigner := testingutils.NewTestingKeyManager()
	storage := testingutils.NewTestingStorage()
	network := testingutils.NewTestingNetwork()

	init := &dkg.Init{
		OperatorIDs: operators,
		Threshold:   uint16(threshold),
	}

	kgps := make(map[types.OperatorID]dkg.Protocol)
	for _, operatorID := range operators {
		p := New(network, operatorID, requestID, dkgsigner, storage, init)
		kgps[operatorID] = p
	}

	if isResharing {
		operatorsOldList := types.OperatorList(operatorsOld).ToUint32List()
		keygenOutcomeOld := oldKeygenOutcomes.ToKeygenOutcomeMap(threshold, operatorsOldList)

		reshare := &dkg.Reshare{
			ValidatorPK: keygenOutcomeOld[operatorsOldList[0]].ValidatorPK,
			OperatorIDs: operators,
			Threshold:   uint16(threshold),
		}

		for _, operatorID := range operatorsOld {
			p := NewResharing(network, operatorID, requestID, dkgsigner, storage, operatorsOld, reshare, keygenOutcomeOld[uint32(operatorID)])
			kgps[operatorID] = p

		}

		for _, operatorID := range operators {
			p := NewResharing(network, operatorID, requestID, dkgsigner, storage, operatorsOld, reshare, nil)
			kgps[operatorID] = p
		}
	}

	alloperators := operators
	if isResharing {
		alloperators = append(alloperators, operatorsOld...)
	}

	for _, operatorID := range alloperators {
		if err := kgps[operatorID].Start(); err != nil {
			return nil, errors.Wrapf(err, "failed to start dkg protocol for operator %d", operatorID)
		}
	}

	outcomes := make(map[uint32]*dkg.ProtocolOutcome)
	for i := 0; i < 3; i++ {

		messages := network.BroadcastedMsgs
		network.BroadcastedMsgs = make([]*types.SSVMessage, 0)

		for _, msg := range messages {

			dkgMsg := &dkg.SignedMessage{}
			if err := dkgMsg.Decode(msg.Data); err != nil {
				return nil, err
			}

			for _, operatorID := range alloperators {

				if operatorID == dkgMsg.Signer {
					continue
				}

				finished, outcome, err := kgps[operatorID].ProcessMsg(dkgMsg)
				if err != nil {
					return nil, err
				}
				if finished {
					outcomes[uint32(operatorID)] = outcome
				}
			}
		}
	}

	return outcomes, nil
}
