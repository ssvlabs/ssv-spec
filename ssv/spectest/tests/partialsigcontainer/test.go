package partialsigcontainer

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

type PartialSigContainerTest struct {
	Name              string
	Type              string
	Documentation     string
	Quorum            uint64
	ValidatorPubKey   []byte
	SignatureMsgs     []*types.PartialSignatureMessage
	ExpectedErrorCode int
	ExpectedResult    []byte
	ExpectedQuorum    bool
	PrivateKeys       *testingutils.PrivateKeyInfo `json:"PrivateKeys,omitempty"`
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
		tests.AssertErrorCode(t, test.ExpectedErrorCode, err)
		if err != nil {
			require.Contains(t, err.Error(), test.ExpectedErrorCode)
			return
		}
		require.EqualValues(t, test.ExpectedResult, result)
	}
}

func (test *PartialSigContainerTest) GetPostState() (interface{}, error) {
	return nil, nil
}

func NewPartialSigContainerTest(name, documentation string, quorum uint64, validatorPubKey []byte, signatureMsgs []*types.PartialSignatureMessage, expectedErrorCode int, expectedResult []byte, expectedQuorum bool, ks *testingutils.TestKeySet) *PartialSigContainerTest {
	return &PartialSigContainerTest{
		Name:              name,
		Type:              testdoc.PartialSigContainerTestType,
		Documentation:     documentation,
		Quorum:            quorum,
		ValidatorPubKey:   validatorPubKey,
		SignatureMsgs:     signatureMsgs,
		ExpectedErrorCode: expectedErrorCode,
		ExpectedResult:    expectedResult,
		ExpectedQuorum:    expectedQuorum,
		PrivateKeys:       testingutils.BuildPrivateKeyInfo(ks),
	}
}
