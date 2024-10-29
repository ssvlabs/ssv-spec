package validation

import (
	"crypto/rsa"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// HappyFlow tests a full happy flow for a duty and confirms that there's no error
func HappyFlow() tests.SpecTest {

	ks := testingutils.DefaultKeySet

	// PreConsensus messages
	preConsensusMsgs := []*types.PartialSignatureMessages{
		testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
		testingutils.PreConsensusRandaoMsg(ks.Shares[2], 2),
		testingutils.PreConsensusRandaoMsg(ks.Shares[3], 3),
		testingutils.PreConsensusRandaoMsg(ks.Shares[4], 4),
	}

	// Consensus messages
	consensusMsgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithID(ks.OperatorKeys[1], 1, testingutils.DefaultMsgID),
		testingutils.TestingPrepareMessageWithIdentifierAndRoot(ks.OperatorKeys[1], 1, testingutils.DefaultMsgID[:], testingutils.TestingQBFTRootData),
		testingutils.TestingPrepareMessageWithIdentifierAndRoot(ks.OperatorKeys[2], 2, testingutils.DefaultMsgID[:], testingutils.TestingQBFTRootData),
		testingutils.TestingPrepareMessageWithIdentifierAndRoot(ks.OperatorKeys[3], 3, testingutils.DefaultMsgID[:], testingutils.TestingQBFTRootData),
		testingutils.TestingPrepareMessageWithIdentifierAndRoot(ks.OperatorKeys[4], 4, testingutils.DefaultMsgID[:], testingutils.TestingQBFTRootData),
		testingutils.TestingCommitMessageWithIdentifierAndRoot(ks.OperatorKeys[1], 1, testingutils.DefaultMsgID[:], testingutils.TestingQBFTRootData),
		testingutils.TestingCommitMessageWithIdentifierAndRoot(ks.OperatorKeys[2], 2, testingutils.DefaultMsgID[:], testingutils.TestingQBFTRootData),
		testingutils.TestingCommitMessageWithIdentifierAndRoot(ks.OperatorKeys[3], 3, testingutils.DefaultMsgID[:], testingutils.TestingQBFTRootData),
		testingutils.TestingCommitMessageWithIdentifierAndRoot(ks.OperatorKeys[4], 4, testingutils.DefaultMsgID[:], testingutils.TestingQBFTRootData),
		testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
			[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
			[]types.OperatorID{1, 2, 3},
			qbft.FirstHeight, testingutils.DefaultMsgID[:], testingutils.TestingQBFTFullData),
		testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
			[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[4]},
			[]types.OperatorID{1, 2, 4},
			qbft.FirstHeight, testingutils.DefaultMsgID[:], testingutils.TestingQBFTFullData),
		testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
			[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[3], ks.OperatorKeys[4]},
			[]types.OperatorID{1, 3, 4},
			qbft.FirstHeight, testingutils.DefaultMsgID[:], testingutils.TestingQBFTFullData),
		testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
			[]*rsa.PrivateKey{ks.OperatorKeys[2], ks.OperatorKeys[3], ks.OperatorKeys[4]},
			[]types.OperatorID{2, 3, 4},
			qbft.FirstHeight, testingutils.DefaultMsgID[:], testingutils.TestingQBFTFullData),
		testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
			[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3], ks.OperatorKeys[4]},
			[]types.OperatorID{1, 2, 3, 4},
			qbft.FirstHeight, testingutils.DefaultMsgID[:], testingutils.TestingQBFTFullData),
	}

	// PostConsensus messages
	postConsensusMsgs := []*types.PartialSignatureMessages{
		testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
		testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, spec.DataVersionDeneb),
		testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, spec.DataVersionDeneb),
		testingutils.PostConsensusProposerMsgV(ks.Shares[4], 4, spec.DataVersionDeneb),
	}

	// Function to transform a partial signature message to a SignedSSVMessage
	getSignedSSVMessageForPartialSignatureMessage := func(pSigMsgs *types.PartialSignatureMessages) *types.SignedSSVMessage {
		pSigMsgs.Slot = phase0.Slot(qbft.FirstHeight)
		ssvMsg := &types.SSVMessage{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   testingutils.DefaultMsgID,
			Data:    testingutils.EncodePartialSignatureMessage(pSigMsgs),
		}
		signer := pSigMsgs.Messages[0].Signer
		return &types.SignedSSVMessage{
			OperatorIDs: []types.OperatorID{signer},
			Signatures:  [][]byte{testingutils.SignSSVMessage(ks.OperatorKeys[signer], ssvMsg)},
			SSVMessage:  ssvMsg,
		}
	}

	// Encode messages
	msgs := [][]byte{}
	for _, msg := range preConsensusMsgs {
		msgs = append(msgs, testingutils.EncodeMessage(getSignedSSVMessageForPartialSignatureMessage(msg)))
	}
	for _, msg := range consensusMsgs {
		msgs = append(msgs, testingutils.EncodeMessage(msg))
	}
	for _, msg := range postConsensusMsgs {
		msgs = append(msgs, testingutils.EncodeMessage(getSignedSSVMessageForPartialSignatureMessage(msg)))
	}

	return &MessageValidationTest{
		Name:     "happy flow",
		Messages: msgs,
	}
}
