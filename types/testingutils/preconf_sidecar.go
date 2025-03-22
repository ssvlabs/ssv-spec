package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/types"
)

var TestingPreconfRoot = phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2}

type TestingPreconfSidecar struct {
}

func NewTestingPreconfSidecar() *TestingPreconfSidecar {
	return &TestingPreconfSidecar{}
}

// GetNewRequest returns a new preconf request
func (sidecar *TestingPreconfSidecar) GetNewRequest() (types.CBSigningRequest, error) {
	request := types.CBSigningRequest{
		Root: TestingPreconfRoot,
	}
	return request, nil
}

// SubmitCommitment submits a commitment to the node
func (sidecar *TestingPreconfSidecar) SubmitCommitment(requestRoot phase0.Root, signature phase0.BLSSignature) error {
	return nil
}
