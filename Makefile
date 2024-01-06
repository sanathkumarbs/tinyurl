ROOTDIR				:= $(shell git rev-parse --show-toplevel)
GO_LINT_VERSION		:= 1.55.1
GO_LINT_BIN			:= $(ROOTDIR)/bin/golangci-lint
GO_IMPORTS_VERSION	:= 0.16.1

.install-oapi-codegen: 
ifeq (, $(shell which oapi-codegen))
	go install -v github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.0.0
endif

.install-golangci-lint: 
ifeq (, $(shell $(GO_LINT_BIN) --version))
	@mkdir -p $(shell dirname $(GO_LINT_BIN))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell dirname $(GO_LINT_BIN)) v$(GO_LINT_VERSION)
endif

.install-goimports: 
ifeq (, $(shell which goimports))
	go install -v golang.org/x/tools/cmd/goimports@v$(GO_IMPORTS_VERSION)
endif

all: gen fmt lint build

clean:
	$(info Running clean...)
	@rm -rf bin/*

build:
	$(info Running build...)
	@CGO_ENABLED=0 go build -o bin/tiny cmd/tiny/main.go

fmt: .install-goimports
	$(info Running goimports...)
	@goimports -w -e $$(find . -type f -name '*.go' -not -path "./.gocache/*")

lint: .install-golangci-lint
	$(info Running lint...)
	time $(GO_LINT_BIN) run ./...

gen: generate

generate: .install-oapi-codegen
	$(info Running generate...)
	@mkdir -p ./pkg/api/services/v1/tiny/
	@oapi-codegen -generate client -package tiny ./api/services/v1/tiny/tiny-openapi.yaml > ./pkg/api/services/v1/tiny/tiny-client.gen.go
	@oapi-codegen -generate server,spec -package tiny ./api/services/v1/tiny/tiny-openapi.yaml > ./pkg/api/services/v1/tiny/tiny-server.gen.go
	@oapi-codegen -generate types -package tiny ./api/services/v1/tiny/tiny-openapi.yaml > ./pkg/api/services/v1/tiny/tiny-types.gen.go