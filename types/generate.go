package types

//go:generate rm -f ./operator_encoding.go
//go:generate go run github.com/ferranbt/fastssz/sszgen --path operator.go --include ./committee_id.go,./domain_type.go --exclude-objs OperatorID

//go:generate rm -f ./share_encoding.go
//go:generate sh -c "go run github.com/ferranbt/fastssz/sszgen --path share.go --include $(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/phase0,./operator.go,./messages.go,./signer.go,./domain_type.go"

//go:generate rm -f ./messages_encoding.go
//go:generate go run github.com/ferranbt/fastssz/sszgen --path messages.go --include ./operator.go --exclude-objs ValidatorPK,MessageID,MsgType,ShareValidatorPK

//go:generate rm -f ./beacon_types_encoding.go
//go:generate sh -c "go run github.com/ferranbt/fastssz/sszgen --path beacon_types.go --include $(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/phase0 --exclude-objs BeaconNetwork,BeaconRole,CommitteeDuty,AggregatorCommitteeDuty"

//go:generate rm -f ./partial_sig_message_encoding.go
//go:generate sh -c "go run github.com/ferranbt/fastssz/sszgen --path partial_sig_message.go --include $(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/phase0,./signer.go,./operator.go --exclude-objs PartialSigMsgType"

//go:generate rm -f ./consensus_data_encoding.go
//go:generate sh -c "go run github.com/ferranbt/fastssz/sszgen --path consensus_data.go --include ./operator.go,./signer.go,./partial_sig_message.go,./beacon_types.go,$(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/phase0,$(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec,$(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/altair --exclude-objs Contributions,BeaconNetwork,BeaconRole"
//go:generate go run golang.org/x/tools/cmd/goimports@v0.23.0 -w consensus_data_encoding.go
