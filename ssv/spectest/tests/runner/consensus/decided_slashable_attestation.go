package consensus

import (
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
		Duty:             testingutils.TestingAttesterDuty,
		Messages:         testingutils.SSVDecidingMsgsForCommitteeRunner(&testingutils.TestBeaconVote, ks, testingutils.TestingDutySlot),
		DecidedSlashable: true,
		OutputMessages:   []*types.PartialSignatureMessages{},
		ExpectedError:    "attempting to sign slashable attestation data: slashable attestation",
	}
}
