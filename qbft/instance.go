package qbft

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/bloxapp/ssv-spec/types"
)

type ProposedValueCheckF func(data []byte) error
type ProposerF func(state *State, round Round) types.OperatorID

// Instance is a single QBFT instance that starts with a Start call (including a value).
// Every new msg the ProcessMsg function needs to be called
type Instance struct {
	State  *State
	config IConfig

	processMsgF *types.ThreadSafeF
	startOnce   sync.Once
	// forceStop will force stop the instance if set to true
	forceStop  bool
	StartValue []byte
}

func NewInstance(
	config IConfig,
	share *types.Share,
	identifier []byte,
	height Height,
) *Instance {
	return &Instance{
		State: &State{
			Share:                share,
			ID:                   identifier,
			Round:                FirstRound,
			Height:               height,
			LastPreparedRound:    NoRound,
			ProposeContainer:     NewMsgContainer(),
			PrepareContainer:     NewMsgContainer(),
			CommitContainer:      NewMsgContainer(),
			RoundChangeContainer: NewMsgContainer(),
		},
		config:      config,
		processMsgF: types.NewThreadSafeF(),
	}
}

func (i *Instance) ForceStop() {
	i.forceStop = true
}

// Start is an interface implementation
func (i *Instance) Start(value []byte, height Height) {
	i.startOnce.Do(func() {
		i.StartValue = value
		i.State.Round = FirstRound
		i.State.Height = height

		i.config.GetTimer().TimeoutForRound(FirstRound)

		// propose if this node is the proposer
		if proposer(i.State, i.GetConfig(), FirstRound) == i.State.Share.OperatorID {
			proposal, err := CreateProposal(i.State, i.config, i.StartValue, nil, nil)
			// nolint
			if err != nil {
				fmt.Printf("%s\n", err.Error())
			}
			// nolint
			if err := i.Broadcast(proposal); err != nil {
				fmt.Printf("%s\n", err.Error())
			}
		}
	})
}

func (i *Instance) Broadcast(msg *Message) error {
	if !i.CanProcessMessages() {
		return errors.New("instance stopped processing messages")
	}
	byts, err := msg.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode message")
	}

	msgID := types.MessageID{}
	copy(msgID[:], msg.Identifier)

	ssvMessage := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   msgID,
		Data:    byts,
	}

	// Sign SSVMessage
	signingData, err := ssvMessage.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode SSVMessage")
	}
	signature, err := i.config.GetSigner().SignNetworkData(signingData, i.State.Share.NetworkPubkey)
	if err != nil {
		return errors.Wrap(err, "could not sign SSVMessage")
	}

	// Create SignedSSVMessage
	msgToBroadcast := &types.SignedSSVMessage{
		OperatorID: []types.OperatorID{i.State.Share.OperatorID},
		Signature:  [][]byte{signature},
		SSVMessage: ssvMessage,
	}

	return i.config.GetNetwork().Broadcast(msgToBroadcast)
}

// ProcessMsg processes a new QBFT msg, returns non nil error on msg processing error
func (i *Instance) ProcessMsg(signedSSVMessage *types.SignedSSVMessage) (decided bool, decidedValue []byte, aggregatedCommit *types.SignedSSVMessage, err error) {
	if !i.CanProcessMessages() {
		return false, nil, nil, errors.New("instance stopped processing messages")
	}

	if err := i.BaseMsgValidation(signedSSVMessage); err != nil {
		return false, nil, nil, errors.Wrap(err, "invalid signed message")
	}

	// Decode
	message := &Message{}
	if err := message.Decode(signedSSVMessage.SSVMessage.Data); err != nil {
		return i.State.Decided, i.State.DecidedValue, nil, errors.Wrap(err, "could not decode Message")
	}

	res := i.processMsgF.Run(func() interface{} {
		switch message.MsgType {
		case ProposalMsgType:
			return i.uponProposal(signedSSVMessage, i.State.ProposeContainer)
		case PrepareMsgType:
			return i.uponPrepare(signedSSVMessage, i.State.PrepareContainer)
		case CommitMsgType:
			decided, decidedValue, aggregatedCommit, err = i.UponCommit(signedSSVMessage, i.State.CommitContainer)
			if decided {
				i.State.Decided = decided
				i.State.DecidedValue = decidedValue
			}
			return err
		case RoundChangeMsgType:
			return i.uponRoundChange(i.StartValue, signedSSVMessage, i.State.RoundChangeContainer, i.config.GetValueCheckF())
		default:
			return errors.New("signed message type not supported")
		}
	})
	if res != nil {
		return false, nil, nil, res.(error)
	}
	return i.State.Decided, i.State.DecidedValue, aggregatedCommit, nil
}

func (i *Instance) BaseMsgValidation(signedSSVMessage *types.SignedSSVMessage) error {
	if err := signedSSVMessage.Validate(); err != nil {
		return errors.Wrap(err, "invalid signed message")
	}

	// Decode
	message := &Message{}
	if err := message.Decode(signedSSVMessage.SSVMessage.Data); err != nil {
		return errors.Wrap(err, "could not decode Message")
	}

	if message.Round < i.State.Round {
		return errors.New("past round")
	}

	switch message.MsgType {
	case ProposalMsgType:
		return isValidProposal(
			i.State,
			i.config,
			signedSSVMessage,
			i.config.GetValueCheckF(),
			i.State.Share.Committee,
		)
	case PrepareMsgType:
		proposedSignedSSVMsg := i.State.ProposalAcceptedForCurrentRound
		if proposedSignedSSVMsg == nil {
			return errors.New("did not receive proposal for this round")
		}
		// Decode
		proposedMessage := &Message{}
		if err := message.Decode(proposedSignedSSVMsg.SSVMessage.Data); err != nil {
			return errors.Wrap(err, "could not decode Message form ProposalAcceptedForCurrentRound")
		}

		return validSignedPrepareForHeightRoundAndRoot(
			i.config,
			signedSSVMessage,
			i.State.Height,
			i.State.Round,
			proposedMessage.Root,
			i.State.Share.Committee,
			false,
		)
	case CommitMsgType:
		proposedMsg := i.State.ProposalAcceptedForCurrentRound
		if proposedMsg == nil {
			return errors.New("did not receive proposal for this round")
		}
		return validateCommit(
			i.config,
			signedSSVMessage,
			i.State.Height,
			i.State.Round,
			i.State.ProposalAcceptedForCurrentRound,
			i.State.Share.Committee,
		)
	case RoundChangeMsgType:
		return validRoundChangeForData(i.State, i.config, signedSSVMessage, i.State.Height, message.Round, message.FullData, false)
	default:
		return errors.New("signed message type not supported")
	}
}

// IsDecided interface implementation
func (i *Instance) IsDecided() (bool, []byte) {
	if state := i.State; state != nil {
		return state.Decided, state.DecidedValue
	}
	return false, nil
}

// GetConfig returns the instance config
func (i *Instance) GetConfig() IConfig {
	return i.config
}

// GetHeight interface implementation
func (i *Instance) GetHeight() Height {
	return i.State.Height
}

// GetRoot returns the state's deterministic root
func (i *Instance) GetRoot() ([32]byte, error) {
	return i.State.GetRoot()
}

// Encode implementation
func (i *Instance) Encode() ([]byte, error) {
	return json.Marshal(i)
}

// Decode implementation
func (i *Instance) Decode(data []byte) error {
	return json.Unmarshal(data, &i)
}

// CanProcessMessages will return true if instance can process messages
func (i *Instance) CanProcessMessages() bool {
	return !i.forceStop && int(i.State.Round) < CutoffRound
}
