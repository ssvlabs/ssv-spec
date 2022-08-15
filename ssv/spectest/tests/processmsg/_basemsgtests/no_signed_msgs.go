package _basemsgtests

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoSignedMsgs tests 0 SignedPostConsensusMessage messages
func NoSignedMsgs() *baseTest {
	ks := testingutils.Testing4SharesSet()

	return &baseTest{
		Name: "no signed messages",
		Msgs: []struct {
			MsgSignerSKs        []byte
			MsgSignerIDs        types.OperatorID
			BeaconRoots         [][]byte
			BeaconSignerSKs     [][]byte
			BeaconRootSignerIDs []types.OperatorID
			Slots               []spec.Slot
		}{
			{
				MsgSignerSKs:        ks.Shares[1].Serialize(),
				MsgSignerIDs:        1,
				BeaconRoots:         [][]byte{},
				BeaconSignerSKs:     [][]byte{},
				BeaconRootSignerIDs: []types.OperatorID{},
				Slots:               []spec.Slot{},
			},
		},
		ExpectedBaseError: "SignedPartialSignatureMessage invalid: no PartialSignatureMessages messages",
	}
}
