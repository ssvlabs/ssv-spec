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
		BatchSize: 2,
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

	controller2 := alea.NewController(
		[]byte{1, 2, 3, 4},
		TestingShare(ks,2),
		testingutils.TestingConfig(ks).Domain,
		config,
	)
	fmt.Println("Controller: identifier, height, stored instances, future msg container, domain, share, config")
	fmt.Println("created Controller2", controller2)

	controller.StartNewInstance([]byte{1})
	instance := controller.InstanceForHeight(0)
	fmt.Println("Instance:", instance)

	controller2.StartNewInstance([]byte{2})
	instance2 := controller2.InstanceForHeight(0)
	fmt.Println("Instance2:", instance2)

	

	proposal1, err := alea.CreateProposal(instance.State,instance.GetConfig(),[]byte{1,1})
	if err != nil {
		errors.Wrap(err, "could not create proposal message")
	}
	controller.ProcessMsg(proposal1)
	proposal2, err := alea.CreateProposal(instance.State,instance.GetConfig(),[]byte{1,2})
	controller.ProcessMsg(proposal2)


	proposalData1,err := proposal1.Message.GetProposalData()
	if err != nil {
		errors.Wrap(err, "could not get proposal data from proposal message created")
	}
	proposalData2,err := proposal2.Message.GetProposalData()


	vcbcMessage,err := alea.CreateVCBC(instance.State,instance.GetConfig(),[]*alea.ProposalData{proposalData1,proposalData2},0)
	if err != nil {
		errors.Wrap(err, "could not create vcbc message")
		return
	}
	controller2.ProcessMsg(vcbcMessage)



	proposal3, err := alea.CreateProposal(instance.State,instance.GetConfig(),[]byte{1,3})
	controller.ProcessMsg(proposal3)
	proposal4, err := alea.CreateProposal(instance.State,instance.GetConfig(),[]byte{1,4})
	controller.ProcessMsg(proposal4)

	proposalData3,err := proposal3.Message.GetProposalData()
	proposalData4,err := proposal4.Message.GetProposalData()
	vcbcMessage2,err := alea.CreateVCBC(instance.State,instance.GetConfig(),[]*alea.ProposalData{proposalData3,proposalData4},0)
	controller2.ProcessMsg(vcbcMessage2)

	// k := 0
	// for {
	// 	if instance.State.StopAgreement {
	// 		break
	// 	} else {
	// 		k += 1
	// 	}
	// 	if (k == 1e9) {
	// 		instance.State.StopAgreement = true
	// 	}
	// }
}


