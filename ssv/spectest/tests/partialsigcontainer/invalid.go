package partialsigcontainer

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func Invalid() tests.SpecTest {

	// Create a test key set
	ks := testingutils.Testing4SharesSet()

	// Create PartialSignatureMessage for testing
	msg1 := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg2 := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 2, qbft.FirstHeight) // invalid signature
	msg3 := testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, qbft.FirstHeight)
	msgs := []*types.PartialSignatureMessage{msg1.Messages[0], msg2.Messages[0], msg3.Messages[0]}

	// Verify the reconstructed signature
	expectedSig, err := types.ReconstructSignatures(map[types.OperatorID][]byte{1: msgs[0].PartialSignature, 2: msgs[1].PartialSignature, 3: msgs[2].PartialSignature})
	if err != nil {
		panic(err.Error())
	}

	return &PartialSigContainerTest{
		Name:            "invalid",
		Quorum:          ks.Threshold,
		ValidatorPubKey: ks.ValidatorPK.Serialize(),
		SignatureMsgs:   msgs,
		ExpectedError:   "could not reconstruct a valid signature",
		ExpectedResult:  expectedSig.Serialize(),
		ExpectedQuorum:  true,
	}
}
