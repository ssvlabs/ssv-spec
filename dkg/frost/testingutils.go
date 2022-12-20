package frost

import (
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	ecies "github.com/ecies/go/v2"
)

func Testing_PreparationMessageBytes(id types.OperatorID, frMsgStore testingutils.FrostMsgStore) []byte {
	pk, _ := hex.DecodeString(frMsgStore.SessionPKs[id])
	msg := &ProtocolMsg{
		Round: Preparation,
		PreparationMessage: &PreparationMessage{
			SessionPk: pk,
		},
	}
	byts, _ := msg.Encode()
	return byts
}

func Testing_Round1MessageBytes(id types.OperatorID, frMsgStore testingutils.FrostMsgStore) []byte {
	commitments := make([][]byte, 0)
	for _, commitment := range frMsgStore.Round1[id].Commitments {
		cbytes, _ := hex.DecodeString(commitment)
		commitments = append(commitments, cbytes)
	}
	proofS, _ := hex.DecodeString(frMsgStore.Round1[id].ProofS)
	proofR, _ := hex.DecodeString(frMsgStore.Round1[id].ProofR)
	shares := map[uint32][]byte{}
	for peerID, share := range frMsgStore.Round1[id].Shares {
		shareBytes, _ := hex.DecodeString(share)
		shares[peerID] = shareBytes
	}
	msg := ProtocolMsg{
		Round: Round1,
		Round1Message: &Round1Message{
			Commitment: commitments,
			ProofS:     proofS,
			ProofR:     proofR,
			Shares:     shares,
		},
	}
	byts, _ := msg.Encode()
	return byts
}

func Testing_Round2MessageBytes(id types.OperatorID, frMsgStore testingutils.FrostMsgStore) []byte {
	vk, _ := hex.DecodeString(frMsgStore.Round2[id].Vk)
	vkshare, _ := hex.DecodeString(frMsgStore.Round2[id].VkShare)
	msg := ProtocolMsg{
		Round: Round2,
		Round2Message: &Round2Message{
			Vk:      vk,
			VkShare: vkshare,
		},
	}
	byts, _ := msg.Encode()
	return byts
}

func BlameMessageBytes(id types.OperatorID, blameType BlameType, blameMessages []*dkg.SignedMessage) []byte {
	blameData := make([][]byte, 0)
	for _, blameMessage := range blameMessages {
		byts, _ := blameMessage.Encode()
		blameData = append(blameData, byts)
	}

	skBytes, _ := hex.DecodeString(testingutils.KeygenMsgStore.SessionSKs[1])
	sk := ecies.NewPrivateKeyFromBytes(skBytes)

	ret, _ := (&ProtocolMsg{
		Round: Blame,
		BlameMessage: &BlameMessage{
			Type:             blameType,
			TargetOperatorID: uint32(id),
			BlameData:        blameData,
			BlamerSessionSk:  sk.Bytes(),
		},
	}).Encode()
	return ret
}
