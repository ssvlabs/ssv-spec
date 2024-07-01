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

var TestingShare = func(keysSet *TestKeySet, valIdx phase0.ValidatorIndex) *types.Share {

	// Decode validator public key
	pkBytesSlice := keysSet.ValidatorPK.Serialize()
	pkBytesArray := [48]byte{}
	copy(pkBytesArray[:], pkBytesSlice)

	return &types.Share{
		ValidatorIndex:      valIdx,
		ValidatorPubKey:     pkBytesArray,
		SharePubKey:         keysSet.Shares[1].GetPublicKey().Serialize(),
		Committee:           keysSet.Committee(),
		DomainType:          TestingSSVDomainType,
		FeeRecipientAddress: TestingFeeRecipient,
		Graffiti:            TestingGraffiti[:],
	}
}

var TestingCommitteeMember = func(keysSet *TestKeySet) *types.CommitteeMember {
	operators := []*types.Operator{}

	for _, key := range keysSet.Committee() {

		// Encode member's public key
		pkBytes, err := types.GetPublicKeyPem(keysSet.OperatorKeys[key.Signer])
		if err != nil {
			panic(err)
		}

		operators = append(operators, &types.Operator{
			OperatorID:        key.Signer,
			SSVOperatorPubKey: pkBytes,
		})
	}

	opIds := []types.OperatorID{}
	for _, key := range keysSet.Committee() {
		opIds = append(opIds, key.Signer)
	}

	operatorPubKeyBytes, err := types.GetPublicKeyPem(keysSet.OperatorKeys[1])
	if err != nil {
		panic(err)
	}

	return &types.CommitteeMember{
		OperatorID:        1,
		CommitteeID:       types.GetCommitteeID(opIds),
		SSVOperatorPubKey: operatorPubKeyBytes,
		FaultyNodes:       (keysSet.Threshold - 1) / 2,
		Committee:         operators,
		DomainType:        TestingSSVDomainType,
	}
}

var BaseInstance = func() *qbft.Instance {
	return baseInstance(TestingCommitteeMember(Testing4SharesSet()), Testing4SharesSet(), []byte{1, 2, 3, 4})
}

var SevenOperatorsInstance = func() *qbft.Instance {
	return baseInstance(TestingCommitteeMember(Testing7SharesSet()), Testing7SharesSet(), []byte{1, 2, 3, 4})
}

var TenOperatorsInstance = func() *qbft.Instance {
	return baseInstance(TestingCommitteeMember(Testing10SharesSet()), Testing10SharesSet(), []byte{1, 2, 3, 4})
}

var ThirteenOperatorsInstance = func() *qbft.Instance {
	return baseInstance(TestingCommitteeMember(Testing13SharesSet()), Testing13SharesSet(), []byte{1, 2, 3, 4})
}

var baseInstance = func(committeeMember *types.CommitteeMember, keySet *TestKeySet, identifier []byte) *qbft.Instance {
	ret := qbft.NewInstance(TestingConfig(keySet), committeeMember, identifier, qbft.FirstHeight)
	ret.StartValue = TestingQBFTFullData
	return ret
}

func NewTestingQBFTController(
	identifier []byte,
	share *types.CommitteeMember,
	config qbft.IConfig,
) *qbft.Controller {
	return qbft.NewController(
		identifier,
		share,
		config,
	)
}
