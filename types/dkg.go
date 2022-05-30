package types

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
)

type DKGInit struct {
	OperatorIDs           []OperatorID
	WithdrawalCredentials []byte
}

type DKGOutput struct {
	EncryptedShare        []byte
	DKGSize               uint16
	ValidatorPubKey       ValidatorPK
	WithdrawalCredentials []byte
}

func (o *DKGOutput) GetRoot() ([]byte, error) {
	uint16Solidity, _ := abi.NewType("uint16", "", nil)
	bytesSolidity, _ := abi.NewType("bytes", "", nil)

	arguments := abi.Arguments{
		{
			Type: bytesSolidity,
		},
		{
			Type: uint16Solidity,
		},
		{
			Type: bytesSolidity,
		},
		{
			Type: bytesSolidity,
		},
	}

	bytes, _ := arguments.Pack(
		o.EncryptedShare,
		o.DKGSize,
		o.ValidatorPubKey,
		o.WithdrawalCredentials,
	)

	return crypto.Keccak256(bytes), nil
}
