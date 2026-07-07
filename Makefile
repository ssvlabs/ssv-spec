GOPATH?=$(shell go env GOPATH)
TEST_PKG?=./...
# Keep in sync with the version pinned in .github/workflows/lint.yaml
GOLANGCI_LINT_VERSION?=v2.12.2

.PHONY: lint-prepare
lint-prepare:
	@tmp=$$(mktemp) \
		&& curl -sfL https://golangci-lint.run/install.sh -o "$$tmp" \
		&& sh "$$tmp" $(GOLANGCI_LINT_VERSION) \
		&& rm -f "$$tmp"

.PHONY: lint
lint:
	./bin/golangci-lint run -v ./...

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: test
test:
	@go test -v -race -timeout 25m "${TEST_PKG}"

.PHONY: generate-jsons
generate-jsons:
	@go generate ./...

.PHONY: generate-ssz
generate-ssz:
	@go generate ./qbft/
	@go generate ./ssv/
	@go generate ./types/
	@go generate ./types/gloas/