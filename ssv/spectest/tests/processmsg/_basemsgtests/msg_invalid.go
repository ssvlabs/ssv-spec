package _basemsgtests

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidMsg tests an invalid SignedPostConsensusMessage
func InvalidMsg() *baseTest {
	ks := testingutils.Testing4SharesSet()

	return &baseTest{
		Name: "valid msg",
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
					nil,
				},
				BeaconSignerSKs: [][]byte{
					ks.Shares[1].Serialize(),
				},
				BeaconRootSignerIDs: []types.OperatorID{1},
				Slots:               []spec.Slot{testingutils.TestingDutySlot},
			},
		},
		ExpectedBaseError: "SignedPartialSignatureMessage invalid: message invalid: SigningRoot invalid",
	}
}
