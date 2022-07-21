package testutils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type MockSigner struct {
	SK            *ecdsa.PrivateKey
	ETHAddress    common.Address
	EncryptionKey *rsa.PrivateKey
}

func (m MockSigner) Decrypt(pk *rsa.PublicKey, cipher []byte) ([]byte, error) {
	if *pk != m.EncryptionKey.PublicKey {
		return nil, errors.New("public key doesn't match the signer's")
	}
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, m.EncryptionKey, cipher, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func (m MockSigner) Encrypt(pk *rsa.PublicKey, data []byte) ([]byte, error) {
	//ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pk, data)
	//if err != nil {
	//	return nil, err
	//}
	return FakeEncryption(data), nil
}

func (m MockSigner) SignRoot(data types.Root, sigType types.SignatureType, pk []byte) (types.Signature, error) {
	panic("implement me")
}

func (m MockSigner) SignDKGOutput(output types.Root, address common.Address) (types.Signature, error) {
	if address != m.ETHAddress {
		return nil, errors.New("address doesn't match the signer's")
	}
	sig := SignDKGMsgRoot(m.SK, output)
	return sig, nil
}

func (m MockSigner) SignETHDepositRoot(root []byte, pk []byte) (types.Signature, error) {
	panic("implement me")
}
