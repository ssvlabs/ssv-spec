package testingutils

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/stretchr/testify/require"
)

const (
	MaxDenebBlock           = 4194304
	MaxConsensusData        = 4219064
	MaxConsensusDataEncoded = 4219184
)

func TestMaxBeaconVote(t *testing.T) {

	root := [32]byte{1, 2, 3, 4}
	bv := types.BeaconVote{
		BlockRoot: root,
		Source: &phase0.Checkpoint{
			Epoch: 10,
			Root:  root,
		},
		Target: &phase0.Checkpoint{
			Epoch: 20,
			Root:  root,
		},
	}

	cdBytes, err := bv.Encode()
	require.NoError(t, err)
	fmt.Printf("BeaconVote encoded length: %v\n", len(cdBytes))
}

func TestMaxConsensusData(t *testing.T) {

	// Pre consensus messages
	pSigMessagesList := make([]*types.PartialSignatureMessages, 0)
	for i := uint64(1); i <= 13; i++ {
		// Pre-consensus message (13 signatures at most)
		pSigMsgList := make([]*types.PartialSignatureMessage, 0)
		for j := uint64(1); j <= 13; j++ {

			signingRoot := [32]byte{1, 2, 3, 4}
			strByt := [96]byte{1, 2, 3, 4}
			pSigMsg := types.PartialSignatureMessage{
				PartialSignature: strByt[:],
				SigningRoot:      signingRoot,
				Signer:           1,
				ValidatorIndex:   10,
			}

			pSigMsgList = append(pSigMsgList, &pSigMsg)
		}

		pSigMsgs := types.PartialSignatureMessages{
			Type:     types.RandaoPartialSig,
			Slot:     phase0.Slot(i),
			Messages: pSigMsgList,
		}
		pSigMessagesList = append(pSigMessagesList, &pSigMsgs)
	}

	// Biggest DataSSZ
	biggestBeaconObject := [MaxDenebBlock]byte{1}

	// Beacon Duty
	blsPK := phase0.BLSPubKey{}
	blsPKStr := [48]byte{1, 2, 3, 4}
	copy(blsPK[:], blsPKStr[:])
	validatorSCIndices := [13]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	duty := types.BeaconDuty{
		Type:                          types.BNRoleAttester,
		PubKey:                        blsPK,
		Slot:                          12,
		ValidatorIndex:                100,
		CommitteeIndex:                20,
		CommitteeLength:               128,
		CommitteesAtSlot:              10,
		ValidatorCommitteeIndex:       22,
		ValidatorSyncCommitteeIndices: validatorSCIndices[:],
	}

	// Consensus Data
	cd := types.ConsensusData{
		Duty:                       duty,
		Version:                    spec.DataVersionAltair,
		PreConsensusJustifications: pSigMessagesList,
		DataSSZ:                    biggestBeaconObject[:],
	}

	cdBytes, err := cd.Encode()
	require.NoError(t, err)
	fmt.Printf("ConsensusData encoded length: %v\n", len(cdBytes))
}

func TestMaxPrepareMessage(t *testing.T) {

	ks := Testing13SharesSet()
	msgID := TestingMessageID

	// Max ConsensusData Encoded
	biggestBeaconObject := [MaxConsensusDataEncoded]byte{1}
	fullData := biggestBeaconObject[:]

	// qbft.Message
	prepareQbftMsg := &qbft.Message{
		MsgType:                  qbft.PrepareMsgType,
		Height:                   qbft.FirstHeight,
		Round:                    1,
		Identifier:               msgID[:],
		Root:                     sha256.Sum256(fullData),
		RoundChangeJustification: make([][]byte, 0),
		PrepareJustification:     make([][]byte, 0),
	}

	prepareQbftMsgBytes, err := prepareQbftMsg.Encode()
	require.NoError(t, err)
	fmt.Printf("Prepare qbft.Message encoded length: %v\n", len(prepareQbftMsgBytes))

	// SSVMessage
	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   msgID,
		Data:    prepareQbftMsgBytes,
	}

	// SignedSSVMessage
	signature, err := SignSSVMessage(ks.OperatorKeys[1], ssvMsg)
	require.NoError(t, err)

	signedSSVMessage := &types.SignedSSVMessage{
		Signatures:  [][]byte{signature},
		OperatorIDs: []types.OperatorID{1},
		SSVMessage:  ssvMsg,
	}

	signedSSVMessageBytes, err := signedSSVMessage.Encode()
	require.NoError(t, err)
	fmt.Printf("Prepare types.SignedSSVMessage encoded length: %v\n", len(signedSSVMessageBytes))
}

