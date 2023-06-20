package common

// ProtocolRound is enum for all the rounds in the protocol
type ProtocolRound int

const (
	Uninitialized ProtocolRound = iota
	Preparation
	Round1
	Round2
	KeygenOutput
	Blame
	Timeout
	KeysignOutput
)

var rounds = []ProtocolRound{
	Uninitialized,
	Preparation,
	Round1,
	Round2,
	KeygenOutput,
	Blame,
	Timeout,
	KeysignOutput,
}

func (round ProtocolRound) String() string {
	m := map[ProtocolRound]string{
		Uninitialized: "Uninitialized",
		Preparation:   "Preparation",
		Round1:        "Round1",
		Round2:        "Round2",
		KeygenOutput:  "KeygenOutput",
		Blame:         "Blame",
		Timeout:       "Timeout",
		KeysignOutput: "KeysignOutput",
	}
	return m[round]
}
