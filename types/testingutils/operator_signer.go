package testingutils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"sync"

	"github.com/bloxapp/ssv-spec/types"
)

type testingOperatorSigner struct {
	SSVOperatorKey map[string]*rsa.PrivateKey
}

var (
	testingOperatorSignerInstance *testingOperatorSigner
	opMu                          sync.Mutex
)

func NewTestingOperatorSigner() *testingOperatorSigner {

	opMu.Lock()
	defer opMu.Unlock()

	if testingOperatorSignerInstance != nil {
		return testingOperatorSignerInstance
	}

	ret := &testingOperatorSigner{
		SSVOperatorKey: map[string]*rsa.PrivateKey{},
	}

	testingSharesSets := []*TestKeySet{Testing4SharesSet(), Testing7SharesSet(), Testing10SharesSet(), Testing13SharesSet()}

	for _, testingShareSet := range testingSharesSets {
		for _, k := range testingShareSet.OperatorKeys {
			_ = ret.AddSSVOperatorKey(k)
		}
	}

	testingOperatorSignerInstance = ret

	return ret
}

func (km *testingOperatorSigner) SignSSVMessage(data []byte, pk []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	sk := km.SSVOperatorKey[hex.EncodeToString(pk)]
	signature, err := rsa.SignPKCS1v15(rand.Reader, sk, crypto.SHA256, hash[:])
	if err != nil {
		return []byte{}, err
	}
	return signature, nil
}

func (km *testingOperatorSigner) AddSSVOperatorKey(sk *rsa.PrivateKey) error {
	pem, err := types.GetPublicKeyPem(sk)
	if err != nil {
		panic(err)
	}
	km.SSVOperatorKey[hex.EncodeToString(pem)] = sk
	return nil
}

func (km *testingOperatorSigner) RemoveSSVOperatorKey(pubKey string) {
	delete(km.SSVOperatorKey, pubKey)
}
