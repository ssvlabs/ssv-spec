#!/bin/bash
set -e

# Get the module directory dynamically
ETH2_CLIENT_DIR=$(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)

if [ -z "$ETH2_CLIENT_DIR" ]; then
    echo "Error: Could not find go-eth2-client module directory"
    exit 1
fi

echo "Using go-eth2-client from: $ETH2_CLIENT_DIR"

cd types

# Remove old generated files
echo "Removing old generated files..."
rm -f ./operator_encoding.go
rm -f ./share_encoding.go
rm -f ./messages_encoding.go
rm -f ./beacon_types_encoding.go
rm -f ./partial_sig_message_encoding.go
rm -f ./consensus_data_encoding.go

# Generate operator encoding
echo "Generating operator_encoding.go..."
go run github.com/ferranbt/fastssz/sszgen --path operator.go \
    --include ./committee_id.go,./domain_type.go \
    --exclude-objs OperatorID

# Generate share encoding
echo "Generating share_encoding.go..."
go run github.com/ferranbt/fastssz/sszgen --path share.go \
    --include "${ETH2_CLIENT_DIR}/spec/phase0,./operator.go,./messages.go,./signer.go,./domain_type.go"

# Generate messages encoding
echo "Generating messages_encoding.go..."
go run github.com/ferranbt/fastssz/sszgen --path messages.go \
    --include ./operator.go \
    --exclude-objs ValidatorPK,MessageID,MsgType,ShareValidatorPK

# Generate beacon types encoding
echo "Generating beacon_types_encoding.go..."
go run github.com/ferranbt/fastssz/sszgen --path beacon_types.go \
    --include "${ETH2_CLIENT_DIR}/spec/phase0" \
    --exclude-objs BeaconNetwork,BeaconRole,CommitteeDuty

# Generate partial sig message encoding
echo "Generating partial_sig_message_encoding.go..."
go run github.com/ferranbt/fastssz/sszgen --path partial_sig_message.go \
    --include "${ETH2_CLIENT_DIR}/spec/phase0,./signer.go,./operator.go" \
    --exclude-objs PartialSigMsgType

# Generate consensus data encoding
echo "Generating consensus_data_encoding.go..."
go run github.com/ferranbt/fastssz/sszgen --path consensus_data.go \
    --include "./operator.go,./signer.go,./partial_sig_message.go,./beacon_types.go,${ETH2_CLIENT_DIR}/spec/phase0,${ETH2_CLIENT_DIR}/spec,${ETH2_CLIENT_DIR}/spec/altair" \
    --exclude-objs Contributions,BeaconNetwork,BeaconRole

# Format the generated consensus data encoding
echo "Formatting consensus_data_encoding.go..."
go run golang.org/x/tools/cmd/goimports@latest -w consensus_data_encoding.go

echo "Types SSZ generation complete!"