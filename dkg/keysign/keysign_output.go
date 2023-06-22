package keysign

import (
	"bytes"
	"fmt"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/common"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

func (instance *Instance) processKeysignOutput() (bool, *dkg.ProtocolOutcome, error) {
	if !instance.canProceedThisRound() {
		return false, nil, nil
	}
	instance.state.SetCurrentRound(common.KeysignOutput)
	instance.state.roundTimer.StartRoundTimeoutTimer(instance.state.GetCurrentRound())

	output := &dkg.ProtocolOutcome{
		KeySignOutput: &dkg.KeySignOutput{
			RequestID:   instance.instanceParams.identifier,
			ValidatorPK: instance.instanceParams.validatorPK,
		},
	}

	prevRoundMsgs := instance.state.msgContainer.AllMessagesForRound(common.Round1)
	var prevSignature []byte
	for _, dkgMessage := range prevRoundMsgs {
		protocolMessage := &ProtocolMsg{}
		if err := protocolMessage.Decode(dkgMessage.Message.Data); err != nil {
			return false, nil, errors.Wrap(err, "failed to decode protocol msg")
		}

		sig := &bls.Sign{}
		if err := sig.Deserialize(protocolMessage.Round1Message.ReconstructedSignature); err != nil {
			return false, nil, fmt.Errorf("failed to deserialize bls sig")
		}

		vk := &bls.PublicKey{}
		if err := vk.Deserialize(instance.instanceParams.validatorPK); err != nil {
			return false, nil, fmt.Errorf("failed to deserialize bls public key")
		}

		if !sig.VerifyByte(vk, instance.instanceParams.SigningRoot) {
			return false, nil, fmt.Errorf("invalid reconstructed bls signature")
		}

		if prevSignature != nil && !bytes.Equal(protocolMessage.Round1Message.ReconstructedSignature, prevSignature) {
			return false, nil, fmt.Errorf("inconsistent signatures between operators")
		}

		prevSignature = protocolMessage.Round1Message.ReconstructedSignature
	}

	output.KeySignOutput.Signature = prevSignature
	return true, output, nil
}
