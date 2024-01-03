package validator

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func ValidPartialSig() tests.SpecTest {
	// KeySet
	ks := testingutils.Testing4SharesSet()

	// Message ID
	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleProposer)

	// Duty
	duty := testingutils.TestingProposerDutyV(spec.DataVersionCapella)

	// Randao
	preConsensus := testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionCapella)
	preConsensusByts, err := preConsensus.Encode()
	if err != nil {
		panic(err.Error())
	}

	// Messages
	msgs := []*types.SSVMessage{
		{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   msgID,
			Data:    preConsensusByts[:],
		},
	}

	return &ValidatorTest{
		Name:                   "valid partial sig",
		Duties:                 []*types.Duty{duty},
		Messages:               msgs,
		OutputMessages:         []*types.SSVMessage{},
		BeaconBroadcastedRoots: []string{},
	}
}
