package keysign

import (
	"github.com/bloxapp/ssv-spec/dkg/common"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func Testing_PreparationMessageBytes(id types.OperatorID, ks *testingutils.TestKeySet, root []byte) []byte {
	msg := &ProtocolMsg{
		Round: common.Preparation,
		PreparationMessage: &PreparationMessage{
			PartialSignature: ks.Shares[id].SignByte(root).Serialize(),
		},
	}
	byts, _ := msg.Encode()
	return byts
}

func Testing_Round1MessageBytes(id types.OperatorID, ks *testingutils.TestKeySet, root []byte) []byte {
	msg := ProtocolMsg{
		Round: common.Round1,
		Round1Message: &Round1Message{
			ReconstructedSignature: ks.ValidatorSK.SignByte(root).Serialize(),
		},
	}
	byts, _ := msg.Encode()
	return byts
}
