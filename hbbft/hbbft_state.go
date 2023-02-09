package hbbft

import (
	"math/rand"
	"time"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

type HBBFTState struct {
	BatchSize int
	Round     ACSRound
	Buffer    []*TransactionData
	ACSState  map[ACSRound]*ACSState
	DECMsgs   map[ACSRound]map[types.OperatorID]map[types.OperatorID][]byte
}

func NewHBBFTState() *HBBFTState {
	hbbftState := &HBBFTState{
		BatchSize: 1,
		Round:     FirstACSRound,
		ACSState:  make(map[ACSRound]*ACSState),
		DECMsgs:   make(map[ACSRound]map[types.OperatorID]map[types.OperatorID][]byte),
	}
	hbbftState.InitMaps(hbbftState.Round)
	return hbbftState
}

func (s *HBBFTState) IncrementRound() {
	s.Round += 1
	s.InitMaps(s.Round)
}

func (s *HBBFTState) InitMaps(round ACSRound) {
	s.DECMsgs[round] = make(map[types.OperatorID]map[types.OperatorID][]byte)
	s.ACSState[round] = NewACSState(round)
}

func (s *HBBFTState) GetBatchSize() int {
	return s.BatchSize
}
func (s *HBBFTState) GetRound() ACSRound {
	return s.Round
}
func (s *HBBFTState) GetABAState(round ACSRound) *ABAState {
	return s.ACSState[round].GetABAState()
}
func (s *HBBFTState) GetCurrentABAState() *ABAState {
	return s.ACSState[s.Round].GetABAState()
}

func (s *HBBFTState) GetRandomTransacations(n int) []*TransactionData {
	if len(s.Buffer) <= n {
		return s.Buffer
	}

	var bufferSlice []*TransactionData
	if len(s.Buffer) <= s.BatchSize {
		bufferSlice = s.Buffer
	} else {
		bufferSlice = s.Buffer[:s.BatchSize]
	}

	rand.Seed(time.Now().UnixMilli())

	bufLen := len(bufferSlice)
	for idx, v := range bufferSlice {
		if idx >= n {
			break
		}
		idx2 := rand.Intn(bufLen)
		bufferSlice[idx] = bufferSlice[idx2]
		bufferSlice[idx2] = v
	}

	return bufferSlice[:n]
}

// func (s *HBBFTState) EncryptProposed(transactions []*TransactionData, valPK types.ValidatorPK) ([]byte, error) {

// 	encodedTransactions, err := EncodeTransactions(transactions)
// 	if err != nil {
// 		return nil, err
// 	}
// 	PK, err := types.PemToPublicKey(valPK)
// 	if err != nil {
// 		return nil, err
// 	}
// 	ans, err := types.Encrypt(PK, encodedTransactions)
// 	if err != nil {
// 		fmt.Println("EncryptProposed: error with encryption", err)
// 		return nil, errors.Wrap(err, "could not encrypt transaction list")
// 	}
// 	return ans, nil
// }

func (s *HBBFTState) RunACS(proposed_encrypted []byte) map[types.OperatorID][]byte {
	// input value to RBC

	s.ACSState[s.Round] = NewACSState(s.Round)
	v := s.ACSState[s.Round].StartACS(proposed_encrypted)
	return v
}

// func (s *HBBFTState) DecryptShare(data byte[], SharePubKey []byte) {

// 	r, _ := types.ComputeSigningRoot(msg, types.ComputeSignatureDomain(domain, sigType))
// 	sig := sk.SignByte(r)

// 	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
// }

func (s *HBBFTState) GetLenDECmsgs(round ACSRound, author types.OperatorID) uint64 {
	if _, exists := s.DECMsgs[round][author]; !exists {
		return 0
	}
	return uint64(len(s.DECMsgs[round][author]))
}

func (s *HBBFTState) GetDECMsgsMap(round ACSRound, author types.OperatorID) map[types.OperatorID][]byte {
	if _, exists := s.DECMsgs[round][author]; !exists {
		return make(map[types.OperatorID][]byte)
	}
	return s.DECMsgs[round][author]
}

// func (s *HBBFTState) DecriptDECSet() {
// 	return nil
// }

// func (s *HBBFTState) SortY() {
// 	return nil
// }

func (s *HBBFTState) StoreBlockAndUpdateBuffer(y map[types.OperatorID][]*TransactionData) {
	for _, transactionsData := range y {
		for _, transactionData := range transactionsData {
			for idx := 0; idx < len(s.Buffer); idx++ {
				tData := s.Buffer[idx]
				if len(tData.Data) == len(transactionData.Data) {
					equal := true
					for data_idx, data_val := range tData.Data {
						if data_val != transactionData.Data[data_idx] {
							equal = false
							break
						}
					}
					if equal {
						s.Buffer = append(s.Buffer[:idx], s.Buffer[idx+1:]...)
						idx -= 1
					}
				}
			}
		}
	}
}
