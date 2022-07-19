package base

const (
	// InitMsgType sent when DKG instance is started by requester
	InitMsgType MsgType = iota
	// ProtocolMsgType is the DKG itself
	ProtocolMsgType
	KeygenOutputType
	PartialOutputMsgType
	// DepositDataMsgType post DKG deposit data signatures
	DepositDataMsgType
	// OutputMsgType final output msg used by requester to make deposits and register validator with SSV
	OutputMsgType
)

const (
	ethAddressSize     = 20
	ethAddressStartPos = 0
	indexSize          = 4
	indexStartPos      = ethAddressStartPos + ethAddressSize
)

