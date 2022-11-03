package qbft

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInstances_FindInstance(t *testing.T) {
	i := NewInMemContainer(HistoricalInstanceCapacity)
	i.AddNewInstance(&Instance{State: &State{Height: 1}})
	i.AddNewInstance(&Instance{State: &State{Height: 2}})
	i.AddNewInstance(&Instance{State: &State{Height: 3}})

	t.Run("find 1", func(t *testing.T) {
		require.NotNil(t, i.FindInstance(1))
	})
	t.Run("find 2", func(t *testing.T) {
		require.NotNil(t, i.FindInstance(2))
	})
	t.Run("find 5", func(t *testing.T) {
		require.Nil(t, i.FindInstance(5))
	})
}

func TestInstances_addNewInstance(t *testing.T) {
	t.Run("add to full", func(t *testing.T) {
		i := NewInMemContainer(HistoricalInstanceCapacity)
		i.AddNewInstance(&Instance{State: &State{Height: 5}})
		i.AddNewInstance(&Instance{State: &State{Height: 4}})
		i.AddNewInstance(&Instance{State: &State{Height: 3}})
		i.AddNewInstance(&Instance{State: &State{Height: 2}})
		i.AddNewInstance(&Instance{State: &State{Height: 1}})

		i.AddNewInstance(&Instance{State: &State{Height: 6}})

		require.EqualValues(t, 6, i.FindInstanceByPosition(0).State.Height)
		require.EqualValues(t, 1, i.FindInstanceByPosition(1).State.Height)
		require.EqualValues(t, 2, i.FindInstanceByPosition(2).State.Height)
		require.EqualValues(t, 3, i.FindInstanceByPosition(3).State.Height)
		require.EqualValues(t, 4, i.FindInstanceByPosition(4).State.Height)
	})

	t.Run("add to empty", func(t *testing.T) {
		i := NewInMemContainer(HistoricalInstanceCapacity)
		i.AddNewInstance(&Instance{State: &State{Height: 1}})

		require.EqualValues(t, 1, i.FindInstanceByPosition(0).State.Height)
		require.Nil(t, i.FindInstanceByPosition(1))
		require.Nil(t, i.FindInstanceByPosition(2))
		require.Nil(t, i.FindInstanceByPosition(3))
		require.Nil(t, i.FindInstanceByPosition(4))
	})

	t.Run("add to semi full", func(t *testing.T) {
		i := NewInMemContainer(HistoricalInstanceCapacity)
		i.AddNewInstance(&Instance{State: &State{Height: 3}})
		i.AddNewInstance(&Instance{State: &State{Height: 2}})
		i.AddNewInstance(&Instance{State: &State{Height: 1}})

		i.AddNewInstance(&Instance{State: &State{Height: 4}})

		require.EqualValues(t, 4, i.FindInstanceByPosition(0).State.Height)
		require.EqualValues(t, 1, i.FindInstanceByPosition(1).State.Height)
		require.EqualValues(t, 2, i.FindInstanceByPosition(2).State.Height)
		require.EqualValues(t, 3, i.FindInstanceByPosition(3).State.Height)
		require.Nil(t, i.FindInstanceByPosition(4))
	})
}

func TestController_Marshaling(t *testing.T) {
	c := testingControllerStruct()

	byts, err := c.Encode()
	require.NoError(t, err)

	fmt.Printf("r - %s", string(byts))
	decoded := &Controller{}
	require.NoError(t, decoded.Decode(byts))

	bytsDecoded, err := decoded.Encode()
	require.NoError(t, err)
	require.EqualValues(t, byts, bytsDecoded)
}
