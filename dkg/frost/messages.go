package frost

type PreparationMessage struct {
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
	BlameData        [][]byte // SignedMessages received from the bad participant
	BlamerSessionSk  []byte
}
