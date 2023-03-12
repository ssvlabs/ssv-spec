package types

import (
	"encoding/binary"
	"fmt"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	ssz "github.com/ferranbt/fastssz"
)

type SSZUint64 uint64

func (s SSZUint64) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(s)
}

func (s SSZUint64) HashTreeRootWith(hh ssz.HashWalker) error {
	indx := hh.Index()
	hh.PutUint64(uint64(s))
	hh.Merkleize(indx)
	return nil
}

// HashTreeRoot --
func (s SSZUint64) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(s))
	var root [32]byte
	copy(root[:], buf)
	return root, nil
}

// SSZBytes --
type SSZBytes []byte

func (b SSZBytes) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(b)
}

func (b SSZBytes) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(b)
}

func (b SSZBytes) HashTreeRootWith(hh ssz.HashWalker) error {
	indx := hh.Index()
	hh.PutBytes(b)
	hh.Merkleize(indx)
	return nil
}

// SSZTransactions --
type SSZTransactions []bellatrix.Transaction

func (b SSZTransactions) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(b)
}

func (b SSZTransactions) HashTreeRootWith(hh ssz.HashWalker) error {
	// taken from https://github.com/prysmaticlabs/prysm/blob/develop/encoding/ssz/htrutils.go#L97-L119
	subIndx := hh.Index()
	num := uint64(len(b))
	if num > 1048576 {
		return ssz.ErrIncorrectListSize
	}
	for _, elem := range b {
		{
			elemIndx := hh.Index()
			byteLen := uint64(len(elem))
			if byteLen > 1073741824 {
				return ssz.ErrIncorrectListSize
			}
			hh.AppendBytes32(elem)
			hh.MerkleizeWithMixin(elemIndx, byteLen, (1073741824+31)/32)
		}
	}
	hh.MerkleizeWithMixin(subIndx, num, 1048576)
	return nil
}

// HashTreeRoot --
func (b SSZTransactions) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(b)
}

// SSZ32Bytes --
type SSZ32Bytes [32]byte

func (b SSZ32Bytes) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(b)
}

func (b SSZ32Bytes) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(b)
}

func (b SSZ32Bytes) HashTreeRootWith(hh ssz.HashWalker) error {
	indx := hh.Index()
	hh.PutBytes(b[:])
	hh.Merkleize(indx)
	return nil
}

// UnmarshalSSZ --
func (b SSZ32Bytes) UnmarshalSSZ(buf []byte) error {
	if len(buf) != b.SizeSSZ() {
		return fmt.Errorf("expected buffer of length %d receiced %d", b.SizeSSZ(), len(buf))
	}
	copy(b[:], buf[:])
	return nil
}

// MarshalSSZTo --
func (b SSZ32Bytes) MarshalSSZTo(dst []byte) ([]byte, error) {
	return append(dst, b[:]...), nil
}

// MarshalSSZ --
func (b SSZ32Bytes) MarshalSSZ() ([]byte, error) {
	return b[:], nil
}

// SizeSSZ returns the size of the serialized object.
func (b SSZ32Bytes) SizeSSZ() int {
	return 32
}
