package qbft

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/types"
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
	committeeMember *types.CommitteeMember,
	identifier []byte,
	height Height,
) *Instance {
	return &Instance{
		State: &State{
			CommitteeMember:      committeeMember,
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
		if proposer(i.State, i.GetConfig(), FirstRound) == i.State.CommitteeMember.OperatorID {
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

func (i *Instance) Broadcast(msg *types.SignedSSVMessage) error {
	if !i.CanProcessMessages() {
		return errors.New("instance stopped processing messages")
	}

	return i.GetConfig().GetNetwork().Broadcast(msg.SSVMessage.GetID(), msg)
}

// ProcessMsg processes a new QBFT msg, returns non nil error on msg processing error
func (i *Instance) ProcessMsg(msg *ProcessingMessage) (decided bool, decidedValue []byte, aggregatedCommit *types.SignedSSVMessage, err error) {
	if !i.CanProcessMessages() {
		return false, nil, nil, errors.New("instance stopped processing messages")
	}

	if err := i.BaseMsgValidation(msg); err != nil {
		return false, nil, nil, errors.Wrap(err, "invalid signed message")
	}

	res := i.processMsgF.Run(func() interface{} {
		switch msg.QBFTMessage.MsgType {
		case ProposalMsgType:
			return i.uponProposal(msg, i.State.ProposeContainer)
		case PrepareMsgType:
			return i.uponPrepare(msg, i.State.PrepareContainer)
		case CommitMsgType:
			decided, decidedValue, aggregatedCommit, err = i.UponCommit(msg, i.State.CommitContainer)
			if decided {
				i.State.Decided = decided
				i.State.DecidedValue = decidedValue
			}
			return err
		case RoundChangeMsgType:
			return i.uponRoundChange(i.StartValue, msg, i.State.RoundChangeContainer, i.config.GetValueCheckF())
		default:
			return errors.New("signed message type not supported")
		}
	})
	if res != nil {
		return false, nil, nil, res.(error)
	}
	return i.State.Decided, i.State.DecidedValue, aggregatedCommit, nil
}

func (i *Instance) BaseMsgValidation(msg *ProcessingMessage) error {
	if err := msg.Validate(); err != nil {
		return err
	}

	if msg.QBFTMessage.Round < i.State.Round {
		return errors.New("past round")
	}

	switch msg.QBFTMessage.MsgType {
	case ProposalMsgType:
		return isValidProposal(
			i.State,
			i.config,
			msg,
			i.config.GetValueCheckF(),
		)
	case PrepareMsgType:
		proposedMsg := i.State.ProposalAcceptedForCurrentRound
		if proposedMsg == nil {
			return errors.New("did not receive proposal for this round")
		}
		return validSignedPrepareForHeightRoundAndRootIgnoreSignature(
			msg,
			i.State.Height,
			i.State.Round,
			proposedMsg.QBFTMessage.Root,
			i.State.CommitteeMember.Committee,
		)
	case CommitMsgType:
		proposedMsg := i.State.ProposalAcceptedForCurrentRound
		if proposedMsg == nil {
			return errors.New("did not receive proposal for this round")
		}
		return validateCommit(
			msg,
			i.State.Height,
			i.State.Round,
			i.State.ProposalAcceptedForCurrentRound,
			i.State.CommitteeMember.Committee,
		)
	case RoundChangeMsgType:
		return validRoundChangeForDataIgnoreSignature(i.State, i.config, msg, i.State.Height, msg.QBFTMessage.Round, msg.SignedMessage.FullData)
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

// CanProcessMessages will return true if instance can process messages
func (i *Instance) CanProcessMessages() bool {
	return !i.forceStop && i.State.Round < i.config.GetCutOffRound()
}
