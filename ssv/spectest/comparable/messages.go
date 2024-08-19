package ssvcomparable

import (
	"encoding/hex"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
)

func SetMessagesInContainer(container *ssv.PartialSigContainer, messages []*types.SignedSSVMessage) *ssv.PartialSigContainer {
	for _, signedSSVMsg := range messages {
		if signedSSVMsg.SSVMessage.MsgType != types.SSVPartialSignatureMsgType {
			continue
		}

		msg := &types.PartialSignatureMessages{}
		if err := msg.Decode(signedSSVMsg.SSVMessage.Data); err != nil {
			panic(err.Error())
		}

		for _, partialSigMsg := range msg.Messages {
			root := hex.EncodeToString(partialSigMsg.SigningRoot[:])

			if container.Signatures[partialSigMsg.ValidatorIndex] == nil {
				container.Signatures[partialSigMsg.ValidatorIndex] = make(map[ssv.SigningRoot]map[uint64]types.Signature)
			}
			if container.Signatures[partialSigMsg.ValidatorIndex][ssv.SigningRoot(root)] == nil {
				container.Signatures[partialSigMsg.ValidatorIndex][ssv.SigningRoot(root)] = make(map[uint64]types.Signature)
			}
			container.Signatures[partialSigMsg.ValidatorIndex][ssv.SigningRoot(root)][partialSigMsg.Signer] = partialSigMsg.PartialSignature
		}
	}
	return container
}
