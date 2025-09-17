package qbft

//go:generate rm -f ./messages_encoding.go
//go:generate go run github.com/ferranbt/fastssz/sszgen --path messages.go --include ./types.go,../types/signer.go,../types/operator.go --exclude-objs OperatorID,ProcessingMessage
