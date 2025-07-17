package partialsigcontainer

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func DuplicateQuorum() tests.SpecTest {

	// Create a test key set
	ks := testingutils.Testing4SharesSet()

	// Create PartialSignatureMessage for testing
	msg1 := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, spec.DataVersionPhase0)
	msg12 := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, spec.DataVersionPhase0)
	msg2 := testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, spec.DataVersionPhase0)
	msg3 := testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, spec.DataVersionPhase0)
	msgs := []*types.PartialSignatureMessage{msg1.Messages[0], msg12.Messages[0], msg2.Messages[0], msg3.Messages[0]}

	// Verify the reconstructed signature
	expectedSig, err := types.ReconstructSignatures(map[types.OperatorID][]byte{1: msgs[0].PartialSignature, 2: msgs[2].PartialSignature, 3: msgs[3].PartialSignature})
	if err != nil {
		panic(err.Error())
	}

	return NewPartialSigContainerTest(
		"duplicate quorum",
		"Tests partial signature container with duplicate signatures but still achieving quorum",
		ks.Threshold,
		ks.ValidatorPK.Serialize(),
		msgs,
		"",
		expectedSig.Serialize(),
		true,
	)
}
