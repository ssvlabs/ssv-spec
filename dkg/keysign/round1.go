package keysign

import (
	"fmt"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/common"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

func (instance *Instance) processRound1() (bool, *dkg.ProtocolOutcome, error) {
	if !instance.canProceedThisRound() {
		return false, nil, nil
	}
	instance.state.SetCurrentRound(common.Round1)
	instance.state.roundTimer.StartRoundTimeoutTimer(instance.state.GetCurrentRound())

	prevRoundMsgs := instance.state.msgContainer.AllMessagesForRound(common.Preparation)
	validSignatureBytes := map[types.OperatorID][]byte{}
	totalValidSignatures := uint64(0)
	for operatorID, dkgMessage := range prevRoundMsgs {
		protocolMessage := &ProtocolMsg{}
		if err := protocolMessage.Decode(dkgMessage.Message.Data); err != nil {
			return false, nil, errors.Wrap(err, "failed to decode protocol msg")
		}

		sharePK, found := instance.instanceParams.OperatorPublicKeyshares[types.OperatorID(operatorID)]
		if !found {
			return false, nil, fmt.Errorf("public keyshare for operator %d not found in keygen output", operatorID)
		}

		sig := &bls.Sign{}
		if err := sig.Deserialize(protocolMessage.PreparationMessage.PartialSignature); err != nil {
			return false, nil, fmt.Errorf("failed to deserialize bls signature")
		}

		if !sig.VerifyByte(sharePK, instance.instanceParams.SigningRoot) {
			return false, nil, fmt.Errorf("invalid partial signature from operator %d", operatorID)
		}

		validSignatureBytes[types.OperatorID(operatorID)] = protocolMessage.PreparationMessage.PartialSignature
		totalValidSignatures += 1
	}

	if totalValidSignatures < instance.instanceParams.threshold {
		return false, nil, fmt.Errorf("insufficient number of valid signatures to reconstruct full signature: expected %d got %d", totalValidSignatures, instance.instanceParams.threshold)
	}

	reconstructedSig, err := types.ReconstructSignatures(validSignatureBytes)
	if err != nil {
		return false, nil, fmt.Errorf("failed to reconstruct full bls signature from t partial signatures")
	}

	msg := &ProtocolMsg{
		Round: common.Round1,
		Round1Message: &Round1Message{
			ReconstructedSignature: reconstructedSig.Serialize(),
		},
	}

	bcastMsg, err := instance.saveSignedMsg(msg)
	if err != nil {
		return false, nil, err
	}
	return false, nil, instance.config.GetNetwork().BroadcastDKGMessage(bcastMsg)
}
