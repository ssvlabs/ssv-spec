package testingutils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ssz "github.com/ferranbt/fastssz"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

type testingKeyManager struct {
	keys           map[string]*bls.SecretKey
	ecdsaKeys      map[string]*ecdsa.PrivateKey
	encryptionKeys map[string]*rsa.PrivateKey
	domain         types.DomainType

	slashableDataRoots [][]byte
}

func NewTestingKeyManager() *testingKeyManager {
	return NewTestingKeyManagerWithSlashableRoots([][]byte{})
}

func NewTestingKeyManagerWithSlashableRoots(slashableDataRoots [][]byte) *testingKeyManager {
	ret := &testingKeyManager{
		keys:           map[string]*bls.SecretKey{},
		ecdsaKeys:      map[string]*ecdsa.PrivateKey{},
		encryptionKeys: nil,
		domain:         types.PrimusTestnet,

		slashableDataRoots: slashableDataRoots,
	}

	_ = ret.AddShare(Testing4SharesSet().ValidatorSK)
	for _, s := range Testing4SharesSet().Shares {
		_ = ret.AddShare(s)
	}

	_ = ret.AddShare(Testing7SharesSet().ValidatorSK)
	for _, s := range Testing7SharesSet().Shares {
		_ = ret.AddShare(s)
	}

	_ = ret.AddShare(Testing10SharesSet().ValidatorSK)
	for _, s := range Testing10SharesSet().Shares {
		_ = ret.AddShare(s)
	}

	_ = ret.AddShare(Testing13SharesSet().ValidatorSK)
	for _, s := range Testing13SharesSet().Shares {
		_ = ret.AddShare(s)
	}
	for _, o := range Testing4SharesSet().DKGOperators {
		ret.ecdsaKeys[o.ETHAddress.String()] = o.SK
	}
	for _, o := range Testing7SharesSet().DKGOperators {
		ret.ecdsaKeys[o.ETHAddress.String()] = o.SK
	}
	for _, o := range Testing10SharesSet().DKGOperators {
		ret.ecdsaKeys[o.ETHAddress.String()] = o.SK
	}
	for _, o := range Testing13SharesSet().DKGOperators {
		ret.ecdsaKeys[o.ETHAddress.String()] = o.SK
	}
	return ret
}

// IsAttestationSlashable returns error if attestation is slashable
func (km *testingKeyManager) IsAttestationSlashable(data *spec.AttestationData) error {
	for _, r := range km.slashableDataRoots {
		r2, _ := data.HashTreeRoot()
		if bytes.Equal(r, r2[:]) {
			return errors.New("slashable attestation")
		}
	}
	return nil
}

func (km *testingKeyManager) SignRoot(data types.Root, sigType types.SignatureType, pk []byte) (types.Signature, error) {
	if k, found := km.keys[hex.EncodeToString(pk)]; found {
		computedRoot, err := types.ComputeSigningRoot(data, types.ComputeSignatureDomain(km.domain, sigType))
		if err != nil {
			return nil, errors.Wrap(err, "could not sign root")
		}

		return k.SignByte(computedRoot).Serialize(), nil
	}
	return nil, errors.New("pk not found")
}

// IsBeaconBlockSlashable returns true if the given block is slashable
func (km *testingKeyManager) IsBeaconBlockSlashable(block *bellatrix.BeaconBlock) error {
	return nil
}

func (km *testingKeyManager) SignBeaconObject(obj ssz.HashRoot, domain spec.Domain, pk []byte) (types.Signature, []byte, error) {
	if k, found := km.keys[hex.EncodeToString(pk)]; found {
		r, err := types.ComputeETHSigningRoot(obj, domain)
		if err != nil {
			return nil, nil, errors.Wrap(err, "could not compute signing root")
		}

		sig := k.SignByte(r[:])
		blsSig := spec.BLSSignature{}
		copy(blsSig[:], sig.Serialize())

		return sig.Serialize(), r[:], nil
	}
	return nil, nil, errors.New("pk not found")
}

// Decrypt given a rsa pubkey and a PKCS1v15 cipher text byte array, returns the decrypted data
func (km *testingKeyManager) Decrypt(pk *rsa.PublicKey, cipher []byte) ([]byte, error) {
	panic("implement")
}

// Encrypt given a rsa pubkey and data returns an PKCS1v15 e
func (km *testingKeyManager) Encrypt(pk *rsa.PublicKey, data []byte) ([]byte, error) {
	return TestingEncryption(pk, data), nil
}

// SignDKGOutput signs output according to the SIP https://docs.google.com/document/d/1TRVUHjFyxINWW2H9FYLNL2pQoLy6gmvaI62KL_4cREQ/edit
func (km *testingKeyManager) SignDKGOutput(output types.Root, address common.Address) (types.Signature, error) {
	root, err := output.GetRoot()
	if err != nil {
		return nil, err
	}
	sk := km.ecdsaKeys[address.String()]
	if sk == nil {
		return nil, errors.New(fmt.Sprintf("unable to find ecdsa key for address %v", address.String()))
	}
	sig, err := crypto.Sign(root, sk)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func (km *testingKeyManager) SignETHDepositRoot(root []byte, address common.Address) (types.Signature, error) {
	panic("implemet")
}

func (km *testingKeyManager) AddShare(shareKey *bls.SecretKey) error {
	km.keys[hex.EncodeToString(shareKey.GetPublicKey().Serialize())] = shareKey
	return nil
}

func (km *testingKeyManager) RemoveShare(pubKey string) error {
	delete(km.keys, pubKey)
	return nil
}
