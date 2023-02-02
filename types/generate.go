package types

//go:generate rm -f ./operator_encoding.go
//go:generate go run .../fastssz/sszgen --path operator.go --exclude-objs OperatorID

//go:generate rm -f ./share_encoding.go
//go:generate go run .../fastssz/sszgen --path share.go --include ./operator.go,./messages.go,./signer.go

//go:generate rm -f ./messages_encoding.go
//go:generate go run .../fastssz/sszgen --path messages.go --exclude-objs ValidatorPK,MessageID,MsgType
