package main

import (
	"bytes"
	"fmt"
	// spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	// "github.com/herumi/bls-eth-go-binary/bls"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/pkg/errors"
)

var TestingConfig = func(keySet *testingutils.TestKeySet) *alea.Config {
	return &alea.Config{
		Signer:    testingutils.NewTestingKeyManager(),
		SigningPK: keySet.Shares[1].GetPublicKey().Serialize(),
		Domain:    types.PrimusTestnet,
		ValueCheckF: func(data []byte) error {
			if bytes.Equal(data, testingutils.TestingInvalidValueCheck) {
				return errors.New("invalid value")
			}

			// as a base validation we do not accept nil values
			if len(data) == 0 {
				return errors.New("invalid value")
			}
			return nil
		},
		ProposerF: func(state *alea.State, round alea.Round) types.OperatorID {
			return 1
		},
		Network: testingutils.NewTestingNetworkAlea(),
		Timer:   testingutils.NewTestingTimerAlea(),
	}
}

var TestingShare = func(keysSet *testingutils.TestKeySet,shareIndex types.OperatorID) *types.Share {
	return &types.Share{
		OperatorID:          1,
		ValidatorPubKey:     keysSet.ValidatorPK.Serialize(),
		SharePubKey:         keysSet.Shares[shareIndex].GetPublicKey().Serialize(),
		DomainType:          types.PrimusTestnet,
		Quorum:              keysSet.Threshold,
		PartialQuorum:       keysSet.PartialThreshold,
		Committee:           keysSet.Committee(),
		FeeRecipientAddress: testingutils.TestingFeeRecipient,
	}
}

func main() {

	// var TestingSK = func() *bls.SecretKey {
	// 	types.InitBLS()
	// 	ret := &bls.SecretKey{}
	// 	ret.SetByCSPRNG()
	// 	return ret
	// }()
	// var testingValidatorPK = spec.BLSPubKey{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}
	// var testingShare = &types.Share{
	// 	OperatorID:      1,
	// 	ValidatorPubKey: testingValidatorPK[:],
	// 	SharePubKey:     TestingSK.GetPublicKey().Serialize(),
	// 	DomainType:      types.PrimusTestnet,
	// 	Quorum:          3,
	// 	PartialQuorum:   2,
	// 	Committee: []*types.Operator{
	// 		{
	// 			OperatorID: 1,
	// 			PubKey:     TestingSK.GetPublicKey().Serialize(),
	// 		},
	// 	},
	// }
	// fmt.Println("created Share", testingShare)
	ks := testingutils.Testing4SharesSet()
	fmt.Println("Ks structure:validator SK, validator PK, shareCount, quorum, partialQuorum, shares (map), DKGOperators map to structure with SK, ETHAddr, EncryptionKey")
	fmt.Println("created Ks", ks)

	config := TestingConfig(ks)
	fmt.Println("Config: signer, signing pk, domain, valueCheckF, ProposerF, network, timer")
	fmt.Println("created Config", config)

	controller := alea.NewController(
		[]byte{1, 2, 3, 4},
		TestingShare(ks,1),
		testingutils.TestingConfig(ks).Domain,
		config,
	)
	fmt.Println("Controller: identifier, height, stored instances, future msg container, domain, share, config")
	fmt.Println("created Controller", controller)

	controller.StartNewInstance([]byte{1})

	// var controller = NewController([]byte{1, 2, 3, 4} /*identifier*/,testingShare,types.PrimusTestnet,TestingConfig(testingutils.Testing4SharesSet()))
	// fmt.Print("created Controller")
}


