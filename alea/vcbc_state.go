package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

type VCBCState struct {
	Priority Priority
	Queues   map[types.OperatorID]*VCBCQueue

	// used to store info about READY messages received
	VCBCr map[types.OperatorID]map[Priority]uint64
	VCBCW map[types.OperatorID]map[Priority][]*SignedMessage
	// used to store data from SEND messages received
	VCBCm map[types.OperatorID]map[Priority][]*ProposalData
	// used to store receipt proofs
	VCBCu map[types.OperatorID]map[Priority][]byte

	// store already received ready messages
	ReceivedReady map[types.OperatorID]map[Priority]map[types.OperatorID]bool
}

func NewVCBCState() *VCBCState {
	return &VCBCState{
		Priority:      FirstPriority,
		Queues:        make(map[types.OperatorID]*VCBCQueue),
		VCBCr:         make(map[types.OperatorID]map[Priority]uint64),
		VCBCW:         make(map[types.OperatorID]map[Priority][]*SignedMessage),
		VCBCm:         make(map[types.OperatorID]map[Priority][]*ProposalData),
		VCBCu:         make(map[types.OperatorID]map[Priority][]byte),
		ReceivedReady: make(map[types.OperatorID]map[Priority]map[types.OperatorID]bool),
	}
}

func (s *VCBCState) GetR(operatorID types.OperatorID, priority Priority) uint64 {
	if _, exists := s.VCBCr[operatorID]; !exists {
		s.VCBCr[operatorID] = make(map[Priority]uint64)
	}
	if _, exists := s.VCBCr[operatorID][priority]; !exists {
		s.VCBCr[operatorID][priority] = 0
	}
	return s.VCBCr[operatorID][priority]
}
func (s *VCBCState) IncrementR(operatorID types.OperatorID, priority Priority) {
	if _, exists := s.VCBCr[operatorID]; !exists {
		s.VCBCr[operatorID] = make(map[Priority]uint64)
	}
	if _, exists := s.VCBCr[operatorID][priority]; !exists {
		s.VCBCr[operatorID][priority] = 0
	}
	s.VCBCr[operatorID][priority] += 1
}

func (s *VCBCState) GetW(operatorID types.OperatorID, priority Priority) []*SignedMessage {
	if _, exists := s.VCBCW[operatorID]; !exists {
		s.VCBCW[operatorID] = make(map[Priority][]*SignedMessage)
	}
	if _, exists := s.VCBCW[operatorID][priority]; !exists {
		s.VCBCW[operatorID][priority] = make([]*SignedMessage, 0)
	}
	return s.VCBCW[operatorID][priority]
}

func (s *VCBCState) AppendToW(operatorID types.OperatorID, priority Priority, signedMessage *SignedMessage) {
	if _, exists := s.VCBCW[operatorID]; !exists {
		s.VCBCW[operatorID] = make(map[Priority][]*SignedMessage)
	}
	if _, exists := s.VCBCW[operatorID][priority]; !exists {
		s.VCBCW[operatorID][priority] = make([]*SignedMessage, 0)
	}
	s.VCBCW[operatorID][priority] = append(s.VCBCW[operatorID][priority], signedMessage)
}

func (s *VCBCState) HasM(operatorID types.OperatorID, priority Priority) bool {
	if _, exists := s.VCBCm[operatorID]; exists {
		if _, exists := s.VCBCm[operatorID][priority]; exists {
			return true
		}
	}
	return false
}
func (s *VCBCState) GetM(operatorID types.OperatorID, priority Priority) []*ProposalData {
	if _, exists := s.VCBCm[operatorID]; !exists {
		s.VCBCm[operatorID] = make(map[Priority][]*ProposalData)
	}
	if _, exists := s.VCBCm[operatorID][priority]; !exists {
		s.VCBCm[operatorID][priority] = make([]*ProposalData, 0)
	}
	return s.VCBCm[operatorID][priority]
}

func (s *VCBCState) AppendToM(operatorID types.OperatorID, priority Priority, proposal *ProposalData) {
	if _, exists := s.VCBCm[operatorID]; !exists {
		s.VCBCm[operatorID] = make(map[Priority][]*ProposalData)
	}
	if _, exists := s.VCBCm[operatorID][priority]; !exists {
		s.VCBCm[operatorID][priority] = make([]*ProposalData, 0)
	}
	s.VCBCm[operatorID][priority] = append(s.VCBCm[operatorID][priority], proposal)
}

func (s *VCBCState) SetM(operatorID types.OperatorID, priority Priority, proposals []*ProposalData) {
	if _, exists := s.VCBCm[operatorID]; !exists {
		s.VCBCm[operatorID] = make(map[Priority][]*ProposalData)
	}
	s.VCBCm[operatorID][priority] = proposals
}

func (s *VCBCState) EqualM(operatorID types.OperatorID, priority Priority, proposals []*ProposalData) bool {
	if !s.HasM(operatorID, priority) {
		return false
	}

	if len(s.VCBCm[operatorID][priority]) != len(proposals) {
		return false
	}
	for idx, proposal := range s.VCBCm[operatorID][priority] {
		if !proposal.Equal(proposals[idx]) {
			return false
		}
	}
	return true
}

func (s *VCBCState) hasU(operatorID types.OperatorID, priority Priority) bool {
	if _, exists := s.VCBCu[operatorID]; exists {
		if _, exists := s.VCBCu[operatorID][priority]; exists {
			return true
		}
	}
	return false
}
func (s *VCBCState) GetU(operatorID types.OperatorID, priority Priority) []byte {
	if _, exists := s.VCBCu[operatorID]; !exists {
		s.VCBCu[operatorID] = make(map[Priority][]byte)
	}
	if _, exists := s.VCBCu[operatorID][priority]; !exists {
		s.VCBCu[operatorID][priority] = nil
	}
	return s.VCBCu[operatorID][priority]
}
func (s *VCBCState) SetU(operatorID types.OperatorID, priority Priority, u []byte) {
	if _, exists := s.VCBCu[operatorID]; !exists {
		s.VCBCu[operatorID] = make(map[Priority][]byte)
	}
	s.VCBCu[operatorID][priority] = u
}

func (s *VCBCState) HasReceivedReady(author types.OperatorID, priority Priority, operatorID types.OperatorID) bool {
	if _, exists := s.ReceivedReady[author]; !exists {
		s.ReceivedReady[author] = make(map[Priority]map[types.OperatorID]bool)
	}
	if _, exists := s.ReceivedReady[author][priority]; !exists {
		s.ReceivedReady[author][priority] = make(map[types.OperatorID]bool)
	}
	if _, exists := s.ReceivedReady[author][priority][operatorID]; !exists {
		s.ReceivedReady[author][priority][operatorID] = false
	}
	return s.ReceivedReady[author][priority][operatorID]
}

func (s *VCBCState) SetReceivedReady(author types.OperatorID, priority Priority, operatorID types.OperatorID, received bool) {
	if _, exists := s.ReceivedReady[author]; !exists {
		s.ReceivedReady[author] = make(map[Priority]map[types.OperatorID]bool)
	}
	if _, exists := s.ReceivedReady[author][priority]; !exists {
		s.ReceivedReady[author][priority] = make(map[types.OperatorID]bool)
	}
	s.ReceivedReady[author][priority][operatorID] = received
}
