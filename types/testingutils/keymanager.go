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

type SignOutput struct {
	Signature []byte
	Root      [32]byte
}

type TestingKeyManager struct {
	keys               map[string]*bls.SecretKey
	ecdsaKeys          map[string]*ecdsa.PrivateKey
	encryptionKeys     map[string]*rsa.PrivateKey
	domain             types.DomainType
	signatureCache     map[string]map[string]map[spec.Domain]*SignOutput
	slashableDataRoots map[string][][]byte
}

var (
	instancesMap = make(map[uint64]*TestingKeyManager)
	mu           sync.Mutex
)

func getHash(data map[string][][]byte) uint64 {
	h := fnv.New64a()
	for k, roots := range data {
		kBytes, _ := hex.DecodeString(k)
		h.Write(kBytes)
		for _, d := range roots {
			h.Write(d)
		}
	}
	return h.Sum64()
}

func NewTestingKeyManager() *TestingKeyManager {
	return NewTestingKeyManagerWithSlashableRoots(map[string][][]byte{})
}

func NewTestingKeyManagerWithSlashableRoots(slashableDataRoots map[string][][]byte) *TestingKeyManager {

	hash := getHash(slashableDataRoots)

	mu.Lock()
	defer mu.Unlock()

	if instance, ok := instancesMap[hash]; ok {
		// clone the instance so we can mutate it
		return instance.Clone()
	}

	ret := &TestingKeyManager{
		keys:           map[string]*bls.SecretKey{},
		ecdsaKeys:      map[string]*ecdsa.PrivateKey{},
		encryptionKeys: nil,
		domain:         TestingSSVDomainType,

		slashableDataRoots: slashableDataRoots,

		signatureCache: make(map[string]map[string]map[spec.Domain]*SignOutput),
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

// AddSlashableDataRoot adds a slashable data root to the key manager
func (km *TestingKeyManager) AddSlashableDataRoot(pk types.ShareValidatorPK, dataRoot []byte) {
	entry := hex.EncodeToString(pk)
	if km.slashableDataRoots[entry] == nil {
		km.slashableDataRoots[entry] = make([][]byte, 0)
	}
	km.slashableDataRoots[entry] = append(km.slashableDataRoots[entry], dataRoot)
}

// Clone returns a new key manager with the same keys
func (km *TestingKeyManager) Clone() *TestingKeyManager {
	ret := &TestingKeyManager{
		keys:               map[string]*bls.SecretKey{},
		ecdsaKeys:          map[string]*ecdsa.PrivateKey{},
		encryptionKeys:     map[string]*rsa.PrivateKey{},
		domain:             km.domain,
		signatureCache:     make(map[string]map[string]map[spec.Domain]*SignOutput),
		slashableDataRoots: make(map[string][][]byte),
	}

	for k, v := range km.keys {
		ret.keys[k] = v
	}

	for k, v := range km.ecdsaKeys {
		ret.ecdsaKeys[k] = v
	}

	// copy encryption keys
	for k, v := range km.encryptionKeys {
		ret.encryptionKeys[k] = v
	}

	// copy slashable data roots
	for k, v := range km.slashableDataRoots {
		ret.slashableDataRoots[k] = make([][]byte, len(v))
		ret.slashableDataRoots[k] = v
	}

	return ret
}

// IsAttestationSlashable returns error if attestation is slashable
func (km *TestingKeyManager) IsAttestationSlashable(pk types.ShareValidatorPK, data *spec.AttestationData) error {
	entry := hex.EncodeToString(pk)
	for _, r := range km.slashableDataRoots[entry] {
		r2, _ := data.HashTreeRoot()
		if bytes.Equal(r, r2[:]) {
			return errors.New("slashable attestation")
		}
	}
	return nil
}

func (km *TestingKeyManager) SignRoot(data types.Root, sigType types.SignatureType, pk []byte) (types.Signature, error) {
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
func (km *TestingKeyManager) IsBeaconBlockSlashable(pk []byte, slot spec.Slot) error {
	return nil
}

func (km *TestingKeyManager) SignBeaconObject(obj ssz.HashRoot, domain spec.Domain, pk []byte, domainType spec.DomainType) (types.Signature, [32]byte, error) {
	mu.Lock()
	defer mu.Unlock()

	pkString := hex.EncodeToString(pk)

	if k, found := km.keys[pkString]; found {

		if signOutput, has := km.hasSignRequest(pkString, obj, domain); has {
			return signOutput.Signature, signOutput.Root, nil
		}

		r, err := types.ComputeETHSigningRoot(obj, domain)
		if err != nil {
			return nil, [32]byte{}, errors.Wrap(err, "could not compute signing root")
		}

		sig := k.SignByte(r[:])
		blsSig := spec.BLSSignature{}
		copy(blsSig[:], sig.Serialize())

		sigString := sig.Serialize()

		km.storeSignRequest(pkString, obj, domain, sigString, r)

		return sigString, r, nil
	}
	return nil, [32]byte{}, errors.New("pk not found")
}

// SignDKGOutput signs output according to the SIP https://docs.google.com/document/d/1TRVUHjFyxINWW2H9FYLNL2pQoLy6gmvaI62KL_4cREQ/edit
func (km *TestingKeyManager) SignDKGOutput(output types.Root, address common.Address) (types.Signature, error) {
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

func (km *TestingKeyManager) SignETHDepositRoot(root []byte, address common.Address) (types.Signature, error) {
	panic("implement")
}

func (km *TestingKeyManager) AddShare(shareKey *bls.SecretKey) error {
	km.keys[hex.EncodeToString(shareKey.GetPublicKey().Serialize())] = shareKey
	return nil
}

func (km *TestingKeyManager) RemoveShare(pubKey string) error {
	delete(km.keys, pubKey)
	return nil
}

func (km *TestingKeyManager) hasSignRequest(pk string, obj ssz.HashRoot, domain spec.Domain) (*SignOutput, bool) {
	if _, exists := km.signatureCache[pk]; !exists {
		return &SignOutput{}, false
	}
	objRoot, err := obj.HashTreeRoot()
	if err != nil {
		return &SignOutput{}, false
	}
	root := hex.EncodeToString(objRoot[:])
	if _, exists := km.signatureCache[pk][root]; !exists {
		return &SignOutput{}, false
	}
	if _, exists := km.signatureCache[pk][root][domain]; !exists {
		return &SignOutput{}, false
	}
	return km.signatureCache[pk][root][domain], true
}

func (km *TestingKeyManager) storeSignRequest(pk string, obj ssz.HashRoot, domain spec.Domain, sig types.Signature, r [32]byte) {
	if _, exists := km.signatureCache[pk]; !exists {
		km.signatureCache[pk] = make(map[string]map[spec.Domain]*SignOutput)
	}
	objRoot, err := obj.HashTreeRoot()
	if err != nil {
		panic(err)
	}
	root := hex.EncodeToString(objRoot[:])
	if _, exists := km.signatureCache[pk][root]; !exists {
		km.signatureCache[pk][root] = make(map[spec.Domain]*SignOutput)
	}
	km.signatureCache[pk][root][domain] = &SignOutput{Signature: sig, Root: r}
}
