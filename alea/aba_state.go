package alea

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

type ABAState struct {
	// message containers
	ABAInitContainer   *MsgContainer
	ABAAuxContainer    *MsgContainer
	ABAConfContainer   *MsgContainer
	ABAFinishContainer *MsgContainer
	// message counters
	InitCounter   map[Round]map[byte][]types.OperatorID
	AuxCounter    map[Round]map[byte][]types.OperatorID
	ConfCounter   map[Round][]types.OperatorID
	ConfValues    map[Round]map[types.OperatorID][]byte
	FinishCounter map[byte][]types.OperatorID
	// already sent message flags
	sentInit   map[Round][]bool
	sentAux    map[Round][]bool
	sentConf   map[Round]bool
	sentFinish []bool
	// current ABA round
	ACRound ACRound
	// value inputed to ABA
	Vin map[Round]byte
	// value decided by ABA
	Vdecided byte
	// current ABA round
	Round Round
	// values that completed strong support of INIT messages
	Values map[Round][]byte
	// terminate channel to announce to ABA caller
	Terminate bool
}

func NewABAState(acRound ACRound) *ABAState {
	abaState := &ABAState{
		ABAInitContainer:   NewMsgContainer(),
		ABAAuxContainer:    NewMsgContainer(),
		ABAConfContainer:   NewMsgContainer(),
		ABAFinishContainer: NewMsgContainer(),
		InitCounter:        make(map[Round]map[byte][]types.OperatorID),
		AuxCounter:         make(map[Round]map[byte][]types.OperatorID),
		ConfCounter:        make(map[Round][]types.OperatorID),
		ConfValues:         make(map[Round]map[types.OperatorID][]byte),
		FinishCounter:      make(map[byte][]types.OperatorID),
		sentInit:           make(map[Round][]bool),
		sentAux:            make(map[Round][]bool),
		sentConf:           make(map[Round]bool),
		sentFinish:         make([]bool, 2),
		ACRound:            acRound,
		Vin:                make(map[Round]byte),
		Vdecided:           byte(2),
		Round:              FirstRound,
		Values:             make(map[Round][]byte),
	}

	abaState.InitializeRound(FirstRound)
	abaState.FinishCounter[0] = make([]types.OperatorID, 0)
	abaState.FinishCounter[1] = make([]types.OperatorID, 0)

	return abaState
}

func (s *ABAState) InitializeRound(round Round) {

	if _, exists := s.InitCounter[round]; !exists {
		s.InitCounter[round] = make(map[byte][]types.OperatorID)
		s.InitCounter[round][0] = make([]types.OperatorID, 0)
		s.InitCounter[round][1] = make([]types.OperatorID, 0)
	}

	if _, exists := s.AuxCounter[round]; !exists {
		s.AuxCounter[round] = make(map[byte][]types.OperatorID, 2)
		s.AuxCounter[round][0] = make([]types.OperatorID, 0)
		s.AuxCounter[round][1] = make([]types.OperatorID, 0)
	}

	if _, exists := s.ConfCounter[round]; !exists {
		s.ConfCounter[round] = make([]types.OperatorID, 0)
	}
	if _, exists := s.ConfValues[round]; !exists {
		s.ConfValues[round] = make(map[types.OperatorID][]byte)
	}

	if len(s.sentInit) == 0 {
		s.sentInit = make(map[Round][]bool)
	}
	if _, exists := s.sentInit[round]; !exists {
		s.sentInit[round] = make([]bool, 2)
	}

	if len(s.sentAux) == 0 {
		s.sentAux = make(map[Round][]bool)
	}
	if _, exists := s.sentAux[round]; !exists {
		s.sentAux[round] = make([]bool, 2)
	}
	if len(s.sentConf) == 0 {
		s.sentConf = make(map[Round]bool)
	}
	if _, exists := s.sentConf[round]; !exists {
		s.sentConf[round] = false
	}

	if _, exists := s.Values[round]; !exists {
		s.Values[round] = make([]byte, 0)
	}
}

func (s *ABAState) IncrementRound() {
	// update info
	s.Round += 1
	s.InitializeRound(s.Round)
}

func (s *ABAState) HasInit(round Round, operatorID types.OperatorID, vote byte) bool {
	for _, opID := range s.InitCounter[round][vote] {
		if opID == operatorID {
			return true
		}
	}
	return false
}
func (s *ABAState) HasAux(round Round, operatorID types.OperatorID, vote byte) bool {
	for _, opID := range s.AuxCounter[round][vote] {
		if opID == operatorID {
			return true
		}
	}
	return false
}
func (s *ABAState) HasConf(round Round, operatorID types.OperatorID) bool {
	for _, opID := range s.ConfCounter[round] {
		if opID == operatorID {
			return true
		}
	}
	return false
}
func (s *ABAState) HasFinish(operatorID types.OperatorID) bool {
	for _, vote := range []byte{0, 1} {
		for _, opID := range s.FinishCounter[vote] {
			if opID == operatorID {
				return true
			}
		}
	}
	return false
}

func (s *ABAState) CountInit(round Round, vote byte) uint64 {
	return uint64(len(s.InitCounter[round][vote]))
}
func (s *ABAState) CountAux(round Round, vote byte) uint64 {
	return uint64(len(s.AuxCounter[round][vote]))
}
func (s *ABAState) CountConf(round Round) uint64 {
	return uint64(len(s.ConfCounter[round]))
}
func (s *ABAState) CountFinish(vote byte) uint64 {
	return uint64(len(s.FinishCounter[vote]))
}