func TestMaxCommitMessage(t *testing.T) {

	ks := Testing13SharesSet()
	msgID := TestingMessageID

	// Max ConsensusData Encoded
	biggestBeaconObject := [MaxConsensusDataEncoded]byte{1}
	fullData := biggestBeaconObject[:]

	// qbft.Message
	qbftMsg := &qbft.Message{
		MsgType:                  qbft.CommitMsgType,
		Height:                   qbft.FirstHeight,
		Round:                    1,
		Identifier:               msgID[:],
		Root:                     sha256.Sum256(fullData),
		RoundChangeJustification: make([][]byte, 0),
		PrepareJustification:     make([][]byte, 0),
	}

	qbftMsgBytes, err := qbftMsg.Encode()
	require.NoError(t, err)
	fmt.Printf("Commit qbft.Message encoded length: %v\n", len(qbftMsgBytes))

	// SSVMessage
	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   msgID,
		Data:    qbftMsgBytes,
	}

	// SignedSSVMessage
	signature, err := SignSSVMessage(ks.OperatorKeys[1], ssvMsg)
	require.NoError(t, err)

	signedSSVMessage := &types.SignedSSVMessage{
		Signatures:  [][]byte{signature},
		OperatorIDs: []types.OperatorID{1},
		SSVMessage:  ssvMsg,
	}

	signedSSVMessageBytes, err := signedSSVMessage.Encode()
	require.NoError(t, err)
	fmt.Printf("Commit types.SignedSSVMessage encoded length: %v\n", len(signedSSVMessageBytes))
}

func TestMaxDecidedMessage(t *testing.T) {

	ks := Testing13SharesSet()
	msgID := TestingMessageID

	// Max ConsensusData Encoded
	biggestBeaconObject := [MaxConsensusDataEncoded]byte{1}
	fullData := biggestBeaconObject[:]

	// qbft.Message
	qbftMsg := &qbft.Message{
		MsgType:                  qbft.CommitMsgType,
		Height:                   qbft.FirstHeight,
		Round:                    1,
		Identifier:               msgID[:],
		Root:                     sha256.Sum256(fullData),
		RoundChangeJustification: make([][]byte, 0),
		PrepareJustification:     make([][]byte, 0),
	}

	qbftMsgBytes, err := qbftMsg.Encode()
	require.NoError(t, err)
	fmt.Printf("Commit qbft.Message encoded length: %v\n", len(qbftMsgBytes))

	// SSVMessage
	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   msgID,
		Data:    qbftMsgBytes,
	}

	// SignedSSVMessage
	signature, err := SignSSVMessage(ks.OperatorKeys[1], ssvMsg)
	require.NoError(t, err)

	sigs := make([][]byte, 0)
	opIds := make([]types.OperatorID, 0)
	for i := uint64(1); i <= 13; i++ {
		sigs = append(sigs, signature)
		opIds = append(opIds, i)
	}

	signedSSVMessage := &types.SignedSSVMessage{
		Signatures:  sigs,
		OperatorIDs: opIds,
		SSVMessage:  ssvMsg,
		FullData:    fullData,
	}

	signedSSVMessageBytes, err := signedSSVMessage.Encode()
	require.NoError(t, err)
	fmt.Printf("Decided types.SignedSSVMessage encoded length: %v\n", len(signedSSVMessageBytes))
}

