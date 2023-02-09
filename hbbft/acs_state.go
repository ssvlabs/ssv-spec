package hbbft

import "github.com/MatheusFranco99/ssv-spec-AleaBFT/types"

type ACSState struct {
	ABAState *ABAState
}

func NewACSState(acsRound ACSRound) *ACSState {
	acsState := &ACSState{
		ABAState: NewABAState(acsRound),
	}
	return acsState
}

func (s *ACSState) GetABAState() *ABAState {
	return s.ABAState
}

func (s *ACSState) StartACS(proposed_encrypted []byte) map[types.OperatorID][]byte {
	return make(map[types.OperatorID][]byte)
}
