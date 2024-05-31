package testingutils

import (
	"bytes"
	"crypto/sha256"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

var TestingQBFTFullData = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9}
var TestingQBFTRootData = func() [32]byte {
	return sha256.Sum256(TestingQBFTFullData)
}()

var TestingCutOffRound = qbft.Round(15)

var TestingConfig = func(keySet *TestKeySet) *qbft.Config {
	return &qbft.Config{
		OperatorSigner: NewTestingOperatorSigner(keySet, 1),
		SigningPK:      keySet.Shares[1].GetPublicKey().Serialize(),
		Domain:         TestingSSVDomainType,
		ValueCheckF: func(data []byte) error {
			if bytes.Equal(data, TestingInvalidValueCheck) {
				return errors.New("invalid value")
			}

			// as a base validation we do not accept nil values
			if len(data) == 0 {
				return errors.New("invalid value")
			}
			return nil
		},
		ProposerF: func(state *qbft.State, round qbft.Round) types.OperatorID {
			return 1
		},
		Network:           NewTestingNetwork(1, keySet.OperatorKeys[1]),
		Timer:             NewTestingTimer(),
		SignatureVerifier: NewTestingVerifier(),
		CutOffRound:       TestingCutOffRound,
	}
}

var TestingInvalidValueCheck = []byte{1, 1, 1, 1}

var TestingGraffiti = [32]byte{1}

var TestingValidatorShare = func(keysSet *TestKeySet, opID types.OperatorID) types.ValidatorShare {

	operatorPKBytes, err := types.MarshalPublicKey(keysSet.OperatorKeys[opID])
	if err != nil {
		panic(err)
	}

	return types.ValidatorShare{
		OperatorID:     opID,
		OperatorPubKey: operatorPKBytes,
		SharePubKey:    keysSet.Shares[opID].GetPublicKey().Serialize(),
	}
}

var TestingShare = func(keysSet *TestKeySet, valIdx phase0.ValidatorIndex) *types.SharedValidator {

	// Decode validator public key
	pkBytesSlice := keysSet.ValidatorPK.Serialize()
	pkBytesArray := [48]byte{}
	copy(pkBytesArray[:], pkBytesSlice)

	operatorsCommittee := make([]types.OperatorID, 0)
	for _, op := range keysSet.Committee() {
		operatorsCommittee = append(operatorsCommittee, op.OperatorID)
	}

	return &types.SharedValidator{
		ValidatorIndex:      valIdx,
		ValidatorPubKey:     pkBytesArray,
		OwnValidatorShare:   TestingValidatorShare(keysSet, 1),
		CommitteeID:         types.GetCommitteeID(operatorsCommittee),
		Committee:           keysSet.Committee(),
		Quorum:              keysSet.Threshold,
		PartialQuorum:       keysSet.PartialThreshold,
		DomainType:          TestingSSVDomainType,
		FeeRecipientAddress: TestingFeeRecipient,
		Graffiti:            TestingGraffiti[:],
	}
}

var BaseInstance = func() *qbft.Instance {
	return baseInstance(TestingShare(Testing4SharesSet(), TestingValidatorIndex), Testing4SharesSet(), []byte{1, 2, 3, 4})
}

var SevenOperatorsInstance = func() *qbft.Instance {
	return baseInstance(TestingShare(Testing7SharesSet(), TestingValidatorIndex), Testing7SharesSet(), []byte{1, 2, 3, 4})
}

var TenOperatorsInstance = func() *qbft.Instance {
	return baseInstance(TestingShare(Testing10SharesSet(), TestingValidatorIndex), Testing10SharesSet(), []byte{1, 2, 3, 4})
}

var ThirteenOperatorsInstance = func() *qbft.Instance {
	return baseInstance(TestingShare(Testing13SharesSet(), TestingValidatorIndex), Testing13SharesSet(), []byte{1, 2, 3, 4})
}

var baseInstance = func(share *types.SharedValidator, keySet *TestKeySet, identifier []byte) *qbft.Instance {
	ret := qbft.NewInstance(TestingConfig(keySet), share, identifier, qbft.FirstHeight)
	ret.StartValue = TestingQBFTFullData
	return ret
}

func NewTestingQBFTController(
	identifier []byte,
	share *types.SharedValidator,
	config qbft.IConfig,
) *qbft.Controller {
	return qbft.NewController(
		identifier,
		share,
		config,
	)
}