func TestMaxRoundChangeMessage(t *testing.T) {

	ks := Testing13SharesSet()
	msgID := TestingMessageID

	// Max ConsensusData Encoded
	biggestBeaconObject := [MaxConsensusDataEncoded]byte{1}
	fullData := biggestBeaconObject[:]

	fmt.Printf("FullData length: %v\n", len(fullData))

	// Justification with prepare messages
	prepareMessages := make([]*types.SignedSSVMessage, 0)
	for i := uint64(1); i <= 13; i++ {

		// qbft.Message
		prepareQbftMsg := &qbft.Message{
			MsgType:                  qbft.PrepareMsgType,
			Height:                   qbft.FirstHeight,
			Round:                    1,
			Identifier:               msgID[:],
			Root:                     sha256.Sum256(fullData),
			RoundChangeJustification: make([][]byte, 0),
			PrepareJustification:     make([][]byte, 0),
		}

		prepareQbftMsgBytes, err := prepareQbftMsg.Encode()
		require.NoError(t, err)
		if i == 13 {
			fmt.Printf("Prepare qbft.Message encoded length: %v\n", len(prepareQbftMsgBytes))
		}

		// SSVMessage
		prepareSSVMessage := &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   msgID,
			Data:    prepareQbftMsgBytes,
		}

		signature, err := SignSSVMessage(ks.OperatorKeys[i], prepareSSVMessage)
		require.NoError(t, err)

		// SignedSSVMessage
		prepareSignedSSVMessage := &types.SignedSSVMessage{
			Signatures:  [][]byte{signature},
			OperatorIDs: []types.OperatorID{i},
			SSVMessage:  prepareSSVMessage,
		}

		prepareSignedSSVMessageByts, _ := prepareSignedSSVMessage.Encode()
		if i == 13 {
			fmt.Printf("Prepare types.SignedSSVMessage encoded length: %v\n", len(prepareSignedSSVMessageByts))
		}

		prepareMessages = append(prepareMessages, prepareSignedSSVMessage)
	}

	prepareMessagesBytes := MarshalJustifications(prepareMessages)
	fmt.Printf("Prepare justification encoded length: [%v][%v] -> %v\n", len(prepareMessagesBytes), len(prepareMessagesBytes[0]), len(prepareMessagesBytes)*len(prepareMessagesBytes[0]))

	// qbft.Message
	rcQbftMsg := &qbft.Message{
		MsgType:                  qbft.RoundChangeMsgType,
		Height:                   qbft.FirstHeight,
		Round:                    2,
		Identifier:               msgID[:],
		Root:                     sha256.Sum256(fullData),
		DataRound:                1,
		RoundChangeJustification: prepareMessagesBytes,
		PrepareJustification:     make([][]byte, 0),
	}

	rcQbftMsgBytes, err := rcQbftMsg.Encode()
	require.NoError(t, err)
	fmt.Printf("Round-Change qbft.Message encoded length: %v\n", len(rcQbftMsgBytes))

	// SSVMessage
	rcSSVMessage := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   msgID,
		Data:    rcQbftMsgBytes,
	}

	signature, err := SignSSVMessage(ks.OperatorKeys[1], rcSSVMessage)
	require.NoError(t, err)

	// SignedSSVMessage
	rcSignedSSVMessage := &types.SignedSSVMessage{
		Signatures:  [][]byte{signature},
		OperatorIDs: []types.OperatorID{1},
		SSVMessage:  rcSSVMessage,
		FullData:    fullData,
	}

	rcSignedSSVMessageByts, _ := rcSignedSSVMessage.Encode()
	fmt.Printf("Round-Change types.SignedSSVMessage with Full-Data encoded length: %v\n", len(rcSignedSSVMessageByts))
}

