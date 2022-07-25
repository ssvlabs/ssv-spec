package testutils

import (
	"crypto/ecdsa"
	"encoding/hex"
	dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var TestingWithdrawalCredentials, _ = hex.DecodeString("010000000000000000000000535953b5a6040074948cf185eaa7d2abbd66808f")
var TestingForkVersion = types.PraterNetwork.ForkVersion()
var TestingAddress = common.HexToAddress("535953b5a6040074948cf185eaa7d2abbd66808f")
var TestingRequestID = dkgtypes.NewRequestID(TestingAddress, 1)

func h2b(str string) []byte {
	b, _ := hex.DecodeString(str)
	return b
}

func FakeEncryption(data []byte) []byte {
	out := []byte("__fake_encrypted(")
	out = append(out, data...)
	out = append(out, []byte(")")...)
	return out
}

func FakeEcdsaSign(root []byte, address []byte) []byte {
	out := []byte("__fake_ecdsa_sign(root=")
	out = append(out, root...)
	out = append(out, []byte(",address=")...)
	out = append(out, address...)
	out = append(out, []byte(")")...)
	return out
}

var SignDKGMsgRoot = func(sk *ecdsa.PrivateKey, msg types.Root) []byte {
	domain := types.PrimusTestnet
	sigType := types.DKGSignatureType

	r, _ := types.ComputeSigningRoot(msg, types.ComputeSignatureDomain(domain, sigType))
	sig, _ := crypto.Sign(r, sk)
	return sig
}

var SignDKGMsg = func(sk *ecdsa.PrivateKey, msg dkgtypes.Signable) dkgtypes.Signable {
	sig := SignDKGMsgRoot(sk, msg)
	_ = msg.SetSignature(sig)

	return msg
}
