package partialsigcontainer

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func OneSignature() tests.SpecTest {

	// Create a test key set
	ks := testingutils.Testing4SharesSet()

	// Create PartialSignatureMessage for testing
	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msgs := []*types.PartialSignatureMessage{msg.Messages[0]}

	// Verify the reconstructed signature
	expectedSig, err := types.ReconstructSignatures(map[types.OperatorID][]byte{1: msgs[0].PartialSignature})
	if err != nil {
		panic(err.Error())
	}

	return &PartialSigContainerTest{
		Name:            "one signature",
		Quorum:          ks.Threshold,
		ValidatorPubKey: ks.ValidatorPK.Serialize(),
		SignatureMsgs:   msgs,
		ExpectedError:   "could not reconstruct a valid signature",
		ExpectedResult:  expectedSig.Serialize(),
		ExpectedQuorum:  false,
	}
}
