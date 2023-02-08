package alea

type ACState struct {
	ACRound  ACRound
	ABAState map[ACRound]*ABAState
}

func NewACState() *ACState {
	acState := &ACState{
		ACRound:  FirstACRound,
		ABAState: make(map[ACRound]*ABAState),
	}
	acState.ABAState[acState.ACRound] = NewABAState(acState.ACRound)
	return acState
}

func (s *ACState) IncrementRound() {
	// update info
	s.ACRound += 1
	s.InitializeRound(s.ACRound)
}

func (s *ACState) InitializeRound(acRound ACRound) {
	if _, exists := s.ABAState[acRound]; !exists {
		s.ABAState[acRound] = NewABAState(acRound)
	}
}

func (s *ACState) GetCurrentABAState() *ABAState {
	return s.ABAState[s.ACRound]
}

func (s *ACState) GetABAState(acRound ACRound) *ABAState {
	return s.ABAState[acRound]
}