func TestMaxRoundChangeMessageNoFullData(t *testing.T) {

	ks := Testing13SharesSet()
	msgID := TestingMessageID

	// Max ConsensusData Encoded
	biggestBeaconObject := [MaxConsensusDataEncoded]byte{1}
	fullData := biggestBeaconObject[:]

	fmt.Printf("FullData length: %v\n", len(fullData))

	// Justification with Prepare messages
	prepareMessages := make([]*types.SignedSSVMessage, 0)
	for i := uint64(1); i <= 13; i++ {

		// qbft.Message
		prepareQbftMsg := &qbft.Message{
			MsgType:                  qbft.PrepareMsgType,
			Height:                   qbft.FirstHeight,
			Round:                    1,
			Identifier:               msgID[:],
			Root:                     sha256.Sum256(fullData),
			RoundChangeJustification: make([][]byte, 0),
			PrepareJustification:     make([][]byte, 0),
		}

		prepareQbftMsgBytes, err := prepareQbftMsg.Encode()
		require.NoError(t, err)
		if i == 13 {
			fmt.Printf("Prepare qbft.Message encoded length: %v\n", len(prepareQbftMsgBytes))
		}

		// SSVMessage
		prepareSSVMessage := &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   msgID,
			Data:    prepareQbftMsgBytes,
		}

		signature, err := SignSSVMessage(ks.OperatorKeys[i], prepareSSVMessage)
		require.NoError(t, err)

		// SignedSSVMessage
		prepareSignedSSVMessage := &types.SignedSSVMessage{
			Signatures:  [][]byte{signature},
			OperatorIDs: []types.OperatorID{i},
			SSVMessage:  prepareSSVMessage,
		}

		prepareSignedSSVMessageByts, _ := prepareSignedSSVMessage.Encode()
		if i == 13 {
			fmt.Printf("Prepare types.SignedSSVMessage encoded length: %v\n", len(prepareSignedSSVMessageByts))
		}

		prepareMessages = append(prepareMessages, prepareSignedSSVMessage)
	}

	prepareMessagesBytes := MarshalJustifications(prepareMessages)
	fmt.Printf("Prepare justification encoded length: [%v][%v] -> %v\n", len(prepareMessagesBytes), len(prepareMessagesBytes[0]), len(prepareMessagesBytes)*len(prepareMessagesBytes[0]))

	// qbft.Message
	rcQbftMsg := &qbft.Message{
		MsgType:                  qbft.RoundChangeMsgType,
		Height:                   qbft.FirstHeight,
		Round:                    2,
		Identifier:               msgID[:],
		Root:                     sha256.Sum256(fullData),
		DataRound:                1,
		RoundChangeJustification: prepareMessagesBytes,
		PrepareJustification:     make([][]byte, 0),
	}

	rcQbftMsgBytes, err := rcQbftMsg.Encode()
	require.NoError(t, err)
	fmt.Printf("Round-Change qbft.Message encoded length: %v\n", len(rcQbftMsgBytes))

	// SSVMessage
	rcSSVMessage := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   msgID,
		Data:    rcQbftMsgBytes,
	}

	signature, err := SignSSVMessage(ks.OperatorKeys[1], rcSSVMessage)
	require.NoError(t, err)

	// SignedSSVMessage
	rcSignedSSVMessage := &types.SignedSSVMessage{
		Signatures:  [][]byte{signature},
		OperatorIDs: []types.OperatorID{1},
		SSVMessage:  rcSSVMessage,
	}

	rcSignedSSVMessageByts, _ := rcSignedSSVMessage.Encode()
	fmt.Printf("Round-Change types.SignedSSVMessage without Full-Data encoded length: %v\n", len(rcSignedSSVMessageByts))

}

