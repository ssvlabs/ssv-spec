package keygen

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/frost/frostutils"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var skFromHex = func(str string) *bls.SecretKey {
	types.InitBLS()
	ret := &bls.SecretKey{}
	if err := ret.DeserializeHexStr(str); err != nil {
		panic(err.Error())
	}
	return ret
}

func SignedOutputBytes(requestID dkg.RequestID, opId types.OperatorID, sk *ecdsa.PrivateKey, pk *rsa.PublicKey, root []byte) []byte {
	d := SignedOutputObject(requestID, opId, sk, pk, root)
	ret, _ := d.Encode()
	return ret
}

func SignedOutputObject(requestID dkg.RequestID, opId types.OperatorID, opSK *ecdsa.PrivateKey, opPK *rsa.PublicKey, root []byte) *dkg.SignedOutput {
	share := skFromHex(frostutils.KeygenMsgStore.Round2[opId].SkShare)
	validatorPublicKey, _ := hex.DecodeString(frostutils.KeygenMsgStore.Round2[1].Vk)

	sk := &bls.SecretKey{}
	_ = sk.DeserializeHexStr(frostutils.KeygenMsgStore.Round2[opId].Sk)

	o := &dkg.Output{
		RequestID:            requestID,
		EncryptedShare:       testingutils.TestingEncryption(opPK, share.Serialize()),
		SharePubKey:          share.GetPublicKey().Serialize(),
		ValidatorPubKey:      validatorPublicKey,
		DepositDataSignature: sk.SignByte(root).Serialize(),
	}

	root1, _ := o.GetRoot()

	sig, _ := crypto.Sign(root1, opSK)

	ret := &dkg.SignedOutput{
		Data:      o,
		Signer:    opId,
		Signature: sig,
	}
	return ret
}
