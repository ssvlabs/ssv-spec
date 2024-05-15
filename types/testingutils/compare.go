package testingutils

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/stretchr/testify/require"
)

func filterPartialSigs(messages []*types.SSVMessage) []*types.SSVMessage {
	ret := make([]*types.SSVMessage, 0)
	for _, msg := range messages {
		if msg.MsgType != types.SSVPartialSignatureMsgType {
			continue
		}
		ret = append(ret, msg)
	}
	return ret
}

func ComparePartialSignatureOutputMessages(t *testing.T, expectedMessages []*types.PartialSignatureMessages, broadcastedSignedMsgs []*types.SignedSSVMessage, committee []*types.CommitteeMember) {

	require.NoError(t, VerifyListOfSignedSSVMessages(broadcastedSignedMsgs, committee))

	broadcastedMsgs := ConvertBroadcastedMessagesToSSVMessages(broadcastedSignedMsgs)

	broadcastedMsgs = filterPartialSigs(broadcastedMsgs)
	require.Len(t, broadcastedMsgs, len(expectedMessages))

	index := 0
	for _, msg := range broadcastedMsgs {

		msg1 := &types.PartialSignatureMessages{}
		require.NoError(t, msg1.Decode(msg.Data))

		msg2 := expectedMessages[index]

		ComparePartialSignatureMessages(t, msg1, msg2)

		index++
	}
}

func ComparePartialSignatureMessages(t *testing.T, msg1 *types.PartialSignatureMessages, msg2 *types.PartialSignatureMessages) {

	require.Len(t, msg1.Messages, len(msg2.Messages))

	// messages are not guaranteed to be in order so we map them and then test all roots to be equal
	roots := make(map[string]string)
	for i, partialSigMsg2 := range msg2.Messages {
		r2, err := partialSigMsg2.GetRoot()
		require.NoError(t, err)
		if _, found := roots[hex.EncodeToString(r2[:])]; !found {
			roots[hex.EncodeToString(r2[:])] = ""
		} else {
			roots[hex.EncodeToString(r2[:])] = hex.EncodeToString(r2[:])
		}

		partialSigMsg1 := msg1.Messages[i]
		r1, err := partialSigMsg1.GetRoot()
		require.NoError(t, err)

		if _, found := roots[hex.EncodeToString(r1[:])]; !found {
			roots[hex.EncodeToString(r1[:])] = ""
		} else {
			roots[hex.EncodeToString(r1[:])] = hex.EncodeToString(r1[:])
		}
	}
	for k, v := range roots {
		require.EqualValues(t, k, v, "missing output msg")
	}

	// test that slot is correct in broadcasted msg
	require.EqualValues(t, msg1.Slot, msg2.Slot, "incorrect broadcasted slot")
}

func CompareSignedSSVMessageOutputMessages(t *testing.T, expectedMessages []*types.SignedSSVMessage, broadcastedSignedMsgs []*types.SignedSSVMessage, committee []*types.CommitteeMember) {

	require.NoError(t, VerifyListOfSignedSSVMessages(broadcastedSignedMsgs, committee))

	require.Len(t, broadcastedSignedMsgs, len(expectedMessages))

	for index, msg := range broadcastedSignedMsgs {
		r1, _ := msg.GetRoot()

		msg2 := broadcastedSignedMsgs[index]
		r2, _ := msg2.GetRoot()

		require.EqualValues(t, r1, r2, fmt.Sprintf("output msg %d roots not equal", index))
	}
}

func CompareBroadcastedBeaconMsgs(t *testing.T, expectedRoots []string, broadcastedRoots []phase0.Root) {
	require.Len(t, broadcastedRoots, len(expectedRoots))
	for _, r1 := range expectedRoots {
		found := false
		for _, r2 := range broadcastedRoots {
			if r1 == hex.EncodeToString(r2[:]) {
				found = true
				break
			}
		}
		require.Truef(t, found, "broadcasted beacon root not found")
	}
}