func TestMaxProposalMessage(t *testing.T) {

	ks := Testing13SharesSet()
	msgID := TestingMessageID

	// Max ConsensusData Encoded
	biggestBeaconObject := [MaxConsensusDataEncoded]byte{1}
	fullData := biggestBeaconObject[:]

	fmt.Printf("FullData length: %v\n", len(fullData))

	// Prepare Justification
	prepareMessages := make([]*types.SignedSSVMessage, 0)
	for i := uint64(1); i <= 13; i++ {
		// qbft.Message
		prepareQbftMsg := &qbft.Message{
			MsgType:                  qbft.PrepareMsgType,
			Height:                   qbft.FirstHeight,
			Round:                    1,
			Identifier:               msgID[:],
			Root:                     sha256.Sum256(fullData),
			RoundChangeJustification: make([][]byte, 0),
			PrepareJustification:     make([][]byte, 0),
		}

		prepareQbftMsgBytes, err := prepareQbftMsg.Encode()
		require.NoError(t, err)
		if i == 13 {
			fmt.Printf("Prepare qbft.Message encoded length: %v\n", len(prepareQbftMsgBytes))
		}

		// SSVMessage
		prepareSSVMessage := &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   msgID,
			Data:    prepareQbftMsgBytes,
		}

		signature, err := SignSSVMessage(ks.OperatorKeys[i], prepareSSVMessage)
		require.NoError(t, err)

		// SignedSSVMessage
		prepareSignedSSVMessage := &types.SignedSSVMessage{
			Signatures:  [][]byte{signature},
			OperatorIDs: []types.OperatorID{i},
			SSVMessage:  prepareSSVMessage,
		}

		prepareSignedSSVMessageByts, _ := prepareSignedSSVMessage.Encode()
		if i == 13 {
			fmt.Printf("Prepare types.SignedSSVMessage encoded length: %v\n", len(prepareSignedSSVMessageByts))
		}

		prepareMessages = append(prepareMessages, prepareSignedSSVMessage)
	}

	prepareMessagesBytes := MarshalJustifications(prepareMessages)
	fmt.Printf("Prepare justification encoded length: [%v][%v] -> %v\n", len(prepareMessagesBytes), len(prepareMessagesBytes[0]), len(prepareMessagesBytes)*len(prepareMessagesBytes[0]))

	// Round-Change Justification
	rcMessages := make([]*types.SignedSSVMessage, 0)
	for i := uint64(1); i <= 13; i++ {

		// qbft.Message
		rcQbftMsg := &qbft.Message{
			MsgType:                  qbft.RoundChangeMsgType,
			Height:                   qbft.FirstHeight,
			Round:                    2,
			Identifier:               msgID[:],
			Root:                     sha256.Sum256(fullData),
			DataRound:                1,
			RoundChangeJustification: prepareMessagesBytes,
			PrepareJustification:     make([][]byte, 0),
		}

		rcQbftMsgBytes, err := rcQbftMsg.Encode()
		require.NoError(t, err)
		if i == 13 {
			fmt.Printf("Round-Change qbft.Message encoded length: %v\n", len(rcQbftMsgBytes))
		}

		// SSVMessage
		rcSSVMessage := &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   msgID,
			Data:    rcQbftMsgBytes,
		}

		signature, err := SignSSVMessage(ks.OperatorKeys[i], rcSSVMessage)
		require.NoError(t, err)

		// SignedSSVMessage without FullData
		rcSignedSSVMessage := &types.SignedSSVMessage{
			Signatures:  [][]byte{signature},
			OperatorIDs: []types.OperatorID{i},
			SSVMessage:  rcSSVMessage,
			FullData:    make([]byte, 0),
		}

		rcSignedSSVMessageByts, _ := rcSignedSSVMessage.Encode()
		if i == 13 {
			fmt.Printf("Round-Change types.SignedSSVMessage with no Full-Data encoded length: %v\n", len(rcSignedSSVMessageByts))
		}

		// SignedSSVMessage with FullData
		rcSignedSSVMessage = &types.SignedSSVMessage{
			Signatures:  [][]byte{signature},
			OperatorIDs: []types.OperatorID{i},
			SSVMessage:  rcSSVMessage,
			FullData:    fullData,
		}

		rcSignedSSVMessageByts, _ = rcSignedSSVMessage.Encode()
		if i == 13 {
			fmt.Printf("Round-Change types.SignedSSVMessage with Full-Data encoded length: %v\n", len(rcSignedSSVMessageByts))
		}

		rcMessages = append(rcMessages, rcSignedSSVMessage)
	}

	// Marshal and remove Full-Data fields
	rcMessagesBytes := MarshalJustifications(rcMessages)
	fmt.Printf("Round-Change justification encoded length: [%v][%v] -> %v\n", len(rcMessagesBytes), len(rcMessagesBytes[0]), len(rcMessagesBytes)*len(rcMessagesBytes[0]))

	// qbft.Message
	proposalQbftMessage := &qbft.Message{
		MsgType:                  qbft.ProposalMsgType,
		Height:                   qbft.FirstHeight,
		Round:                    2,
		Identifier:               msgID[:],
		Root:                     sha256.Sum256(fullData),
		DataRound:                1,
		RoundChangeJustification: rcMessagesBytes,
		PrepareJustification:     prepareMessagesBytes,
	}

	proposalQbftMsgBytes, err := proposalQbftMessage.Encode()
	require.NoError(t, err)
	fmt.Printf("Proposal qbft.Message encoded length: %v\n", len(proposalQbftMsgBytes))

	// SSVMessage
	proposalSSVMessage := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   msgID,
		Data:    proposalQbftMsgBytes,
	}

	// SignedSSVMessage
	signature, err := SignSSVMessage(ks.OperatorKeys[1], proposalSSVMessage)
	require.NoError(t, err)

	proposalSignedSSVMessage := &types.SignedSSVMessage{
		Signatures:  [][]byte{signature},
		OperatorIDs: []types.OperatorID{1},
		SSVMessage:  proposalSSVMessage,
		FullData:    fullData,
	}

	proposalSignedSSVMessageBytes, err := proposalSignedSSVMessage.Encode()
	require.NoError(t, err)
	fmt.Printf("Proposal types.SignedSSVMessage encoded length: %v\n", len(proposalSignedSSVMessageBytes))
}

