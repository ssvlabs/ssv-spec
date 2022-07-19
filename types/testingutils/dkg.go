package testingutils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/dkg/base"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var TestingWithdrawalCredentials, _ = hex.DecodeString("0x010000000000000000000000535953b5a6040074948cf185eaa7d2abbd66808f")
/*
var TestingDKGNode = func(keySet *TestKeySet) *dkg.Node {
	network := NewTestingNetwork()
	config := &dkg.Config{
		Protocol: func(init *dkg.Init, operatorID types.OperatorID, identifier dkg.RequestID) dkg.Protocol {
			ret := stubdkg.New(init, identifier, dkg.ProtocolConfig{
				Identifier:    identifier,
				Operator:      &dkg.Operator{
					OperatorID:       operatorID,
					ETHAddress:       common.Address{},
					EncryptionPubKey: nil,
				}, // TODO: Fix this
				BeaconNetwork: "",
				Signer:        nil,
			})
			ret.(*stubdkg.DKG).SetOperators(
				Testing4SharesSet().ValidatorPK.Serialize(),
				Testing4SharesSet().Shares,
			)
			return ret
		},
		Network:             network,
		Storage:             NewTestingStorage(),
		SignatureDomainType: types.PrimusTestnet,
		Signer:              NewTestingKeyManager(),
	}

	return dkg.NewNode(&dkg.Operator{
		OperatorID:       1,
		ETHAddress:       keySet.DKGOperators[1].ETHAddress,
		EncryptionPubKey: &keySet.DKGOperators[1].EncryptionKey.PublicKey,
	}, config)
}
*/

var SignDKGMsg = func(sk *ecdsa.PrivateKey, id types.OperatorID, msg *base.Message) *base.Message {
	domain := types.PrimusTestnet
	sigType := types.DKGSignatureType

	r, _ := types.ComputeSigningRoot(msg, types.ComputeSignatureDomain(domain, sigType))
	sig, _ := crypto.Sign(r, sk)

	return &base.Message{
		Header: msg.Header,
		Data:   msg.Data,
		Signature: sig,
	}
}

var InitMessageDataBytes = func(operators []types.OperatorID, threshold uint16, withdrawalCred []byte) []byte {
	m := &base.Init{
		OperatorIDs:           operators,
		Threshold:             threshold,
		WithdrawalCredentials: withdrawalCred,
	}
	byts, _ := m.Encode()
	return byts
}

//var ProtocolMsgDataBytes = func(stage stubdkg.Stage) []byte {
//	d := &stubdkg.ProtocolMsg{
//		Stage: stage,
//	}
//
//	ret, _ := d.Encode()
//	return ret
//}
