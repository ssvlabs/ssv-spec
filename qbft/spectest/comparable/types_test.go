package comparable

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	testingutilscomparable "github.com/bloxapp/ssv-spec/types/testingutils/comparable"
	"testing"
)

func TestSprintf(t *testing.T) {
	a := &qbft.State{
		ID: []byte{1, 2, 3, 4},
		Share: &types.Share{
			OperatorID: 1,
		},
	}
	b := &qbft.State{
		ID: []byte{1, 2, 3, 4},
		Share: &types.Share{
			OperatorID: 2,
		},
	}

	testingutilscomparable.Print(CompareQBFTState(a, b))
}
