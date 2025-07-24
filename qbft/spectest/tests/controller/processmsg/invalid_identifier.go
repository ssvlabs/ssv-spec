package processmsg

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidIdentifier tests a process msg with the wrong identifier
func InvalidIdentifier() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	invalidPK := make([]byte, 32)
	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, invalidPK, types.RoleCommittee)

	test := tests.NewControllerSpecTest(
		"invalid identifier",
		testdoc.ControllerProcessMsgInvalidIdentifierDoc,
		[]*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					testingutils.SignQBFTMsg(ks.OperatorKeys[1], 1, &qbft.Message{
						MsgType:    qbft.ProposalMsgType,
						Height:     qbft.FirstHeight,
						Round:      qbft.FirstRound,
						Identifier: identifier[:],
						Root:       testingutils.TestingQBFTRootData,
					}),
				},
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: nil,
				},
			},
		},
		nil,
		"invalid msg: message doesn't belong to Identifier",
		nil,
	)

	test.SetPrivateKeys(ks)

	return test
}
