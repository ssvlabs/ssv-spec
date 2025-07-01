package partialsigcontainer

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
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
	PrivateKeys     *PrivateKeyInfo `json:"PrivateKeys,omitempty"`
}

type PrivateKeyInfo struct {
	ValidatorSK  string
	Shares       map[types.OperatorID]string
	OperatorKeys map[types.OperatorID]string
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

func (test *PartialSigContainerTest) SetPrivateKeys(ks *testingutils.TestKeySet) {
	privateKeyInfo := &PrivateKeyInfo{
		ValidatorSK:  hex.EncodeToString(ks.ValidatorSK.Serialize()),
		Shares:       make(map[types.OperatorID]string),
		OperatorKeys: make(map[types.OperatorID]string),
	}

	// Add share keys
	for operatorID, shareSK := range ks.Shares {
		privateKeyInfo.Shares[operatorID] = hex.EncodeToString(shareSK.Serialize())
	}

	// Add operator keys (RSA private keys used for signing)
	for operatorID, operatorKey := range ks.OperatorKeys {
		privateKeyInfo.OperatorKeys[operatorID] = fmt.Sprintf("N:%s,E:%d",
			operatorKey.N.String(), operatorKey.E)
	}

	test.PrivateKeys = privateKeyInfo
}
