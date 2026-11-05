package valcheckattestations

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CrossIndexEquivocation tests that the decided attestation index goes into the slashability data (SIP #94 §2):
// with a payload-present vote signed for the slot, re-voting it passes, while a payload-absent vote differing
// only in the index is a double vote.
func CrossIndexEquivocation() tests.SpecTest {
	keySet := testingutils.Testing4SharesSet()
	sharePKBytes := keySet.Shares[1].Serialize()
	slot := phase0.Slot(testingutils.TestingDutySlot)

	payloadPresent := testingutils.TestBeaconVote // AttestationDataIndex = 1
	payloadAbsent := payloadPresent
	payloadAbsent.AttestationDataIndex = 0

	encode := func(bv types.BeaconVote) []byte {
		byts, err := bv.Encode()
		if err != nil {
			panic(err.Error())
		}
		return byts
	}
	signed := map[string][]*phase0.AttestationData{
		hex.EncodeToString(sharePKBytes): {{
			Slot:            slot,
			Index:           payloadPresent.AttestationDataIndex,
			BeaconBlockRoot: payloadPresent.BlockRoot,
			Source:          payloadPresent.Source,
			Target:          payloadPresent.Target,
		}},
	}

	return valcheck.NewMultiSpecTest(
		"attestation value check cross-index equivocation",
		testdoc.ValCheckAttestationCrossIndexEquivocationDoc,
		[]*valcheck.SpecTest{
			{
				Name:               "same vote",
				Network:            types.BeaconTestNetwork,
				RunnerRole:         types.RoleCommittee,
				DutySlot:           slot,
				Input:              encode(payloadPresent),
				ExpectedSource:     *payloadPresent.Source,
				ExpectedTarget:     *payloadPresent.Target,
				ShareValidatorsPK:  []types.ShareValidatorPK{sharePKBytes},
				SignedAttestations: signed,
			},
			{
				Name:               "other index",
				Network:            types.BeaconTestNetwork,
				RunnerRole:         types.RoleCommittee,
				DutySlot:           slot,
				Input:              encode(payloadAbsent),
				ExpectedSource:     *payloadAbsent.Source,
				ExpectedTarget:     *payloadAbsent.Target,
				ShareValidatorsPK:  []types.ShareValidatorPK{sharePKBytes},
				SignedAttestations: signed,
				ExpectedErrorCode:  types.SlashableAttestationErrorCode,
			},
		},
	)
}
