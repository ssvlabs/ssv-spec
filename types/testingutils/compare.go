package testingutils

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/stretchr/testify/require"
)

func filterPartialSigs(messages []*types.SSVMessage) []*types.SSVMessage {
	ret := make([]*types.SSVMessage, 0)
	for _, msg := range messages {
		if msg.MsgType == types.SSVPartialSignatureMsgType {
			ret = append(ret, msg)
		} else if msg.MsgType == types.CommitBoostPartialSignatureMsgType {
			CBMsg := &types.CBPartialSignatures{}
			CBMsg.Decode(msg.Data)
			partialSigMsg, err := CBMsg.PartialSig.Encode()
			if err != nil {
				panic(err)
			}
			ret = append(ret, &types.SSVMessage{
				MsgType: types.SSVPartialSignatureMsgType,
				MsgID:   msg.MsgID,
				Data:    partialSigMsg,
			})
		}
	}
	return ret
}

func ComparePartialSignatureOutputMessages(t *testing.T, expectedMessages []*types.PartialSignatureMessages, broadcastedSignedMsgs []*types.SignedSSVMessage, committee []*types.Operator) {

	require.NoError(t, VerifyListOfSignedSSVMessages(broadcastedSignedMsgs, committee))

	broadcastedMsgs := ConvertBroadcastedMessagesToSSVMessages(broadcastedSignedMsgs)

	broadcastedMsgs = filterPartialSigs(broadcastedMsgs)
	require.Len(t, broadcastedMsgs, len(expectedMessages))

	for index, msg := range broadcastedMsgs {

		msg1 := &types.PartialSignatureMessages{}
		require.NoError(t, msg1.Decode(msg.Data))

		msg2 := expectedMessages[index]

		err := ComparePartialSignatureMessages(msg1, msg2)
		require.NoError(t, err)
	}
}

// Compare partial sig output messages without assuming any order between messages (asynchonous)
func ComparePartialSignatureOutputMessagesInAsynchronousOrder(t *testing.T, expectedMessages []*types.PartialSignatureMessages, broadcastedSignedMsgs []*types.SignedSSVMessage, committee []*types.Operator) {

	require.NoError(t, VerifyListOfSignedSSVMessages(broadcastedSignedMsgs, committee))

	broadcastedMsgs := ConvertBroadcastedMessagesToSSVMessages(broadcastedSignedMsgs)
	broadcastedMsgs = filterPartialSigs(broadcastedMsgs)

	// Require that:
	// - the broadcasted and expected messages have equal length
	// - every broadcasted message is linked (equal) to an expected message
	// - two broadcasted messages are not linked to the same expected message
	// i.e. a bijection between the lists
	require.Len(t, broadcastedMsgs, len(expectedMessages))

	expectedMsgAlreadyLinked := make([]bool, len(expectedMessages))
	for i := range expectedMsgAlreadyLinked {
		expectedMsgAlreadyLinked[i] = false
	}
	for _, msg := range broadcastedMsgs {
		msg1 := &types.PartialSignatureMessages{}
		require.NoError(t, msg1.Decode(msg.Data))

		found := false
		for expectedMsgIndex, msg2 := range expectedMessages {
			if expectedMsgAlreadyLinked[expectedMsgIndex] {
				continue
			}
			err := ComparePartialSignatureMessages(msg1, msg2)
			if err == nil {
				found = true
				expectedMsgAlreadyLinked[expectedMsgIndex] = true
				break
			}
		}
		require.True(t, found)
	}

	// Assert that all expected messages are linked.
	// An expected message not linked should be an impossible state (i.e. an error should be triggered by the above checks)
	for _, linked := range expectedMsgAlreadyLinked {
		require.True(t, linked)
	}
}

func RootCountMapForPartialSignatureMessages(msg *types.PartialSignatureMessages) map[string]int {
	roots := make(map[string]int)

	for _, partialSigMessage := range msg.Messages {
		root, err := partialSigMessage.GetRoot()
		if err != nil {
			panic(err)
		}
		rootStr := hex.EncodeToString(root[:])
		if _, found := roots[rootStr]; !found {
			roots[rootStr] = 0
		}
		roots[rootStr] += 1
	}

	return roots
}

func ComparePartialSignatureMessages(msg1 *types.PartialSignatureMessages, msg2 *types.PartialSignatureMessages) error {

	if len(msg1.Messages) != len(msg2.Messages) {
		return errors.New("different messages length")
	}

	// messages are not guaranteed to be in order so we map their roots and then test all roots to match and have the same multiplicity
	roots1 := RootCountMapForPartialSignatureMessages(msg1)
	roots2 := RootCountMapForPartialSignatureMessages(msg2)

	// Compare roots and their multiplicity
	if len(roots1) != len(roots2) {
		return errors.New("messages have different sets of roots")
	}
	for r1, r1Count := range roots1 {
		foundSameRootAndSameCount := false
		for r2, r2Count := range roots2 {
			if r1 == r2 {
				foundSameRootAndSameCount = (r1Count == r2Count)
				break
			}
		}
		if !foundSameRootAndSameCount {
			return errors.New("missing output msg")
		}
	}

	// test that slot is correct in broadcasted msg
	if msg1.Slot != msg2.Slot {
		return errors.New("incorrect broadcasted slot")
	}
	// test that type is correct in broadcasted msg
	if msg1.Type != msg2.Type {
		return errors.New("incorrect broadcasted type")
	}
	return nil
}

func CompareSignedSSVMessageOutputMessages(t *testing.T, expectedMessages []*types.SignedSSVMessage, broadcastedSignedMsgs []*types.SignedSSVMessage, committee []*types.Operator) {

	require.NoError(t, VerifyListOfSignedSSVMessages(broadcastedSignedMsgs, committee))

	require.Len(t, broadcastedSignedMsgs, len(expectedMessages))

	for index, msg := range broadcastedSignedMsgs {
		r1, _ := msg.GetRoot()

		msg2 := expectedMessages[index]
		r2, _ := msg2.GetRoot()

		require.EqualValues(t, r1, r2, fmt.Sprintf("output msg %d roots not equal", index))
	}
}

func CompareBroadcastedBeaconMsgs(t *testing.T, expectedRoots []string, broadcastedRoots []phase0.Root) {
	require.Len(t, broadcastedRoots, len(expectedRoots))

	// broadcastedRootAlreadyLinked has to purpose of not using the same
	// broadcasted root twice when confirming that an expected root exists
	broadcastedRootAlreadyLinked := make([]bool, len(broadcastedRoots))
	for i := range broadcastedRootAlreadyLinked {
		broadcastedRootAlreadyLinked[i] = false
	}
	for _, r1 := range expectedRoots {
		found := false
		for index2, r2 := range broadcastedRoots {
			if broadcastedRootAlreadyLinked[index2] {
				continue
			}
			if r1 == hex.EncodeToString(r2[:]) {
				found = true
				broadcastedRootAlreadyLinked[index2] = true
				break
			}
		}
		require.Truef(t, found, "broadcasted beacon root not found")
	}
}
