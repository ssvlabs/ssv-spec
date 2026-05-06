GOPATH?=$(shell go env GOPATH)
TEST_PKG?=./...

.PHONY: lint-prepare
lint-prepare:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

.PHONY: lint
lint:
	./bin/golangci-lint run -v ./...

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: test
test:
	@go test -v -race -timeout 25m "${TEST_PKG}"

.PHONY: regen-proposer-fixtures
regen-proposer-fixtures:
	@go generate ./types/testingutils

.PHONY: generate-jsons
generate-jsons:
	@go test ./types/testingutils -run TestProposerFixturesBLSDecodable -count=1
	@go generate ./qbft/spectest/generate
	@go generate ./ssv/spectest/generate
	@go generate ./types/spectest/generate

.PHONY: generate-ssz
generate-ssz:
	@go generate ./qbft/
	@go generate ./ssv/
	@go generate ./types/