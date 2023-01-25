package alea


type ABAState struct {
	// message containers
	ABAInitContainer     		*MsgContainer
	ABAAuxContainer     		*MsgContainer
	ABAConfContainer     		*MsgContainer
	ABAFinishContainer     		*MsgContainer
	// message counters
	Init1Counter				uint64
	Init0Counter				uint64
	Aux1Counter					uint64
	Aux0Counter					uint64
	ConfCounter					uint64
	Finish1Counter				uint64
	Finish0Counter				uint64
	// already sent message flags
	SentInit1					bool
	SentInit0					bool
	SentAux1					bool
	SentAux0					bool
	SentConf					bool
	SentFinish1					bool
	SentFinish0					bool
	// current ABA round
	ACRound						Round
	// value inputed to ABA
	Vin							byte
	// value decided by ABA
	Vdecided					byte
	// current ABA round
	Round						Round
	// values that completed strong support of INIT messages
	Values						map[Round][]byte
	// terminate channel to announce to ABA caller
	Terminate					chan bool
}

func NewABAState(acRound Round) *ABAState {
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
		SentInit1:				false,
		SentInit0:				false,
		SentAux1:				false,
		SentAux0:				false,
		SentConf:				false,
		SentFinish1:			false,
		SentFinish0:			false,
		ACRound:				acRound,
		Vin:					byte(2),
		Vdecided:				byte(2),
		Round:					0,
		Values:					make(map[Round][]byte),
	}
}


func (s *ABAState) Coin(round Round) byte {
	// FIX ME : implement a RANDOM coin generator given the round number
	return byte(round%2)
}


func (s *ABAState) IncrementRound() {
	// update info
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
	s.SentInit1 = false
	s.SentInit0 = false
	s.SentAux1 = false
	s.SentAux0 = false
	s.SentConf = false
	s.SentFinish1 = false
	s.SentFinish0 = false
	s.Values = make(map[Round][]byte)
}
