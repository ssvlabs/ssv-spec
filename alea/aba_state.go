package alea


type ABAState struct {
	ABAInitContainer     	*MsgContainer
	ABAAuxContainer     	*MsgContainer
	ABAConfContainer     	*MsgContainer
	ABAFinishContainer     	*MsgContainer

	Init1Counter				uint64
	Init0Counter				uint64
	Aux1Counter					uint64
	Aux0Counter					uint64
	ConfCounter					uint64
	Finish1Counter				uint64
	Finish0Counter				uint64

	SentFinish1					bool
	SentFinish0					bool

	ABARound				Round

	Vin						byte
	Vdecided				byte

	Round					Round
	Values					map[Round][]byte

	Terminate				chan bool
}

func NewABAState(abaRound Round) *ABAState {
	return &ABAState{
		ABAInitContainer:		NewMsgContainer(),
		ABAAuxContainer:		NewMsgContainer(),
		ABAConfContainer:		NewMsgContainer(),
		ABAFinishContainer:		NewMsgContainer(),
		Init1Counter:			0,
		Init0Counter:			0,
		Aux1Counter:			0,
		Aux0Counter:			0,
		ConfCounter:			0,
		Finish1Counter:			0,
		Finish0Counter:			0,
		SentFinish1:			false,
		SentFinish0:			false,
		Vin:					byte(2),
		Vdecided:				byte(2),
		Round:					0,
		Values:					make(map[Round][]byte),
	}
}


func (s *ABAState) Coin(round Round) byte {
	// FIX ME : implement a random generator given the round number
	return byte(round%2)
}


func (s *ABAState) IncrementRound() {
	s.Round += 1
	s.ABAInitContainer.Clear()
	s.ABAAuxContainer.Clear()
	s.ABAConfContainer.Clear()
	s.ABAFinishContainer.Clear()
	s.Init1Counter = 0
	s.Init0Counter = 0
	s.Aux1Counter = 0
	s.Aux0Counter = 0
	s.ConfCounter = 0
	s.Finish1Counter = 0
	s.Finish0Counter = 0
	s.SentFinish1 = false
	s.SentFinish0 = false
	s.Values = make(map[Round][]byte)
}
