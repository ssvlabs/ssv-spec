package frost

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_isResharing(t *testing.T) {
	t.Run("keygen", func(t *testing.T) {
		params := &InstanceParams{}

		isResharing := params.isResharing()
		require.EqualValues(t, false, isResharing)
	})

	t.Run("resharing", func(t *testing.T) {
		params := &InstanceParams{
			operatorsOld: []uint32{1, 2, 3},
		}

		isResharing := params.isResharing()
		require.EqualValues(t, true, isResharing)
	})
}

func Test_inOldCommittee(t *testing.T) {
	t.Run("true case", func(t *testing.T) {
		params := &InstanceParams{
			operatorID:   1,
			operators:    []uint32{5, 6, 7, 8, 9},
			operatorsOld: []uint32{1, 2, 3},
		}

		ok := params.inOldCommittee()
		require.EqualValues(t, true, ok)
	})

	t.Run("false case", func(t *testing.T) {
		params := &InstanceParams{
			operatorID:   5,
			operators:    []uint32{5, 6, 7, 8, 9},
			operatorsOld: []uint32{1, 2, 3},
		}

		ok := params.inOldCommittee()
		require.EqualValues(t, false, ok)
	})
}

func Test_inNewCommittee(t *testing.T) {
	t.Run("true case", func(t *testing.T) {
		params := &InstanceParams{
			operatorID:   5,
			operators:    []uint32{5, 6, 7, 8, 9},
			operatorsOld: []uint32{1, 2, 3},
		}

		ok := params.inNewCommittee()
		require.EqualValues(t, true, ok)
	})

	t.Run("false case", func(t *testing.T) {
		params := &InstanceParams{
			operatorID:   1,
			operators:    []uint32{5, 6, 7, 8, 9},
			operatorsOld: []uint32{1, 2, 3},
		}

		ok := params.inNewCommittee()
		require.EqualValues(t, false, ok)
	})
}
