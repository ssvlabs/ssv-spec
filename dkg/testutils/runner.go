package testutils

import (
	"encoding/json"
	"fmt"
	"github.com/bloxapp/ssv-spec/dkg"
	dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func TestNode(dataset DepositSignDataSet, index uint64) *dkg.Node {
	jstr, _ := json.Marshal(dataset)
	//o := dataset.DKGOperators[types.OperatorID(index)]
	fmt.Printf("data set is %v\n", string(jstr))
	operator := &dkgtypes.Operator{
		OperatorID:       types.OperatorID(index),
		ETHAddress:       dataset.DKGOperators[types.OperatorID(index)].ETHAddress,
		EncryptionPubKey: &dataset.DKGOperators[types.OperatorID(index)].EncryptionKey.PublicKey,
	}

	var protocol = func(init *dkgtypes.Init, operatorID types.OperatorID, identifier dkgtypes.RequestID) dkgtypes.Protocol {
		return mockProtocol{localKeyShare: dataset.MakeLocalKeyShare(index)}
	}
	config := &dkgtypes.Config{
		Protocol:            protocol,
		BeaconNetwork:       types.PraterNetwork,
		Network:             newMockNetwork(),
		Storage:             newMockStorage(dataset.TestKeySet),
		SignatureDomainType: types.PrimusTestnet,
		Signer: &MockSigner{
			SK:            dataset.DKGOperators[types.OperatorID(index)].SK,
			ETHAddress:    dataset.DKGOperators[types.OperatorID(index)].ETHAddress,
			EncryptionKey: dataset.DKGOperators[types.OperatorID(index)].EncryptionKey,
		},
	}
	return dkg.NewNode(operator, config)
}

func TestConfig(dataset DepositSignDataSet, index uint64) *dkgtypes.Config {
	var protocol = func(init *dkgtypes.Init, operatorID types.OperatorID, identifier dkgtypes.RequestID) dkgtypes.Protocol {
		return mockProtocol{localKeyShare: dataset.MakeLocalKeyShare(index)}
	}
	return &dkgtypes.Config{
		Protocol:            protocol,
		BeaconNetwork:       types.PraterNetwork,
		Network:             newMockNetwork(),
		Storage:             newMockStorage(dataset.TestKeySet),
		SignatureDomainType: types.PrimusTestnet,
		Signer: &MockSigner{
			SK:            dataset.DKGOperators[types.OperatorID(index)].SK,
			ETHAddress:    dataset.DKGOperators[types.OperatorID(index)].ETHAddress,
			EncryptionKey: dataset.DKGOperators[types.OperatorID(index)].EncryptionKey,
		},
	}
}

func TestRunner(dataset DepositSignDataSet, keyset testingutils.TestKeySet, index uint64) *dkg.Runner {
	kgProtocol := mockProtocol{localKeyShare: dataset.MakeLocalKeyShare(index)}
	partyData := dataset.PartyData[types.OperatorID(index)]
	threshold := uint64(len(partyData.Coefficients) - 1)
	init := dkgtypes.Init{
		OperatorIDs:           dataset.IndicesVec(),
		Threshold:             threshold,
		WithdrawalCredentials: TestingWithdrawalCredentials,
		Fork:                  TestingForkVersion[:],
	}
	return &dkg.Runner{
		Operator: &dkgtypes.Operator{
			OperatorID:       types.OperatorID(index),
			ETHAddress:       keyset.DKGOperators[types.OperatorID(index)].ETHAddress,
			EncryptionPubKey: &keyset.DKGOperators[types.OperatorID(index)].EncryptionKey.PublicKey,
		},
		InitMsg:               &init,
		Identifier:            dkgtypes.RequestID{},
		DepositDataRoot:       nil,
		DepositDataSignatures: nil,
		PartialSignatures:     nil,
		OutputMsgs:            map[types.OperatorID]*dkgtypes.ParsedSignedDepositDataMessage{},
		KeygenSubProtocol:     kgProtocol,
		Config: &dkgtypes.Config{
			Protocol:            nil,
			BeaconNetwork:       "",
			Network:             newMockNetwork(),
			Storage:             nil,
			SignatureDomainType: nil,
			Signer: &MockSigner{
				SK:            keyset.DKGOperators[types.OperatorID(index)].SK,
				ETHAddress:    keyset.DKGOperators[types.OperatorID(index)].ETHAddress,
				EncryptionKey: keyset.DKGOperators[types.OperatorID(index)].EncryptionKey,
			},
		},
	}
}
