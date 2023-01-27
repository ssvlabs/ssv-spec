package alea

type ABAState struct {
	// message containers
	ABAInitContainer   *MsgContainer
	ABAAuxContainer    *MsgContainer
	ABAConfContainer   *MsgContainer
	ABAFinishContainer *MsgContainer
	// message counters
	InitCounter   []uint64
	AuxCounter    []uint64
	ConfCounter   uint64
	FinishCounter []uint64
	// already sent message flags
	SentInit   []bool
	SentAux    []bool
	SentConf   bool
	SentFinish []bool
	// current ABA round
	ACRound Round
	// value inputed to ABA
	Vin byte
	// value decided by ABA
	Vdecided byte
	// current ABA round
	Round Round
	// values that completed strong support of INIT messages
	Values map[Round][]byte
	// terminate channel to announce to ABA caller
	Terminate bool
}

func NewABAState(acRound Round) *ABAState {
	return &ABAState{
		ABAInitContainer:   NewMsgContainer(),
		ABAAuxContainer:    NewMsgContainer(),
		ABAConfContainer:   NewMsgContainer(),
		ABAFinishContainer: NewMsgContainer(),
		InitCounter:        make([]uint64, 2),
		AuxCounter:         make([]uint64, 2),
		ConfCounter:        0,
		FinishCounter:      make([]uint64, 2),
		SentInit:           make([]bool, 2),
		SentAux:            make([]bool, 2),
		SentConf:           false,
		SentFinish:         make([]bool, 2),
		ACRound:            acRound,
		Vin:                byte(2),
		Vdecided:           byte(2),
		Round:              FirstRound,
		Values:             make(map[Round][]byte),
	}
}

func (s *ABAState) Coin(round Round) byte {
	// FIX ME : implement a RANDOM coin generator given the round number
	return byte(round % 2)
}

func (s *ABAState) IncrementRound() {
	// update info
	s.Round += 1
	// s.ABAInitContainer.Clear()
	// s.ABAAuxContainer.Clear()
	// s.ABAConfContainer.Clear()
	// s.ABAFinishContainer.Clear()

	s.InitCounter = make([]uint64, 2)
	s.AuxCounter = make([]uint64, 2)
	s.ConfCounter = 0
	// s.FinishCounter = make([]uint64, 2)
	s.SentInit = make([]bool, 2)
	s.SentAux = make([]bool, 2)
	s.SentConf = false
	// s.SentFinish = make([]bool, 2)

	s.Values = make(map[Round][]byte)
}
