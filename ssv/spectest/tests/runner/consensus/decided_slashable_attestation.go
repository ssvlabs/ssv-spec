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

	return tests.NewMsgProcessingSpecTest(
		"decide on slashable attestation",
		testingutils.CommitteeRunner(ks),
		testingutils.TestingAttesterDuty(spec.DataVersionPhase0),
		testingutils.SSVDecidingMsgsForCommitteeRunner(&testingutils.TestBeaconVote, ks, testingutils.TestingDutySlot),
		true,
		"",
		nil,
		[]*types.PartialSignatureMessages{},
		[]string{},
		false,
		"failed processing consensus message: decided ValidatorConsensusData invalid: decided value is invalid: slashable attestation",
	)
}
