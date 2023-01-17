package dkg

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInit_Validate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		init := Init{
			OperatorIDs:           []types.OperatorID{1, 2, 3, 4},
			Threshold:             3,
			WithdrawalCredentials: common.Hex2Bytes("010000000000000000000000535953b5a6040074948cf185eaa7d2abbd66808f"),
			Fork:                  spec.Version{},
		}
		require.NoError(t, init.Validate())
		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6, 7}
		init.Threshold = 5
		require.NoError(t, init.Validate())
		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		init.Threshold = 7
		require.NoError(t, init.Validate())
		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
		init.Threshold = 9
		require.NoError(t, init.Validate())
	})
	t.Run("invalid number of operators", func(t *testing.T) {
		init := Init{
			OperatorIDs:           []types.OperatorID{1, 2, 3},
			Threshold:             3,
			WithdrawalCredentials: common.Hex2Bytes("010000000000000000000000535953b5a6040074948cf185eaa7d2abbd66808f"),
			Fork:                  spec.Version{},
		}
		require.Error(t, init.Validate())
		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6}
		init.Threshold = 3
		require.Error(t, init.Validate())
		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6, 7, 8}
		init.Threshold = 5
		require.Error(t, init.Validate())
		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
		init.Threshold = 7
		require.Error(t, init.Validate())
	})
	t.Run("invalid threshold", func(t *testing.T) {
		init := Init{
			OperatorIDs:           []types.OperatorID{1, 2, 3, 4},
			Threshold:             2,
			WithdrawalCredentials: common.Hex2Bytes("010000000000000000000000535953b5a6040074948cf185eaa7d2abbd66808f"),
			Fork:                  spec.Version{},
		}
		require.Error(t, init.Validate())
		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6, 7}
		init.Threshold = 6
		require.Error(t, init.Validate())
		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		init.Threshold = 8
		require.Error(t, init.Validate())
		init.OperatorIDs = []types.OperatorID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
		init.Threshold = 8
		require.Error(t, init.Validate())
	})
	t.Run("short WithdrawalCredentials", func(t *testing.T) {
		init := Init{
			OperatorIDs:           []types.OperatorID{1, 2, 3, 4},
			Threshold:             3,
			WithdrawalCredentials: common.Hex2Bytes("010000000000000000000000535953b5a6040074948cf185eaa7d2abbd6680"),
			Fork:                  spec.Version{},
		}
		require.Error(t, init.Validate())
	})
	t.Run("long WithdrawalCredentials", func(t *testing.T) {
		init := Init{
			OperatorIDs:           []types.OperatorID{1, 2, 3, 4},
			Threshold:             3,
			WithdrawalCredentials: common.Hex2Bytes("010000000000000000000000535953b5a6040074948cf185eaa7d2abbd66808faa"),
			Fork:                  spec.Version{},
		}
		require.Error(t, init.Validate())
	})
}
