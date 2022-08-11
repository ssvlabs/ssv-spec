package testingutils

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
