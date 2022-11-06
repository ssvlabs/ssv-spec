package qbft

type InstanceContainer [HistoricalInstanceCapacity]*Instance

func (i InstanceContainer) FindInstance(height Height) *Instance {
	for _, inst := range i {
		if inst != nil {
			if inst.GetHeight() == height {
				return inst
			}
		}
	}
	return nil
}

// addNewInstance will add the new instance at index 0, pushing all other stored InstanceContainer one index up (ejecting last one if existing)
func (i *InstanceContainer) addNewInstance(instance *Instance) {
	for idx := HistoricalInstanceCapacity - 1; idx > 0; idx-- {
		i[idx] = i[idx-1]
	}
	i[0] = instance
}
