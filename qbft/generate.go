package qbft

//go:generate rm -f ./commit_extra_load_encoding.go
//go:generate go run github.com/ferranbt/fastssz/sszgen --path commit_extra_load.go --include ./types.go,../types/signer.go,$GOPATH/pkg/mod/github.com/attestantio/go-eth2-client@v0.19.7/spec/phase0
//go:generate go install golang.org/x/tools/cmd/goimports@latest
//go:generate goimports -w commit_extra_load_encoding.go

//go:generate rm -f ./messages_encoding.go
//go:generate go run github.com/ferranbt/fastssz/sszgen --path messages.go --include ./types.go,../types/signer.go,../types/operator.go,./commit_extra_load.go,$GOPATH/pkg/mod/github.com/attestantio/go-eth2-client@v0.19.7/spec/phase0 --exclude-objs OperatorID
