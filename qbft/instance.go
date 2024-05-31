package qbft

import (
	"encoding/json"
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
	share *types.SharedValidator,
	identifier []byte,
	height Height,
) *Instance {
	return &Instance{
		State: &State{
			SharedValidator:      share,
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
		if proposer(i.State, i.GetConfig(), FirstRound) == i.State.SharedValidator.OwnValidatorShare.OperatorID {
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
func (i *Instance) ProcessMsg(signedMsg *types.SignedSSVMessage) (decided bool, decidedValue []byte, aggregatedCommit *types.SignedSSVMessage, err error) {
	if !i.CanProcessMessages() {
		return false, nil, nil, errors.New("instance stopped processing messages")
	}

	msg, err := DecodeMessage(signedMsg.SSVMessage.Data)
	if err != nil {
		return false, nil, nil, err
	}

	if err := i.BaseMsgValidation(signedMsg); err != nil {
		return false, nil, nil, errors.Wrap(err, "invalid signed message")
	}

	res := i.processMsgF.Run(func() interface{} {
		switch msg.MsgType {
		case ProposalMsgType:
			return i.uponProposal(signedMsg, i.State.ProposeContainer)
		case PrepareMsgType:
			return i.uponPrepare(signedMsg, i.State.PrepareContainer)
		case CommitMsgType:
			decided, decidedValue, aggregatedCommit, err = i.UponCommit(signedMsg, i.State.CommitContainer)
			if decided {
				i.State.Decided = decided
				i.State.DecidedValue = decidedValue
			}
			return err
		case RoundChangeMsgType:
			return i.uponRoundChange(i.StartValue, signedMsg, i.State.RoundChangeContainer, i.config.GetValueCheckF())
		default:
			return errors.New("signed message type not supported")
		}
	})
	if res != nil {
		return false, nil, nil, res.(error)
	}
	return i.State.Decided, i.State.DecidedValue, aggregatedCommit, nil
}

func (i *Instance) BaseMsgValidation(signedMsg *types.SignedSSVMessage) error {
	if err := signedMsg.Validate(); err != nil {
		return errors.Wrap(err, "invalid SignedSSVMessage")
	}

	msg, err := DecodeMessage(signedMsg.SSVMessage.Data)
	if err != nil {
		return err
	}

	if err := msg.Validate(); err != nil {
		return errors.Wrap(err, "invalid Message")
	}

	if msg.Round < i.State.Round {
		return errors.New("past round")
	}

	switch msg.MsgType {
	case ProposalMsgType:
		return isValidProposal(
			i.State,
			i.config,
			signedMsg,
			i.config.GetValueCheckF(),
		)
	case PrepareMsgType:
		proposedSignedMsg := i.State.ProposalAcceptedForCurrentRound
		if proposedSignedMsg == nil {
			return errors.New("did not receive proposal for this round")
		}

		proposedMsg, err := DecodeMessage(proposedSignedMsg.SSVMessage.Data)
		if err != nil {
			return errors.Wrap(err, "proposal saved for this round is invalid")
		}

		return validSignedPrepareForHeightRoundAndRootIgnoreSignature(
			signedMsg,
			i.State.Height,
			i.State.Round,
			proposedMsg.Root,
			i.State.SharedValidator.Committee,
		)
	case CommitMsgType:
		proposedMsg := i.State.ProposalAcceptedForCurrentRound
		if proposedMsg == nil {
			return errors.New("did not receive proposal for this round")
		}
		return validateCommit(
			signedMsg,
			i.State.Height,
			i.State.Round,
			i.State.ProposalAcceptedForCurrentRound,
			i.State.SharedValidator.Committee,
		)
	case RoundChangeMsgType:
		return validRoundChangeForDataIgnoreSignature(i.State, i.config, signedMsg, i.State.Height, msg.Round, signedMsg.FullData)
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
	return !i.forceStop && i.State.Round < i.config.GetCutOffRound()
}
