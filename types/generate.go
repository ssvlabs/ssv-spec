package types

//go:generate rm -f ./operator_encoding.go
//go:generate go run github.com/ferranbt/fastssz/sszgen --path operator.go --include ./committee_id.go,./domain_type.go --exclude-objs OperatorID

//go:generate rm -f ./share_encoding.go
//go:generate go run github.com/ferranbt/fastssz/sszgen --path share.go --include $GOPATH/pkg/mod/github.com/attestantio/go-eth2-client@v0.21.3/spec/phase0,./operator.go,./messages.go,./signer.go,./domain_type.go

//go:generate rm -f ./messages_encoding.go
//go:generate go run github.com/ferranbt/fastssz/sszgen --path messages.go --include ./operator.go --exclude-objs ValidatorPK,MessageID,MsgType,ShareValidatorPK

//go:generate rm -f ./beacon_types_encoding.go
//go:generate go run github.com/ferranbt/fastssz/sszgen --path beacon_types.go --include $GOPATH/pkg/mod/github.com/attestantio/go-eth2-client@v0.21.3/spec/phase0 --exclude-objs BeaconNetwork,BeaconRole,CommitteeDuty

//go:generate rm -f ./partial_sig_message_encoding.go
//go:generate go run github.com/ferranbt/fastssz/sszgen --path partial_sig_message.go --include $GOPATH/pkg/mod/github.com/attestantio/go-eth2-client@v0.21.3/spec/phase0,./signer.go,./operator.go --exclude-objs PartialSigMsgType

//go:generate rm -f ./consensus_data_encoding.go
//go:generate go run github.com/ferranbt/fastssz/sszgen --path consensus_data.go --include ./operator.go,./signer.go,./partial_sig_message.go,./beacon_types.go,$GOPATH/pkg/mod/github.com/attestantio/go-eth2-client@v0.21.3/spec/phase0,$GOPATH/pkg/mod/github.com/attestantio/go-eth2-client@v0.21.3/spec,$GOPATH/pkg/mod/github.com/attestantio/go-eth2-client@v0.21.3/spec/altair --exclude-objs Contributions,BeaconNetwork,BeaconRole
//go:generate go install golang.org/x/tools/cmd/goimports@latest
//go:generate goimports -w consensus_data_encoding.go
