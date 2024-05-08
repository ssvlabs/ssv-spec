package testingutils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"sync"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	ssz "github.com/ferranbt/fastssz"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

type testingKeyManager struct {
	keys           map[string]*bls.SecretKey
	ecdsaKeys      map[string]*ecdsa.PrivateKey
	encryptionKeys map[string]*rsa.PrivateKey
	domain         types.DomainType

	slashableDataRoots [][]byte
}

var (
	instancesMap = make(map[uint64]*testingKeyManager)
	mu           sync.Mutex
)

func getHash(data [][]byte) uint64 {
	h := fnv.New64a()
	for _, d := range data {
		h.Write(d)
	}
	return h.Sum64()
}

func NewTestingKeyManager() *testingKeyManager {
	return NewTestingKeyManagerWithSlashableRoots([][]byte{})
}

func NewTestingKeyManagerWithSlashableRoots(slashableDataRoots [][]byte) *testingKeyManager {

	hash := getHash(slashableDataRoots)

	mu.Lock()
	defer mu.Unlock()

	if instance, ok := instancesMap[hash]; ok {
		return instance
	}

	ret := &testingKeyManager{
		keys:           map[string]*bls.SecretKey{},
		ecdsaKeys:      map[string]*ecdsa.PrivateKey{},
		encryptionKeys: nil,
		domain:         TestingSSVDomainType,

		slashableDataRoots: slashableDataRoots,
	}

	testingSharesSets := []*TestKeySet{Testing4SharesSet(), Testing7SharesSet(), Testing10SharesSet(), Testing13SharesSet()}

	for _, testingShareSet := range testingSharesSets {
		_ = ret.AddShare(testingShareSet.ValidatorSK)
		for _, s := range testingShareSet.Shares {
			_ = ret.AddShare(s)
		}
		for _, o := range testingShareSet.DKGOperators {
			ret.ecdsaKeys[o.ETHAddress.String()] = o.SK
		}
	}

	for _, keySet := range TestingKeySetMap {
		_ = ret.AddShare(keySet.ValidatorSK)
		for _, s := range keySet.Shares {
			_ = ret.AddShare(s)
		}
		for _, o := range keySet.DKGOperators {
			ret.ecdsaKeys[o.ETHAddress.String()] = o.SK
		}
	}

	instancesMap[hash] = ret

	return ret
}

// IsAttestationSlashable returns error if attestation is slashable
func (km *testingKeyManager) IsAttestationSlashable(pk []byte, data *spec.AttestationData) error {
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

		return k.SignByte(computedRoot[:]).Serialize(), nil
	}
	return nil, errors.New("pk not found")
}

// IsBeaconBlockSlashable returns error if the given block is slashable
func (km *testingKeyManager) IsBeaconBlockSlashable(pk []byte, slot spec.Slot) error {
	return nil
}

func (km *testingKeyManager) SignBeaconObject(obj ssz.HashRoot, domain spec.Domain, pk []byte, domainType spec.DomainType) (types.Signature, [32]byte, error) {
	if k, found := km.keys[hex.EncodeToString(pk)]; found {
		r, err := types.ComputeETHSigningRoot(obj, domain)
		if err != nil {
			return nil, [32]byte{}, errors.Wrap(err, "could not compute signing root")
		}

		sig := k.SignByte(r[:])
		blsSig := spec.BLSSignature{}
		copy(blsSig[:], sig.Serialize())

		return sig.Serialize(), r, nil
	}
	return nil, [32]byte{}, errors.New("pk not found")
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
	sig, err := ethcrypto.Sign(root[:], sk)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func (km *testingKeyManager) SignETHDepositRoot(root []byte, address common.Address) (types.Signature, error) {
	panic("implement")
}

func (km *testingKeyManager) AddShare(shareKey *bls.SecretKey) error {
	km.keys[hex.EncodeToString(shareKey.GetPublicKey().Serialize())] = shareKey
	return nil
}

func (km *testingKeyManager) RemoveShare(pubKey string) error {
	delete(km.keys, pubKey)
	return nil
}
