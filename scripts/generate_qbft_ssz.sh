#!/bin/bash
set -e

echo "Generating SSZ encodings for qbft..."

cd qbft

# Remove old generated file
rm -f ./messages_encoding.go

# Generate messages encoding
echo "Generating messages_encoding.go..."
go run github.com/ferranbt/fastssz/sszgen --path messages.go \
    --include ./types.go,../types/signer.go,../types/operator.go \
    --exclude-objs OperatorID,ProcessingMessage

echo "QBFT SSZ generation complete!"