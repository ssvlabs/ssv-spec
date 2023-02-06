package alea

import "math/rand"

func SharedCoinF(round Round) byte {
	rand.Seed(int64(round))
	return byte(rand.Intn(2))
}
