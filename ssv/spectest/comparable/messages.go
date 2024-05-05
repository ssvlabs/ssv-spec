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
			if container.Signatures[root] == nil {
				container.Signatures[root] = map[types.OperatorID][]byte{}
			}
			container.Signatures[root][partialSigMsg.Signer] = partialSigMsg.PartialSignature
		}
	}
	return container
}
