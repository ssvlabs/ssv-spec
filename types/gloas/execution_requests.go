package gloas

import (
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// BuilderDepositRequest is the EIP-8282 builder deposit request — a fixed-size container in the Gloas
// ExecutionRequests.
type BuilderDepositRequest struct {
	Pubkey                phase0.BLSPubKey `ssz-size:"48"`
	WithdrawalCredentials [32]byte         `ssz-size:"32"`
	Amount                phase0.Gwei
	Signature             phase0.BLSSignature `ssz-size:"96"`
}

// BuilderExitRequest is the EIP-8282 builder exit request — a fixed-size container in the Gloas
// ExecutionRequests.
type BuilderExitRequest struct {
	SourceAddress bellatrix.ExecutionAddress `ssz-size:"20"`
	Pubkey        phase0.BLSPubKey           `ssz-size:"48"`
}

// ExecutionRequests is the Gloas execution requests: the Electra three (deposits, withdrawals,
// consolidations) plus the EIP-8282 builder deposit/exit requests. A Gloas CL encodes all five lists, so
// electra.ExecutionRequests (three) marshals a block two offsets short and a Gloas CL rejects the §4 submit
// as invalid SSZ — hence this five-list variant. List bounds are the spec MAX_* values.
type ExecutionRequests struct {
	Deposits        []*electra.DepositRequest       `ssz-index:"0" ssz-type:"progressive-list"`
	Withdrawals     []*electra.WithdrawalRequest    `ssz-index:"1" ssz-type:"progressive-list"`
	Consolidations  []*electra.ConsolidationRequest `ssz-index:"2" ssz-type:"progressive-list"`
	BuilderDeposits []*BuilderDepositRequest        `ssz-index:"3" ssz-type:"progressive-list"`
	BuilderExits    []*BuilderExitRequest           `ssz-index:"4" ssz-type:"progressive-list"`
}
