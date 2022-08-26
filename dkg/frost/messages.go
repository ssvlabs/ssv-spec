package frost

type PreparationRoundMessage struct {
	SessionPk []byte
}

type Round1Message struct {
	Commitment [][]byte
	ProofS     []byte
	ProofR     []byte
	Shares     [][]byte
}

type Round2Message struct {
	Vk      []byte
	VkShare []byte
}

type BlameMessage struct {
	TargetOperatorID uint64
	BlameData        []byte // Received signed Round1Message
	BlamerSessionSk  []byte
}
