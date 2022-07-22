package testutils

import (
	"github.com/bloxapp/ssv-spec/dkg"
	dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func TestRunner(dataset DepositSignDataSet, keyset testingutils.TestKeySet, index uint64) *dkg.Runner {
	kgProtocol := mockProtocol{localKeyShare: dataset.MakeLocalKeyShare(index)}
	partyData := dataset.PartyData[types.OperatorID(index)]
	threshold := uint64(len(partyData.Coefficients) - 1)
	var opIds []types.OperatorID
	for _, id := range dataset.IndicesVec() {
		opIds = append(opIds, types.OperatorID(id))
	}
	init := dkgtypes.Init{
		Nonce:                 0,
		OperatorIDs:           opIds,
		Threshold:             threshold,
		WithdrawalCredentials: TestingWithdrawalCredentials,
		Fork:                  TestingForkVersion,
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
