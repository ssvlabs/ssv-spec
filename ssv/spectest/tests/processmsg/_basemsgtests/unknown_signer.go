package _basemsgtests

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// UnknownSigner tests Unknown SignedPostConsensusMessage signer
func UnknownSigner() *baseTest {
	ks := testingutils.Testing4SharesSet()

	domain, _ := types.ComputeETHDomain(types.DomainRandao, types.GenesisForkVersion, types.GenesisValidatorsRoot)
	r, _ := types.ComputeETHSigningRoot(types.SSZUint64(testingutils.TestingDutyEpoch), domain)

	return &baseTest{
		Name: "unknown signer",
		Msgs: []struct {
			MsgSignerSKs        []byte
			MsgSignerIDs        types.OperatorID
			BeaconRoots         [][]byte
			BeaconSignerSKs     [][]byte
			BeaconRootSignerIDs []types.OperatorID
			Slots               []spec.Slot
		}{
			{
				MsgSignerSKs: ks.Shares[1].Serialize(),
				MsgSignerIDs: 10,
				BeaconRoots: [][]byte{
					r[:],
				},
				BeaconSignerSKs: [][]byte{
					ks.Shares[1].Serialize(),
				},
				BeaconRootSignerIDs: []types.OperatorID{1},
				Slots:               []spec.Slot{testingutils.TestingDutySlot},
			},
		},
		ExpectedBaseError: "failed to verify PartialSignature: signer not found in operators",
	}
}
