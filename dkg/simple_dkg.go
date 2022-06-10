package dkg

// SimpleDKG is a simple DKG protocol with no verification
type SimpleDKG struct {
	network Network
}

func NewSimpleDKG(network Network) Protocol {
	return &SimpleDKG{
		network: network,
	}
}

func (s *SimpleDKG) Start() error {
	panic("implement")
}

func (s *SimpleDKG) ProcessMsg(msg *SignedMessage) (bool, *ProtocolOutput, error) {
	panic("implement")
}
