package dkg

// Network is a collection of funcs for DKG
type Network interface {
	// StreamDKGOutput will stream to any subscriber the result of the DKG
	StreamDKGOutput(output *SignedOutput) error
	// Broadcast will broadcast a msg to the dkg network
	Broadcast(msg *SignedMessage) error
}
