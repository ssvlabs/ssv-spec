package types

import (
	"encoding/binary"
	ssz "github.com/ferranbt/fastssz"
)

type SSZUint64 uint64

// HashTreeRoot --
func (s SSZUint64) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(s))
	var root [32]byte
	copy(root[:], buf)
	return root, nil
}

// HashTreeRootWith --
func (s SSZUint64) HashTreeRootWith(hh *ssz.Hasher) error {
	indx := hh.Index()
	hh.PutUint64(uint64(s))
	hh.Merkleize(indx)
	return nil
}

// SSZBytes --
type SSZBytes []byte

// HashTreeRoot --
func (b SSZBytes) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(b)
}

// HashTreeRootWith --
func (b SSZBytes) HashTreeRootWith(hh *ssz.Hasher) error {
	indx := hh.Index()
	hh.PutBytes(b)
	hh.Merkleize(indx)
	return nil
}
