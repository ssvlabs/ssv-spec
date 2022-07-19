package types

import (
	"encoding/binary"
	"github.com/ethereum/go-ethereum/common"
)

type RequestID [24]byte

func (msg RequestID) GetETHAddress() common.Address {
	ret := common.Address{}
	copy(ret[:], msg[ethAddressStartPos:ethAddressStartPos+ethAddressSize])
	return ret
}

func (msg RequestID) GetRoleType() uint32 {
	indexByts := msg[indexStartPos : indexStartPos+indexSize]
	return binary.LittleEndian.Uint32(indexByts)
}

func NewRequestID(ethAddress common.Address, index uint32) RequestID {
	indexByts := make([]byte, 4)
	binary.LittleEndian.PutUint32(indexByts, index)

	ret := RequestID{}
	copy(ret[ethAddressStartPos:ethAddressStartPos+ethAddressSize], ethAddress[:])
	copy(ret[indexStartPos:indexStartPos+indexSize], indexByts[:])
	return ret
}
