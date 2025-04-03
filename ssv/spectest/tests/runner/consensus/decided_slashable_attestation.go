package consensus

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DecidedSlashableAttestation tests that a slashable attestation is not signed
func DecidedSlashableAttestation() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	return &tests.MsgProcessingSpecTest{
		Name:             "decide on slashable attestation",
		Runner:           testingutils.CommitteeRunner(ks),
		Duty:             testingutils.TestingAttesterDuty(spec.DataVersionPhase0),
		Messages:         testingutils.SSVDecidingMsgsForCommitteeRunner(&testingutils.TestBeaconVote, ks, testingutils.TestingDutySlot),
		DecidedSlashable: true,
		OutputMessages:   []*types.PartialSignatureMessages{},
		ExpectedError:    "failed processing consensus message: decided ValidatorConsensusData invalid: decided value is invalid: slashable attestation",
	}
}
