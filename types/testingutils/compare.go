package testingutils

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
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

// ComparePartialSignatureOutputMessagesStrictOrder compares output messages with strict ordering
// Both the order of messages and the order of partial signatures within each message must match
func ComparePartialSignatureOutputMessagesStrictOrder(t *testing.T, expectedMessages []*types.PartialSignatureMessages, broadcastedSignedMsgs []*types.SignedSSVMessage, committee []*types.Operator) {

	require.NoError(t, VerifyListOfSignedSSVMessages(broadcastedSignedMsgs, committee))

	broadcastedMsgs := ConvertBroadcastedMessagesToSSVMessages(broadcastedSignedMsgs)
	broadcastedMsgs = filterPartialSigs(broadcastedMsgs)

	require.Len(t, broadcastedMsgs, len(expectedMessages))

	for index, msg := range broadcastedMsgs {
		msg1 := &types.PartialSignatureMessages{}
		require.NoError(t, msg1.Decode(msg.Data))

		msg2 := expectedMessages[index]

		err := ComparePartialSignatureMessagesStrict(msg1, msg2)
		require.NoError(t, err, "message %d comparison failed", index)
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
	return ComparePartialSignatureMessagesWithOptions(msg1, msg2, false)
}

// ComparePartialSignatureMessagesStrict compares two PartialSignatureMessages with strict ordering
func ComparePartialSignatureMessagesStrict(msg1 *types.PartialSignatureMessages, msg2 *types.PartialSignatureMessages) error {
	return ComparePartialSignatureMessagesWithOptions(msg1, msg2, true)
}

// ComparePartialSignatureMessagesWithOptions compares two PartialSignatureMessages with optional strict ordering
func ComparePartialSignatureMessagesWithOptions(msg1 *types.PartialSignatureMessages, msg2 *types.PartialSignatureMessages, strictOrder bool) error {

	if len(msg1.Messages) != len(msg2.Messages) {
		return errors.New("different messages length")
	}

	if strictOrder {
		// Compare messages in order - each message must match exactly at the same index
		for i := range msg1.Messages {
			m1 := msg1.Messages[i]
			m2 := msg2.Messages[i]

			if m1.ValidatorIndex != m2.ValidatorIndex {
				return errors.Errorf("message %d: validator index mismatch: got %d, expected %d", i, m1.ValidatorIndex, m2.ValidatorIndex)
			}
			if m1.SigningRoot != m2.SigningRoot {
				return errors.Errorf("message %d: signing root mismatch", i)
			}
			// Compare partial signature
			r1, err := m1.GetRoot()
			if err != nil {
				return errors.Wrap(err, "failed to get root for message 1")
			}
			r2, err := m2.GetRoot()
			if err != nil {
				return errors.Wrap(err, "failed to get root for message 2")
			}
			if r1 != r2 {
				return errors.Errorf("message %d: root mismatch", i)
			}
		}
	} else {
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

func CompareConsensusData(t *testing.T, expectedData [][]byte, actualData [][]byte) {
	require.Len(t, actualData, len(expectedData))
	// expectedData should be compared to actualData in order
	for i := range expectedData {
		compareConsensusDataSample(t, expectedData[i], actualData[i])
	}
}

// compareConsensusDataSample compares two consensus data samples.
// It tries decoding into ProposerConsensusData, AggregatorCommitteeConsensusData, or BeaconVote
// and compare each type accordingly.
func compareConsensusDataSample(t *testing.T, expectedData []byte, actualData []byte) {
	expectedBeaconVote := &types.BeaconVote{}
	if err := expectedBeaconVote.Decode(expectedData); err == nil {
		actualBeaconVote := &types.BeaconVote{}
		require.NoError(t, actualBeaconVote.Decode(actualData))
		compareBeaconVotes(t, expectedBeaconVote, actualBeaconVote)
		return
	}

	expectedProposerData := &types.ProposerConsensusData{}
	if err := expectedProposerData.Decode(expectedData); err == nil {
		actualProposerData := &types.ProposerConsensusData{}
		require.NoError(t, actualProposerData.Decode(actualData))
		compareProposerConsensusData(t, expectedProposerData, actualProposerData)
		return
	}

	expectedAggCommData := &types.AggregatorCommitteeConsensusData{}
	if err := expectedAggCommData.Decode(expectedData); err == nil {
		actualAggCommData := &types.AggregatorCommitteeConsensusData{}
		require.NoError(t, actualAggCommData.Decode(actualData))
		compareAggregatorCommitteeConsensusData(t, expectedAggCommData, actualAggCommData)
		return
	}

	require.Fail(t, "could not decode consensus data as any known type")
}

// Compare beacon votes
func compareBeaconVotes(t *testing.T, expectedVotes *types.BeaconVote, actualVotes *types.BeaconVote) {
	require.Equal(t, expectedVotes.BlockRoot, actualVotes.BlockRoot, "beacon vote block root mismatch")
	require.Equal(t, expectedVotes.Source.Epoch, actualVotes.Source.Epoch, "beacon vote source epoch mismatch")
	require.Equal(t, expectedVotes.Source.Root, actualVotes.Source.Root, "beacon vote source root mismatch")
	require.Equal(t, expectedVotes.Target.Epoch, actualVotes.Target.Epoch, "beacon vote target epoch mismatch")
	require.Equal(t, expectedVotes.Target.Root, actualVotes.Target.Root, "beacon vote target root mismatch")
}

// Compare proposer consensus data
func compareProposerConsensusData(t *testing.T, expectedData *types.ProposerConsensusData, actualData *types.ProposerConsensusData) {
	require.Equal(t, expectedData.DataSSZ, actualData.DataSSZ, "proposer consensus data mismatch")
	require.Equal(t, expectedData.Version, actualData.Version, "proposer consensus data version mismatch")
	compareValidatorDuty(t, &expectedData.Duty, &actualData.Duty)
}

// Compare duty
func compareValidatorDuty(t *testing.T, expectedDuty *types.ValidatorDuty, actualDuty *types.ValidatorDuty) {
	require.Equal(t, expectedDuty.Type, actualDuty.Type, "validator duty type mismatch")
	require.Equal(t, expectedDuty.PubKey, actualDuty.PubKey, "validator duty pubkey mismatch")
	require.Equal(t, expectedDuty.Slot, actualDuty.Slot, "validator duty slot mismatch")
	require.Equal(t, expectedDuty.ValidatorIndex, actualDuty.ValidatorIndex, "validator duty validator index mismatch")
	require.Equal(t, expectedDuty.CommitteeIndex, actualDuty.CommitteeIndex, "validator duty committee index mismatch")
	require.Equal(t, expectedDuty.CommitteeLength, actualDuty.CommitteeLength, "validator duty committee length mismatch")
	require.Equal(t, expectedDuty.CommitteesAtSlot, actualDuty.CommitteesAtSlot, "validator duty committees at slot mismatch")
	require.Equal(t, expectedDuty.ValidatorCommitteeIndex, actualDuty.ValidatorCommitteeIndex, "validator duty validator committee index mismatch")
	require.Equal(t, expectedDuty.ValidatorSyncCommitteeIndices, actualDuty.ValidatorSyncCommitteeIndices, "validator duty validator sync committee indices mismatch")
}

// Compare aggregator committee consensus data
func compareAggregatorCommitteeConsensusData(t *testing.T, expectedData *types.AggregatorCommitteeConsensusData, actualData *types.AggregatorCommitteeConsensusData) {
	// Compare version
	require.Equal(t, expectedData.Version, actualData.Version, "aggregator committee consensus data version mismatch")

	// Compare lengths
	require.Equal(t, len(expectedData.Aggregators), len(actualData.Aggregators), "aggregator committee consensus data aggregators length mismatch")
	require.Equal(t, len(expectedData.AggregatorsCommitteeIndexes), len(actualData.AggregatorsCommitteeIndexes), "aggregator committee consensus data aggregators committee indexes length mismatch")
	require.Equal(t, len(expectedData.AggregatedAttestations), len(actualData.AggregatedAttestations), "aggregator committee consensus data aggregated attestations length mismatch")
	require.Equal(t, len(expectedData.Contributors), len(actualData.Contributors), "aggregator committee consensus data contributors length mismatch")
	require.Equal(t, len(expectedData.SyncCommitteeContributions), len(actualData.SyncCommitteeContributions), "aggregator committee consensus data sync committee contributions length mismatch")

	// Compare each Aggregator, without requiring the same order
	for _, expectedAgg := range expectedData.Aggregators {
		found := false
		for _, actualAgg := range actualData.Aggregators {
			if isEqualAssignedAggregator(&expectedAgg, &actualAgg) {
				found = true
				break
			}
		}
		require.Truef(t, found, "aggregator not found in actual consensus data: validator index %d, committee index %d", expectedAgg.ValidatorIndex, expectedAgg.CommitteeIndex)
	}

	// Compare each AggregatedAttestation, without requiring the same order
	for _, expectedAtt := range expectedData.AggregatedAttestations {
		found := false
		for _, actualAtt := range actualData.AggregatedAttestations {
			if isEqualAggregatedAttestation(expectedAtt, actualAtt) {
				found = true
				break
			}
		}
		require.Truef(t, found, "aggregated attestation not found in actual consensus data: %s", hex.EncodeToString(expectedAtt))
	}

	// Compare each SyncCommitteeContribution, without requiring the same order
	for _, expectedSCC := range expectedData.SyncCommitteeContributions {
		found := false
		for _, actualSCC := range actualData.SyncCommitteeContributions {
			if isEqualSyncCommitteeContribution(&expectedSCC, &actualSCC) {
				found = true
				break
			}

		}
		require.Truef(t, found, "sync committee contribution not found in actual consensus data: subcommittee index %d", expectedSCC.SubcommitteeIndex)
	}

	// Compare each Contributor, without requiring the same order
	for _, expectedContributor := range expectedData.Contributors {
		found := false
		for _, actualContributor := range actualData.Contributors {
			if isEqualAssignedAggregator(&expectedContributor, &actualContributor) {
				found = true
				break
			}
		}
		require.Truef(t, found, "contributor not found in actual consensus data: validator index %d, committee index %d", expectedContributor.ValidatorIndex, expectedContributor.CommitteeIndex)
	}
}

// Is equal assigned aggregator
func isEqualAssignedAggregator(expected *types.AssignedAggregator, actual *types.AssignedAggregator) bool {
	return expected.ValidatorIndex == actual.ValidatorIndex &&
		expected.CommitteeIndex == actual.CommitteeIndex &&
		cmp.Equal(expected.SelectionProof, actual.SelectionProof)
}

// Is equal sync committee contribution
func isEqualSyncCommitteeContribution(expected *altair.SyncCommitteeContribution, actual *altair.SyncCommitteeContribution) bool {
	return expected.Slot == actual.Slot &&
		cmp.Equal(expected.BeaconBlockRoot, actual.BeaconBlockRoot) &&
		expected.SubcommitteeIndex == actual.SubcommitteeIndex &&
		cmp.Equal(expected.AggregationBits, actual.AggregationBits) &&
		cmp.Equal(expected.Signature, actual.Signature)
}

// Is equal aggregated attestation
func isEqualAggregatedAttestation(expected []byte, actual []byte) bool {
	return string(expected) == string(actual)
}
