package types

type Root interface {
	GetRoot() ([32]byte, error)
}
