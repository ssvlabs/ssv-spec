package validation

import (
	"context"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/ssvlabs/ssv-spec/types"
)

// Message validator interface for the libp2p pubsub framework
type IMessageValidator interface {
	ValidatorForTopic(topic string) func(ctx context.Context, peerID peer.ID, msg *pubsub.Message) pubsub.ValidationResult
}

type MessageValidator struct {
	PeerState   map[peer.ID]*PeerState
	DutyCounter map[peer.ID]*DutyCounter
	Beacon      types.BeaconNetwork
	Network     NetworkDataFetcher
	Config      Config
	DutyFetcher DutyFetcher
}

func NewMessageValidator(beacon types.BeaconNetwork, network NetworkDataFetcher, dutyFetcher DutyFetcher, config Config) *MessageValidator {
	return &MessageValidator{
		Beacon:      beacon,
		Network:     network,
		Config:      config,
		DutyFetcher: dutyFetcher,
		PeerState:   make(map[peer.ID]*PeerState),
		DutyCounter: make(map[peer.ID]*DutyCounter),
	}
}

// Used to validate messages within the libp2p pubsub framework
func (mv *MessageValidator) ValidatorForTopic(_ string) func(ctx context.Context, p peer.ID, pmsg *pubsub.Message) pubsub.ValidationResult {
	return mv.Validate
}

// Validates a pubsub message
func (mv *MessageValidator) Validate(_ context.Context, peerID peer.ID, pmsg *pubsub.Message) pubsub.ValidationResult {

	// Validate
	err := mv.HandlePubsubMessage(pmsg, time.Now())
	if err != nil {
		return handleValidationError(err)
	}

	// Update state
	if err := mv.UpdateState(peerID, pmsg); err != nil {
		return handleValidationError(err)
	}

	return pubsub.ValidationAccept
}

// Handles a pubsub message and performs the full validation chain
func (mv *MessageValidator) HandlePubsubMessage(pMsg *pubsub.Message, receivedAt time.Time) error {

	// Validate pubsub message
	if err := mv.validatePubSubMessage(pMsg); err != nil {
		return err
	}

	msg, err := decodeSignedSSVMessage(pMsg)
	if err != nil {
		return err
	}

	// Validate syntax
	if err := mv.validateSyntax(msg); err != nil {
		return err
	}

	// Validate general semantics
	if err := mv.validateGeneralSemantics(msg, *pMsg.Topic); err != nil {
		return err
	}

	switch msg.SSVMessage.MsgType {
	case types.SSVConsensusMsgType:
		// Validate consensus message for: semantics, qbft logic and duty logic
		if err := mv.ValidateConsensusMessageSemantics(pMsg.ReceivedFrom, msg); err != nil {
			return err
		}
		if err := mv.ValidateConsensusMessageQBFTLogic(pMsg.ReceivedFrom, msg, receivedAt); err != nil {
			return err
		}
		if err := mv.ValidateMessageDutyLogic(pMsg.ReceivedFrom, msg, receivedAt); err != nil {
			return err
		}
		if err := mv.ValidateConsensusMessageDutyLogic(pMsg.ReceivedFrom, msg, receivedAt); err != nil {
			return err
		}
	case types.SSVPartialSignatureMsgType:
		// Validate partial signature message for: semantics and duty logic
		if err := mv.ValidatePartialSignatureMessageSemantics(pMsg.ReceivedFrom, msg); err != nil {
			return err
		}
		if err := mv.ValidateMessageDutyLogic(pMsg.ReceivedFrom, msg, receivedAt); err != nil {
			return err
		}
		if err := mv.ValidatePartialSignatureMessageDutyLogic(pMsg.ReceivedFrom, msg, receivedAt); err != nil {
			return err
		}
	default:
		return ErrUnknownSSVMessageType
	}

	// Verify signature
	if err := types.Verify(msg, mv.Network.GetCommitteeInfo(msg.SSVMessage.MsgID).Operators); err != nil {
		return ErrSignatureVerification
	}

	return nil
}

func (mv *MessageValidator) GetDutyCounter(peerID peer.ID) *DutyCounter {
	if _, exists := mv.DutyCounter[peerID]; !exists {
		mv.DutyCounter[peerID] = NewDutyCounter()
	}
	return mv.DutyCounter[peerID]
}

func (mv *MessageValidator) GetPeerState(peerID peer.ID) *PeerState {
	if _, exists := mv.PeerState[peerID]; !exists {
		mv.PeerState[peerID] = NewPeerState()
	}
	return mv.PeerState[peerID]
}
