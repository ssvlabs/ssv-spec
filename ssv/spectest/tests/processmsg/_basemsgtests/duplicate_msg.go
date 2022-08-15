package _basemsgtests

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateMsg tests duplicate SignedPostConsensusMessage msgs
func DuplicateMsg() *baseTest {
	ks := testingutils.Testing4SharesSet()

	domain, _ := types.ComputeETHDomain(types.DomainRandao, types.GenesisForkVersion, types.GenesisValidatorsRoot)
	r, _ := types.ComputeETHSigningRoot(types.SSZUint64(testingutils.TestingDutyEpoch), domain)

	return &baseTest{
		Name: "duplicate message",
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
				MsgSignerIDs: 1,
				BeaconRoots: [][]byte{
					r[:],
				},
				BeaconSignerSKs: [][]byte{
					ks.Shares[1].Serialize(),
				},
				BeaconRootSignerIDs: []types.OperatorID{1},
				Slots:               []spec.Slot{testingutils.TestingDutySlot},
			},
			{
				MsgSignerSKs: ks.Shares[1].Serialize(),
				MsgSignerIDs: 1,
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
		msgAccepted: true,
	}
}
