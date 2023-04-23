package types

type RootGetter interface {
	GetRoot() ([]byte, error)
}
