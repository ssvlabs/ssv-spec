package partialsigcontainer

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/stretchr/testify/require"
)

type PartialSigContainerTest struct {
	Name            string
	Quorum          uint64
	ValidatorPubKey []byte
	SignatureMsgs   []*types.PartialSignatureMessage
	ExpectedError   string
	ExpectedResult  []byte
	ExpectedQuorum  bool
}

func (test *PartialSigContainerTest) TestName() string {
	return "PartialSigContainer " + test.Name
}

func (test *PartialSigContainerTest) Run(t *testing.T) {
	ps := ssv.NewPartialSigContainer(test.Quorum)

	validatorIndexRoots := make(map[phase0.ValidatorIndex][32]byte)
	// Add signature messages
	for _, sigMsg := range test.SignatureMsgs {
		ps.AddSignature(sigMsg)
		validatorIndexRoots[sigMsg.ValidatorIndex] = sigMsg.SigningRoot
	}

	for validatorIndex, root := range validatorIndexRoots {

		// Check quorum
		require.Equal(t, test.ExpectedQuorum, ps.HasQuorum(validatorIndex, root))

		result, err := ps.ReconstructSignature(root, test.ValidatorPubKey, validatorIndex)
		// Check the result and error
		if len(test.ExpectedError) > 0 {
			require.Error(t, err)
			require.Contains(t, err.Error(), test.ExpectedError)
		} else {
			require.NoError(t, err)
			require.EqualValues(t, test.ExpectedResult, result)
		}
	}
}

func (test *PartialSigContainerTest) GetPostState() (interface{}, error) {
	return nil, nil
}
