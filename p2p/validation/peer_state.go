package validation

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

// PeerState represents the state of a peer's view on the Network.
// It's useful to validate upcoming messages from this specific peer.
// For each possible duty (identified by a MessageID (validator/committee and role) and a slot), it holds a DutyState
type PeerState struct {
	dutyState map[types.MessageID]map[phase0.Slot]*DutyState
}

func NewPeerState() *PeerState {
	return &PeerState{
		dutyState: make(map[types.MessageID]map[phase0.Slot]*DutyState),
	}
}

// Returns the duty state given an identifier and slot
func (ps *PeerState) GetDutyState(msgID types.MessageID, slot phase0.Slot) *DutyState {
	if _, exists := ps.dutyState[msgID]; !exists {
		ps.dutyState[msgID] = make(map[phase0.Slot]*DutyState)
	}
	if _, exists := ps.dutyState[msgID][slot]; !exists {
		ps.dutyState[msgID][slot] = NewDutyState()
	}
	return ps.dutyState[msgID][slot]
}

// For a given identifier and signer, returns the highest slot in which the signer participates
func (ps *PeerState) GetHighestSlotForSigner(msgID types.MessageID, signer types.OperatorID) phase0.Slot {
	highestSlot := phase0.Slot(0)
	if _, exists := ps.dutyState[msgID]; !exists {
		return highestSlot
	}
	for slot, dutyState := range ps.dutyState[msgID] {
		if slot > highestSlot {
			if _, exists := dutyState.signerState[signer]; exists {
				highestSlot = slot
			}
		}
	}
	return highestSlot
}

// DutyState manages a SignerState for each
// operator that belongs to the duty and
// a set of decided messages (by recording the signers)
type DutyState struct {
	signerState map[types.OperatorID]*SignerState
	decideds    map[[sha256.Size]byte]struct{}
}

func NewDutyState() *DutyState {
	return &DutyState{
		signerState: make(map[types.OperatorID]*SignerState),
		decideds:    make(map[[sha256.Size]byte]struct{}),
	}
}

// Returns the SignerState for a given operator
func (ds *DutyState) GetSignerState(operator types.OperatorID) *SignerState {
	if _, exists := ds.signerState[operator]; !exists {
		ds.signerState[operator] = NewSignerState()
	}
	return ds.signerState[operator]
}

// Record decided message
func (ds *DutyState) RecordDecided(signers []types.OperatorID) {
	encodedOperators := encodeOperators(signers)
	ds.decideds[encodedOperators] = struct{}{}
}

// Check if it has a decided message
func (ds *DutyState) HasDecided(signers []types.OperatorID) bool {
	encodedOperators := encodeOperators(signers)
	_, exists := ds.decideds[encodedOperators]
	return exists
}

// SignerState is the state of a specific operator regarding a duty.
// It informs the current consensus round of the operator,
// the number of messages per type and
// the proposal data for the current consensus round.
// The proposal data is useful to detect multiple proposals with different data
type SignerState struct {
	Round         qbft.Round
	MessageCounts MessageCounts
	ProposalData  []byte
}

func NewSignerState() *SignerState {
	return &SignerState{
		Round:         qbft.FirstRound,
		MessageCounts: MessageCounts{},
	}
}

// Process a consensus message by updating the message counters,
// and, if necessary, bumping the round and storing the proposal data
func (ss *SignerState) ProcessConsensusMessage(msg *qbft.Message, fullData []byte, signers int) {
	if msg.Round > ss.Round {
		ss.BumpToRound(msg.Round)
	}
	if msg.MsgType == qbft.ProposalMsgType {
		ss.ProposalData = fullData
	}
	ss.MessageCounts.RecordConsensusMessage(msg.MsgType, signers)
}

// Process a partial signature message by updating the message counters
func (ss *SignerState) ProcessPartialSignatureMessage(msg *types.PartialSignatureMessages) {
	ss.MessageCounts.RecordPartialSignatureMessage(msg.Type)
}

// Bumps the state to a given round erasing the proposal data and
// reseting the message counters
func (ss *SignerState) BumpToRound(round qbft.Round) {
	ss.Round = round
	ss.MessageCounts.ResetForRound()
	ss.ProposalData = make([]byte, 0)
}

// Auxiliary method to transform a list of OperatorID to a [32]byte
// so that it can be used as key in a map
func encodeOperators(operators []types.OperatorID) [sha256.Size]byte {
	// Convert to bytes
	committeeBytes := make([]byte, len(operators)*4)
	for i, v := range operators {
		binary.LittleEndian.PutUint32(committeeBytes[i*4:], uint32(v))
	}
	// Hash
	hash := sha256.Sum256(committeeBytes)
	return hash
}

// After a successfull validation, updates the state considering the peer's message by
// updating the PeerState and DutyCounter state
func (mv *MessageValidator) UpdateState(peerID peer.ID, pmsg *pubsub.Message) error {

	msg, err := decodeSignedSSVMessage(pmsg)
	if err != nil {
		return err
	}

	// Update the PeerState according to the message type
	if msg.SSVMessage.MsgType == types.SSVConsensusMsgType {
		if err := mv.updatePeerStateForConsensusMessage(peerID, msg); err != nil {
			return err
		}
	} else if msg.SSVMessage.MsgType == types.SSVPartialSignatureMsgType {
		if err := mv.updatePeerStateForPartialSignatureMessage(peerID, msg); err != nil {
			return err
		}
	} else {
		return ErrUnknownSSVMessageType
	}

	// Updates duty counter
	slot, err := getMessageSlot(msg)
	if err != nil {
		return err
	}
	mv.GetDutyCounter(peerID).RecordDuty(msg.SSVMessage.MsgID, mv.Beacon.EstimatedEpochAtSlot(slot), slot)

	return nil
}

// Updates the PeerState for a consensus message
func (mv *MessageValidator) updatePeerStateForConsensusMessage(peerID peer.ID, msg *types.SignedSSVMessage) error {

	// Decode to qbft.Message
	qbftMsg := &qbft.Message{}
	err := qbftMsg.Decode(msg.SSVMessage.Data)
	if err != nil {
		return ErrUndecodableData
	}

	dutyState := mv.GetPeerState(peerID).GetDutyState(msg.SSVMessage.MsgID, phase0.Slot(qbftMsg.Height))

	if qbftMsg.MsgType == qbft.CommitMsgType && len(msg.OperatorIDs) > 1 {
		// If decided message, record decided signers
		dutyState.RecordDecided(msg.OperatorIDs)
	} else {
		// Update signer state
		signer := msg.OperatorIDs[0]
		signerState := dutyState.GetSignerState(signer)
		signerState.ProcessConsensusMessage(qbftMsg, msg.FullData, len(msg.OperatorIDs))
	}

	return nil
}

// Updates the PeerState for a partial signature message
func (mv *MessageValidator) updatePeerStateForPartialSignatureMessage(peerID peer.ID, msg *types.SignedSSVMessage) error {

	// Decode to partial signature message
	pSigMsgs := &types.PartialSignatureMessages{}
	err := pSigMsgs.Decode(msg.SSVMessage.Data)
	if err != nil {
		return ErrUndecodableData
	}

	// Update signer state
	dutyState := mv.GetPeerState(peerID).GetDutyState(msg.SSVMessage.MsgID, pSigMsgs.Slot)
	signer := msg.OperatorIDs[0]
	signerState := dutyState.GetSignerState(signer)
	signerState.ProcessPartialSignatureMessage(pSigMsgs)

	return nil
}
