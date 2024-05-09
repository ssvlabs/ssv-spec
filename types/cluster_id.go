package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"sort"
)

type ClusterID [32]byte

func (cid ClusterID) MessageIDBelongs(msgID MessageID) bool {
	id := msgID.GetDutyExecutorID()[16:]
	return bytes.Equal(id, cid[:])
}

// Return a 32 bytes ID for the cluster of operators
func GetClusterID(committee []OperatorID) ClusterID {
	// sort
	sort.Slice(committee, func(i, j int) bool {
		return committee[i] < committee[j]
	})
	// Convert to bytes
	bytes := make([]byte, len(committee)*4)
	for i, v := range committee {
		binary.LittleEndian.PutUint32(bytes[i*4:], uint32(v))
	}
	// Hash
	return sha256.Sum256(bytes)
}

// [ERR]: failed to encode Operator: failed to encode ClusterID: failed to encode ClusterID: no ssz-size or ssz-max tags found for element. tag=
// exit status 1
// types/generate.go:4: running "go": exit status 1
// make: *** [generate-ssz] Error 1