func TestMaxPartialSig(t *testing.T) {

	msgID := TestingMessageID
	ks := Testing4SharesSet()

	// 1000 PartialSignatureMessage
	pSigMsgList := make([]*types.PartialSignatureMessage, 0)
	for i := uint64(1); i <= 1000; i++ {

		signingRoot := [32]byte{1, 2, 3, 4}

		strByt := [96]byte{1, 2, 3, 4}

		pSigMsg := types.PartialSignatureMessage{
			PartialSignature: strByt[:],
			SigningRoot:      signingRoot,
			Signer:           1,
			ValidatorIndex:   10,
		}

		pSigMsgList = append(pSigMsgList, &pSigMsg)

		if i == 13 {
			pSigByts, err := pSigMsg.Encode()
			require.NoError(t, err)
			fmt.Printf("PartialSignatureMessage encoded length: %v\n", len(pSigByts))
		}
	}

	// PartialSignatureMessages
	pSigMessages := types.PartialSignatureMessages{
		Type:     types.RandaoPartialSig,
		Slot:     12,
		Messages: pSigMsgList,
	}
	pSigMsgsByts, err := pSigMessages.Encode()
	require.NoError(t, err)
	fmt.Printf("PartialSignatureMessages encoded length: %v\n", len(pSigMsgsByts))

	// SSVMessage
	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   msgID,
		Data:    pSigMsgsByts,
	}

	signature, err := SignSSVMessage(ks.OperatorKeys[1], ssvMsg)
	require.NoError(t, err)

	// SignedSSVMessage
	signedSSVMessage := &types.SignedSSVMessage{
		Signatures:  [][]byte{signature},
		OperatorIDs: []types.OperatorID{1},
		SSVMessage:  ssvMsg,
	}

	signedSSVMessageByts, err := signedSSVMessage.Encode()
	require.NoError(t, err)
	fmt.Printf("PartialSignatureMessage types.SignedSSVMessage encoded length: %v\n", len(signedSSVMessageByts))

}
