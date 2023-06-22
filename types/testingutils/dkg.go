package testingutils

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"strconv"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/stubdkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var TestingWithdrawalCredentials, _ = hex.DecodeString("010000000000000000000000535953b5a6040074948cf185eaa7d2abbd66808f")
var TestingForkVersion = types.PraterNetwork.ForkVersion()

var TestingDKGNode = func(keySet *TestKeySet) *dkg.Node {
	network := NewTestingNetwork()
	km := NewTestingKeyManager()
	config := &dkg.Config{
		KeygenProtocol: func(dkg.RequestID, types.OperatorID, dkg.IConfig, *dkg.Init) dkg.Protocol {
			return &TestingKeygenProtocol{
				KeyGenOutput: keySet.KeyGenOutput(1),
			}
		},
		ReshareProtocol: func(dkg.RequestID, types.OperatorID, dkg.IConfig, *dkg.Reshare, *dkg.ReshareParams) dkg.Protocol {
			return &TestingKeygenProtocol{
				KeyGenOutput: keySet.KeyGenOutput(1),
			}
		},
		Network:             network,
		Storage:             NewTestingStorage(),
		SignatureDomainType: types.PrimusTestnet,
		Signer:              km,
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
	byts, _ := InitMessageData(operators, threshold, withdrawalCred, fork).Encode()
	return byts
}

var InitMessageData = func(operators []types.OperatorID, threshold uint16, withdrawalCred []byte, fork spec.Version) *dkg.Init {
	return &dkg.Init{
		OperatorIDs:           operators,
		Threshold:             threshold,
		WithdrawalCredentials: withdrawalCred,
		Fork:                  fork,
	}
}

var ReshareMessageDataBytes = func(operators []types.OperatorID, threshold uint16, validatorPK types.ValidatorPK, oldOperators []types.OperatorID) []byte {
	byts, _ := ReshareMessageData(operators, threshold, validatorPK, oldOperators).Encode()
	return byts
}

var ReshareMessageData = func(operators []types.OperatorID, threshold uint16, validatorPK types.ValidatorPK, oldOperators []types.OperatorID) *dkg.Reshare {
	return &dkg.Reshare{
		ValidatorPK:    validatorPK,
		OperatorIDs:    operators,
		Threshold:      threshold,
		OldOperatorIDs: oldOperators,
	}
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
var (
	encryptedDataCache = map[string][]byte{}
	decryptedDataCache = map[string][]byte{}
)

func TestingEncryption(pk *rsa.PublicKey, data []byte) []byte {
	id := hex.EncodeToString(pk.N.Bytes()) + fmt.Sprintf("%x", pk.E) + hex.EncodeToString(data)
	if found := encryptedDataCache[id]; found != nil {
		return found
	}
	cipherText, _ := types.Encrypt(pk, data)
	encryptedDataCache[id] = cipherText
	return cipherText
}

func TestingDecryption(sk *rsa.PrivateKey, data []byte) []byte {
	id := hex.EncodeToString(sk.N.Bytes()) + fmt.Sprintf("%x", sk.E) + hex.EncodeToString(data)
	if found := decryptedDataCache[id]; found != nil {
		return found
	}
	plaintext, _ := types.Decrypt(sk, data)
	decryptedDataCache[id] = plaintext
	return plaintext
}

func (ks *TestKeySet) KeyGenOutput(opId types.OperatorID) *dkg.KeyGenOutput {
	opPks := make(map[types.OperatorID]*bls.PublicKey)
	for id, share := range ks.Shares {
		opPks[id] = share.GetPublicKey()
	}

	return &dkg.KeyGenOutput{
		Share:           ks.Shares[opId],
		OperatorPubKeys: opPks,
		ValidatorPK:     ks.ValidatorPK.Serialize(),
		Threshold:       ks.Threshold,
	}
}

var (
	signedOutputCache = map[string]*dkg.SignedOutput{}
)

func (ks *TestKeySet) SignedOutputObject(requestID dkg.RequestID, opId types.OperatorID, root []byte) *dkg.SignedOutput {
	id := hex.EncodeToString(requestID[:]) + strconv.FormatUint(uint64(opId), 10) + hex.EncodeToString(root)
	if found := signedOutputCache[id]; found != nil {
		return found
	}
	share := ks.Shares[opId]
	o := &dkg.Output{
		RequestID:       requestID,
		EncryptedShare:  TestingEncryption(&ks.DKGOperators[opId].EncryptionKey.PublicKey, []byte("0x"+share.SerializeToHexStr())),
		SharePubKey:     share.GetPublicKey().Serialize(),
		ValidatorPubKey: ks.ValidatorPK.Serialize(),
	}
	if root != nil {
		o.DepositDataSignature = ks.ValidatorSK.SignByte(root).Serialize()
	}
	// root1, _ := o.GetRoot()
	root1, _ := types.ComputeSigningRoot(o, types.ComputeSignatureDomain(types.PrimusTestnet, types.DKGSignatureType))

	sig, _ := crypto.Sign(root1, ks.DKGOperators[opId].SK)

	ret := &dkg.SignedOutput{
		Data:      o,
		Signer:    opId,
		Signature: sig,
	}
	signedOutputCache[id] = ret
	return ret
}

func (ks *TestKeySet) SignedOutputBytes(requestID dkg.RequestID, opId types.OperatorID, root []byte) []byte {
	d := ks.SignedOutputObject(requestID, opId, root)
	ret, _ := d.Encode()
	return ret
}

func (ks *TestKeySet) SignedKeySignOutputObject(requestID dkg.RequestID, opID types.OperatorID, signingRoot []byte) *dkg.SignedOutput {
	id := hex.EncodeToString(requestID[:]) + strconv.FormatUint(uint64(opID), 10) + hex.EncodeToString(signingRoot)
	if found := signedOutputCache[id]; found != nil {
		return found
	}

	sk := ks.ValidatorSK
	o := &dkg.KeySignOutput{
		RequestID:   requestID,
		Signature:   sk.SignByte(signingRoot).Serialize(),
		ValidatorPK: ks.ValidatorPK.Serialize(),
	}

	root, _ := types.ComputeSigningRoot(o, types.ComputeSignatureDomain(types.PrimusTestnet, types.DKGSignatureType))
	sig, _ := crypto.Sign(root, ks.DKGOperators[opID].SK)

	ret := &dkg.SignedOutput{
		KeySignData: o,
		Signer:      opID,
		Signature:   sig,
	}
	signedOutputCache[id] = ret
	return ret
}

func (ks *TestKeySet) SignedKeySignOutputBytes(requestID dkg.RequestID, opID types.OperatorID, signingRoot []byte) []byte {
	d := ks.SignedKeySignOutputObject(requestID, opID, signingRoot)
	ret, _ := d.Encode()
	return ret
}
