package types

const (
	UnknownMsgType MsgType = iota
	// InitMsgType sent when DKG instance is started by requester
	InitMsgType
	// ProtocolMsgType is the DKG itself
	ProtocolMsgType
	// PartialSingatureMsgType post DKG deposit data signatures
	PartialSingatureMsgType
	// SignedDepositDataMsgType final output msg used by requester to make deposits and register validator with SSV
	SignedDepositDataMsgType
)

const (
	ethAddressSize     = 20
	ethAddressStartPos = 0
	indexSize          = 4
	indexStartPos      = ethAddressStartPos + ethAddressSize
)
