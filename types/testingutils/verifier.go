package testingutils

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"sync"

	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

type testingVerifier struct {
	signaturesCache map[types.OperatorID]map[[32]byte][]byte
}

var (
	testingVerifierInstance             *testingVerifier
	testingVerifierInstanceLock         sync.Mutex
	testingVerifierSingletonConstructor sync.Once
)

func NewTestingVerifier() types.SignatureVerifier {
	testingVerifierSingletonConstructor.Do(func() {
		testingVerifierInstance = &testingVerifier{
			signaturesCache: make(map[types.OperatorID]map[[32]byte][]byte),
		}
	})
	return testingVerifierInstance
}

func (v *testingVerifier) Verify(msg *types.SignedSSVMessage, operators []*types.Operator) error {

	encodedMsg, err := msg.SSVMessage.Encode()
	if err != nil {
		return err
	}

	// Get message hash
	hash := sha256.Sum256(encodedMsg)

	// Find operator that matches ID with the signer and verify signature
	for i, signer := range msg.OperatorIDs {
		if err := v.VerifySignatureForSigner(hash, msg.Signatures[i], signer, operators); err != nil {
			return err
		}
	}

	return nil
}

func (v *testingVerifier) VerifySignatureForSigner(root [32]byte, signature []byte, signer types.OperatorID, operators []*types.Operator) error {

	for _, op := range operators {
		// Find signer
		if signer == op.OperatorID {

			// Check cache
			if v.HasSignature(op.OperatorID, root, signature) {
				return nil
			}

			// Get public key
			pk, err := types.PemToPublicKey(op.SSVOperatorPubKey)
			if err != nil {
				return errors.Wrap(err, "could not parse signer public key")
			}

			// Verify
			err = rsa.VerifyPKCS1v15(pk, crypto.SHA256, root[:], signature)

			if err == nil {
				v.SaveSignature(op.OperatorID, root, signature)
			}
			return err
		}
	}
	return errors.New("unknown signer")
}

func (v *testingVerifier) HasSignature(operatorID types.OperatorID, root [32]byte, signature []byte) bool {
	testingVerifierInstanceLock.Lock()
	defer testingVerifierInstanceLock.Unlock()
	if _, found := v.signaturesCache[operatorID]; !found {
		return false
	}

	storedSignature, found := v.signaturesCache[operatorID][root]
	if !found {
		return false
	}

	return bytes.Equal(storedSignature[:], signature[:])
}

func (v *testingVerifier) SaveSignature(operatorID types.OperatorID, root [32]byte, signature []byte) {
	testingVerifierInstanceLock.Lock()
	defer testingVerifierInstanceLock.Unlock()
	if _, found := v.signaturesCache[operatorID]; !found {
		v.signaturesCache[operatorID] = make(map[[32]byte][]byte)
	}
	v.signaturesCache[operatorID][root] = signature
}

// Verifies a list of SignedSSVMessage using the operators list
func VerifyListOfSignedSSVMessages(msgs []*types.SignedSSVMessage, operators []*types.Operator) error {
	verifier := NewTestingVerifier()

	for _, msg := range msgs {
		err := verifier.Verify(msg, operators)
		if err != nil {
			return err
		}
	}
	return nil
}