func (s *ABAState) SetInit(round Round, operatorID types.OperatorID, vote byte) {
	s.InitCounter[round][vote] = append(s.InitCounter[round][vote], operatorID)
}
func (s *ABAState) SetAux(round Round, operatorID types.OperatorID, vote byte) {
	s.AuxCounter[round][vote] = append(s.AuxCounter[round][vote], operatorID)
}
func (s *ABAState) SetConf(round Round, operatorID types.OperatorID, votes []byte) {
	s.ConfCounter[round] = append(s.ConfCounter[round], operatorID)

	if _, exists := s.ConfValues[round]; !exists {
		s.ConfValues[round] = make(map[types.OperatorID][]byte)
	}
	s.ConfValues[round][operatorID] = make([]byte, 0)
	for _, vote := range votes {
		s.ConfValues[round][operatorID] = append(s.ConfValues[round][operatorID], vote)
	}
}
func (s *ABAState) SetFinish(operatorID types.OperatorID, vote byte) {
	s.FinishCounter[vote] = append(s.FinishCounter[vote], operatorID)
}

func (s *ABAState) SentInit(round Round, vote byte) bool {
	if len(s.sentInit) == 0 {
		s.sentInit = make(map[Round][]bool)
	}
	if _, exists := s.sentInit[round]; !exists {
		s.sentInit[round] = make([]bool, 2)
	}
	return s.sentInit[round][vote]
}
func (s *ABAState) SentAux(round Round, vote byte) bool {
	if len(s.sentAux) == 0 {
		s.sentAux = make(map[Round][]bool)
	}
	if _, exists := s.sentAux[round]; !exists {
		s.sentAux[round] = make([]bool, 2)
	}
	return s.sentAux[round][vote]
}
func (s *ABAState) SentConf(round Round) bool {
	if len(s.sentConf) == 0 {
		s.sentConf = make(map[Round]bool)
	}
	if _, exists := s.sentConf[round]; !exists {
		s.sentConf[round] = false
	}
	return s.sentConf[round]
}
func (s *ABAState) SentFinish(vote byte) bool {
	if len(s.sentFinish) == 0 {
		s.sentFinish = make([]bool, 2)
	}
	return s.sentFinish[vote]
}

func (s *ABAState) SetSentInit(round Round, vote byte, value bool) {
	if len(s.sentInit) == 0 {
		s.sentInit = make(map[Round][]bool)
	}
	if _, exists := s.sentInit[round]; !exists {
		s.sentInit[round] = make([]bool, 2)
	}
	s.sentInit[round][vote] = value
}
func (s *ABAState) SetSentAux(round Round, vote byte, value bool) {
	if len(s.sentAux) == 0 {
		s.sentAux = make(map[Round][]bool)
	}
	if _, exists := s.sentAux[round]; !exists {
		s.sentAux[round] = make([]bool, 2)
	}
	s.sentAux[round][vote] = value
}
func (s *ABAState) SetSentConf(round Round, value bool) {
	if len(s.sentConf) == 0 {
		s.sentConf = make(map[Round]bool)
	}
	s.sentConf[round] = value
}
func (s *ABAState) SetSentFinish(vote byte, value bool) {
	if len(s.sentConf) == 0 {
		s.sentConf = make(map[Round]bool)
	}
	s.sentFinish[vote] = value
}

func (s *ABAState) GetValues(round Round) []byte {
	return s.Values[round]
}

func (s *ABAState) AddToValues(round Round, vote byte) {
	for _, value := range s.Values[round] {
		if value == vote {
			return
		}
	}
	s.Values[round] = append(s.Values[round], vote)
}

func (s *ABAState) isContainedInValues(round Round, values []byte) bool {
	num_equal := 0
	for _, value := range values {
		for _, storedValue := range s.Values[round] {
			if value == storedValue {
				num_equal += 1
			}
		}
	}
	return num_equal == len(values)
}
func (s *ABAState) existsInValues(round Round, value byte) bool {
	for _, storedValue := range s.Values[round] {
		if value == storedValue {
			return true
		}
	}
	return false
}

func (s *ABAState) CountAuxInValues(round Round) uint64 {
	ans := uint64(0)
	if s.existsInValues(round, byte(0)) {
		ans += uint64(len(s.AuxCounter[round][byte(0)]))
	}
	if s.existsInValues(round, byte(1)) {
		ans += uint64(len(s.AuxCounter[round][byte(1)]))
	}
	return ans
}

func (s *ABAState) CountConfContainedInValues(round Round) uint64 {

	if _, exists := s.ConfCounter[round]; !exists {
		return uint64(0)
	}

	ans := uint64(0)

	for _, opID := range s.ConfCounter[round] {
		if votes, exists := s.ConfValues[round][opID]; exists {
			if s.isContainedInValues(round, votes) {
				ans += uint64(1)
			}
		}
	}
	return ans
}

func (s *ABAState) SetVInput(round Round, vote byte) {
	s.Vin[round] = vote
}
func (s *ABAState) GetVInput(round Round) byte {
	return s.Vin[round]
}

func (s *ABAState) SetDecided(vote byte) {
	s.Vdecided = vote
}

func (s *ABAState) SetTerminate(value bool) {
	s.Terminate = value
}
