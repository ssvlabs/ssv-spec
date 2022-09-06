package qbft

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

//// UponDecided returns true if a decided messages was received.
//func (i *Instance) UponDecided(signedDecided *SignedMessage, commitMsgContainer *MsgContainer) (bool, []byte, error) {
//	if i.State.Decided {
//		return true, i.State.DecidedValue, nil
//	}
//
//	if err := validateDecided(
//		i.State.Share,
//		i.config,
//		signedDecided,
//	); err != nil {
//		return false, nil, errors.Wrap(err, "invalid decided msg")
//	}
//
//	addMsg, err := commitMsgContainer.AddFirstMsgForSignerAndRound(signedDecided)
//	if err != nil {
//		return false, nil, errors.Wrap(err, "could not add commit msg to container")
//	}
//	if !addMsg {
//		return false, nil, nil // UponCommit was already called
//	}
//
//	msgDecidedData, err := signedDecided.Message.GetCommitData()
//	if err != nil {
//		return false, nil, errors.Wrap(err, "could not get msg decided data")
//	}
//
//	return true, msgDecidedData.Data, nil
//}

func validateDecided(
	share *types.Share,
	config IConfig,
	signedDecided *SignedMessage,
) error {
	if !isDecidedMsg(share, signedDecided) {
		return errors.New("not a decided msg")
	}

	operators := share.Committee
	if err := baseCommitValidation(config, signedDecided, signedDecided.Message.Height, operators); err != nil {
		return errors.Wrap(err, "invalid decided msg")
	}

	msgDecidedData, err := signedDecided.Message.GetCommitData()
	if err != nil {
		return errors.Wrap(err, "could not get msg decided data")
	}

	if err := config.GetValueCheckF()(msgDecidedData.Data); err != nil {
		return errors.Wrap(err, "decided value invalid")
	}

	return nil
}

// returns true if signed commit has all quorum sigs
func isDecidedMsg(share *types.Share, signedDecided *SignedMessage) bool {
	return share.HasQuorum(len(signedDecided.Signers))
}
