package qbft

//go:generate rm -f ./commit_extra_load_encoding.go
//go:generate go run github.com/ferranbt/fastssz/sszgen --path commit_extra_load.go --include ./types.go,../types/signer.go

//go:generate rm -f ./messages_encoding.go
//go:generate go run github.com/ferranbt/fastssz/sszgen --path messages.go --include ./types.go,../types/signer.go,../types/operator.go,./commit_extra_load.go --exclude-objs OperatorID
