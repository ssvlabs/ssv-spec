package qbft

import "encoding/json"

// InstanceContainer interface for managing multi instance lifecycle
type InstanceContainer interface {
	All() []*Instance
	FindInstance(height Height) *Instance
	FindInstanceByPosition(position int) *Instance
	AddNewInstance(instance *Instance)
	AddNewInstanceAtPosition(position int, instance *Instance)
}

// inMemContainer implements InstanceContainer
type inMemContainer struct {
	capacity  int
	container []*Instance
}

// NewInMemContainer return new instance store with capacity
func NewInMemContainer(capacity int) InstanceContainer {
	imc := &inMemContainer{
		capacity,
		make([]*Instance, capacity),
	}

	// in order to pass tests. (already with 5 nil items)
	for i := 0; i < capacity; i++ {
		imc.container[i] = nil
	}

	return imc
}

func (i *inMemContainer) All() []*Instance {
	return i.container
}

// FindInstance by height. return instance if exist or nil if not.
func (i *inMemContainer) FindInstance(height Height) *Instance {
	for _, inst := range i.container {
		if inst != nil {
			if inst.GetHeight() == height {
				return inst
			}
		}
	}
	return nil
}

// FindInstanceByPosition return instance by the store position
func (i *inMemContainer) FindInstanceByPosition(position int) *Instance {
	return i.container[position]
}

// AddNewInstance will add the new instance at index 0, pushing all others stored InstanceContainer one index up (ejecting last one if existing)
func (i *inMemContainer) AddNewInstance(instance *Instance) {
	for idx := i.capacity - 1; idx > 0; idx-- {
		i.container[idx] = i.container[idx-1]
	}
	i.container[0] = instance
}

// AddNewInstanceAtPosition adding instance at specific position in store
func (i *inMemContainer) AddNewInstanceAtPosition(position int, instance *Instance) {
	if position < 0 || position >= i.capacity {
		return
	}
	i.container[position] = instance
}

func (i *inMemContainer) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.container)
}

func (i *inMemContainer) UnmarshalJSON(input []byte) error {
	var container []*Instance
	if err := json.Unmarshal(input, &container); err != nil {
		return err
	}
	i.container = container
	i.capacity = len(container)
	return nil
}
