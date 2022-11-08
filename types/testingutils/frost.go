package testingutils

import (
	crand "crypto/rand"
	"encoding/hex"
	"math/big"
	mrand "math/rand"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

type TestOutcome struct {
	KeygenOutcome TestKeygenOutcome
	BlameOutcome  TestBlameOutcome
}

type TestKeygenOutcome struct {
	ValidatorPK     string
	Share           map[uint32]string
	OperatorPubKeys map[uint32]string
}

func (o TestKeygenOutcome) ToKeygenOutcomeMap(threshold uint64, operators []uint32) map[uint32]*dkg.KeyGenOutput {
	m := make(map[uint32]*dkg.KeyGenOutput)

	opPublicKeys := make(map[types.OperatorID]*bls.PublicKey)
	for _, operatorID := range operators {

		pk := &bls.PublicKey{}
		_ = pk.DeserializeHexStr(o.OperatorPubKeys[operatorID])
		opPublicKeys[types.OperatorID(operatorID)] = pk

		share := o.Share[operatorID]
		sk := &bls.SecretKey{}
		_ = sk.DeserializeHexStr(share)

		vk, _ := hex.DecodeString(o.ValidatorPK)

		m[operatorID] = &dkg.KeyGenOutput{
			Share:           sk,
			ValidatorPK:     vk,
			OperatorPubKeys: opPublicKeys,
			Threshold:       threshold,
		}
	}

	return m
}

func ResetRandSeed() {
	src := mrand.NewSource(1)
	src.Seed(12345)
	crand.Reader = mrand.New(src)
}

func GetRandRequestID() dkg.RequestID {
	requestID := dkg.RequestID{}
	for i := range requestID {
		rndInt, _ := crand.Int(crand.Reader, big.NewInt(255))
		if len(rndInt.Bytes()) == 0 {
			requestID[i] = 0
		} else {
			requestID[i] = rndInt.Bytes()[0]
		}
	}
	return requestID
}

type TestBlameOutcome struct {
	Valid        bool
	BlameMessage []byte
}
