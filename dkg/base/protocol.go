package base

// Protocol is an interface for all DKG protocol to support a variety of protocols for future upgrades
type Protocol interface {
	Start() ([]Message, error)
	// ProcessMsg returns true and a bls share if finished
	ProcessMsg(msg *Message) ([]Message, error)
	Output() ([]byte, error)
}
