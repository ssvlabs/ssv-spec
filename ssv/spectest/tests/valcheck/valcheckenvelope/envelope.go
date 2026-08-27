package valcheckenvelope

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// EnvelopeBlocks covers the Gloas (ePBS, SIP #94 §6) envelope value-check rules: the QBFT value must
// decode into an envelope-proposer duty for this validator carrying a blinded envelope whose builder
// index is the self-build sentinel and whose beacon block root equals the recorded §4-decided block
// root for the duty's slot — the linkage that is the reason §6 exists. The envelope content itself is
// leader-trusted (blinded-block trust model).
func EnvelopeBlocks() tests.SpecTest {
	slot := phase0.Slot(testingutils.TestingDutySlotGloas)

	encode := func(cd *types.EnvelopeConsensusData) []byte {
		byts, err := cd.Encode()
		if err != nil {
			panic(err.Error())
		}
		return byts
	}
	withEnvelope := func(mutate func(*gloas.BlindedExecutionPayloadEnvelope)) []byte {
		envelope := testingutils.TestingBlindedExecutionPayloadEnvelope(slot)
		mutate(envelope)
		byts, err := envelope.Encode()
		if err != nil {
			panic(err.Error())
		}
		cd := testingutils.TestingEnvelopeConsensusData(slot)
		cd.DataSSZ = byts
		return encode(cd)
	}

	wrongRoleDuty := testingutils.TestingEnvelopeProposerDuty()
	wrongRoleDuty.Type = types.BNRoleProposer
	wrongRoleCd := testingutils.TestingEnvelopeConsensusData(slot)
	wrongRoleCd.Duty = *wrongRoleDuty

	// The store is seeded for the Gloas duty slot only, so a duty at another slot has no recorded root.
	noRootCd := testingutils.TestingEnvelopeConsensusData(testingutils.TestingDutySlotGloasNextEpoch)

	return valcheck.NewMultiSpecTest(
		"envelope blocks",
		testdoc.ValCheckEnvelopeBlocksDoc,
		[]*valcheck.SpecTest{
			{
				Name:       "valid envelope",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleEnvelopeProposer,
				Input:      testingutils.TestingEnvelopeConsensusDataByts(slot),
			},
			{
				Name:              "undecodable consensus data",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleEnvelopeProposer,
				Input:             []byte("garbage"),
				ExpectedErrorCode: types.EnvelopeConsensusDataDecodeErrorCode,
			},
			{
				Name:       "undecodable envelope",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleEnvelopeProposer,
				Input: encode(func() *types.EnvelopeConsensusData {
					cd := testingutils.TestingEnvelopeConsensusData(slot)
					cd.DataSSZ = []byte("garbage")
					return cd
				}()),
				ExpectedErrorCode: types.UnmarshalSSZErrorCode,
			},
			{
				Name:              "wrong duty role",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleEnvelopeProposer,
				Input:             encode(wrongRoleCd),
				ExpectedErrorCode: types.WrongBeaconRoleTypeErrorCode,
			},
			{
				Name:       "wrong builder index",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleEnvelopeProposer,
				Input: withEnvelope(func(envelope *gloas.BlindedExecutionPayloadEnvelope) {
					envelope.BuilderIndex = gloas.BuilderIndexSelfBuild - 1
				}),
				ExpectedErrorCode: types.EnvelopeWrongBuilderIndexErrorCode,
			},
			{
				Name:       "beacon block root does not match the decided block",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleEnvelopeProposer,
				Input: withEnvelope(func(envelope *gloas.BlindedExecutionPayloadEnvelope) {
					envelope.BeaconBlockRoot = testingutils.TestingBlockRoot
				}),
				ExpectedErrorCode: types.EnvelopeBlockRootMismatchErrorCode,
			},
			{
				Name:              "no decided block root for the slot",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleEnvelopeProposer,
				Input:             encode(noRootCd),
				ExpectedErrorCode: types.EnvelopeNoProposedBlockRootErrorCode,
			},
		},
	)
}
