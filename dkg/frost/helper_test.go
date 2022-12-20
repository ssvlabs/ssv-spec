package frost

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_haveSameRoot(t *testing.T) {
	t.Run("true case", func(t *testing.T) {
		msg := testSignedMessage(Preparation, 1)
		msg2 := testSignedMessage(Preparation, 1)
		actual := haveSameRoot(msg, msg2)
		require.EqualValues(t, true, actual)
	})

	t.Run("false case", func(t *testing.T) {
		msg := testSignedMessage(Preparation, 1)
		msg2 := testSignedMessage(Preparation, 2)
		actual := haveSameRoot(msg, msg2)
		require.EqualValues(t, false, actual)
	})
}
