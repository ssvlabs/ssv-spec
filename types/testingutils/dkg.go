package testingutils

import (
	"crypto/ecdsa"
	"encoding/hex"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/stubdkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var TestingWithdrawalCredentials, _ = hex.DecodeString("010000000000000000000000535953b5a6040074948cf185eaa7d2abbd66808f")
var TestingForkVersion = types.PraterNetwork.ForkVersion()

var TestingDKGNode = func(keySet *TestKeySet) *dkg.Node {
	network := NewTestingNetwork()
	km := NewTestingKeyManager()
	config := &dkg.Config{
		Protocol: func(network dkg.Network, operatorID types.OperatorID, identifier dkg.RequestID) dkg.KeyGenProtocol {
			return &MockKeygenProtocol{
				KeyGenOutput: keySet.KeyGenOutput(1),
			}
		},
		Network:             network,
		Storage:             NewTestingStorage(),
		SignatureDomainType: types.PrimusTestnet,
		Signer:              km,
		Verifier:            km,
	}

	return dkg.NewNode(&dkg.Operator{
		OperatorID:       1,
		ETHAddress:       keySet.DKGOperators[1].ETHAddress,
		EncryptionPubKey: &keySet.DKGOperators[1].EncryptionKey.PublicKey,
	}, config)
}

var SignDKGMsg = func(sk *ecdsa.PrivateKey, id types.OperatorID, msg *dkg.Message) *dkg.SignedMessage {
	domain := types.PrimusTestnet
	sigType := types.DKGSignatureType

	r, _ := types.ComputeSigningRoot(msg, types.ComputeSignatureDomain(domain, sigType))
	sig, _ := crypto.Sign(r, sk)

	return &dkg.SignedMessage{
		Message:   msg,
		Signer:    id,
		Signature: sig,
	}
}

var InitMessageDataBytes = func(operators []types.OperatorID, threshold uint16, withdrawalCred []byte, fork spec.Version) []byte {
	m := &dkg.Init{
		OperatorIDs:           operators,
		Threshold:             threshold,
		WithdrawalCredentials: withdrawalCred,
		Fork:                  fork,
	}
	byts, _ := m.Encode()
	return byts
}

var ProtocolMsgDataBytes = func(stage stubdkg.Stage) []byte {
	d := &stubdkg.ProtocolMsg{
		Stage: stage,
	}

	ret, _ := d.Encode()
	return ret
}

var PartialDepositDataBytes = func(signer types.OperatorID, root []byte, sk *bls.SecretKey) []byte {
	d := &dkg.PartialDepositData{
		Signer:    signer,
		Root:      root,
		Signature: sk.SignByte(root).Serialize(),
	}
	ret, _ := d.Encode()
	return ret
}

var DespositDataSigningRoot = func(keySet *TestKeySet, initMsg *dkg.Init) []byte {
	root, _, _ := types.GenerateETHDepositData(
		keySet.ValidatorPK.Serialize(),
		initMsg.WithdrawalCredentials,
		initMsg.Fork,
		types.DomainDeposit,
	)
	return root
}

var SignedOutputObject = func(requestID dkg.RequestID, signer types.OperatorID, root []byte, address common.Address, share *bls.SecretKey, validatorSk *bls.SecretKey) *dkg.SignedOutput {
	// TODO: Move FakeEncryption and FakeEcdsaSign to before calling this function?
	o := &dkg.Output{
		RequestID:            requestID,
		EncryptedShare:       FakeEncryption(share.Serialize()),
		SharePubKey:          share.GetPublicKey().Serialize(),
		ValidatorPubKey:      validatorSk.GetPublicKey().Serialize(),
		DepositDataSignature: validatorSk.SignByte(root).Serialize(),
	}
	root1, _ := o.GetRoot()
	return &dkg.SignedOutput{
		Data:      o,
		Signer:    signer,
		Signature: FakeEcdsaSign(root1, address[:]),
	}
}

var SignedOutputBytes = func(requestID dkg.RequestID, signer types.OperatorID, root []byte, address common.Address, share *bls.SecretKey, validatorSk *bls.SecretKey) []byte {
	d := SignedOutputObject(requestID, signer, root, address, share, validatorSk)
	ret, _ := d.Encode()
	return ret
}
