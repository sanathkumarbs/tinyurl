ROOTDIR				:= $(shell git rev-parse --show-toplevel)
COMMIT_SHA 			:= $(git rev-parse --short HEAD)
GO_LINT_VERSION		:= 1.55.1
GO_LINT_BIN			:= $(ROOTDIR)/bin/golangci-lint
GO_IMPORTS_VERSION	:= 0.16.1

.ONESHELL:

_install-oapi-codegen: 
ifeq (, $(shell which oapi-codegen))
	go install -v github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.0.0
endif

_install-golangci-lint: 
ifeq (, $(shell $(GO_LINT_BIN) --version))
	@mkdir -p $(shell dirname $(GO_LINT_BIN))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell dirname $(GO_LINT_BIN)) v$(GO_LINT_VERSION)
endif

_install-goimports: 
ifeq (, $(shell which goimports))
	go install -v golang.org/x/tools/cmd/goimports@v$(GO_IMPORTS_VERSION)
endif

_install-squawk-cli: 
ifeq (, $(shell which squawk))
	npm install -g squawk-cli@0.26.0
endif


all: gen fmt lint build

ru: run-up
rd: run-down
rdv: run-down-volumes

run-up: 
	$(info Running run-up...)
	@$(MAKE) build
	docker compose -f devenv/docker-compose.yaml up -d

run-down: 
	$(info Running run-down...)
	docker compose -f devenv/docker-compose.yaml down
	rm -rf tmp/*

run-down-volumes: 
	$(info Running run-down...)
	docker compose -f devenv/docker-compose.yaml down --volumes
	rm -rf tmp/*

docker-prune:
	$(info Running docker-prune...)
	docker system prune -a --force
	docker builder prune -a --force

psql:
	$(info Running exec into postgres...)
	docker exec -it devenv-postgres-1 psql -d tiny -U postgres

clean:
	$(info Running clean...)
	@rm -rf bin/*

build:
	@$(MAKE) build-go
	@$(MAKE) build-image

build-go:
	$(info Running build-go...)
	CGO_ENABLED=0 GOOS=linux go build -o bin/tiny cmd/tiny/main.go

build-image:
	$(info Running build-image...)
	@$(ROOTDIR)/scripts/build-image.sh

fmt:
	$(info Running goimports...)
	@$(MAKE) --quiet _install-goimports
	@goimports -w -e $$(find . -type f -name '*.go' -not -path "./.gocache/*")

lint: 
	$(info Running lint...)
	@$(MAKE) lint-go
	@$(MAKE) lint-sql

lint-go:
	$(info Running lint-go...)
	@$(MAKE) --quiet _install-golangci-lint	
	time $(GO_LINT_BIN) run ./...

lint-sql:
	$(info Running lint-sql...)
	@$(MAKE) --quiet _install-squawk-cli
	time squawk internal/tiny/db/migrations/*

gen:
	@$(MAKE) generate
	squawk internal/tiny/db/migrations/*

generate:
	$(info Running generate...)
	@$(MAKE) --quiet _install-oapi-codegen	
	@mkdir -p ./pkg/api/services/v1/tiny/
	@oapi-codegen -generate client -package tiny ./api/services/v1/tiny/tiny-openapi.yaml > ./pkg/api/services/v1/tiny/tiny-client.gen.go
	@oapi-codegen -generate server,spec -package tiny ./api/services/v1/tiny/tiny-openapi.yaml > ./pkg/api/services/v1/tiny/tiny-server.gen.go
	@oapi-codegen -generate types -package tiny ./api/services/v1/tiny/tiny-openapi.yaml > ./pkg/api/services/v1/tiny/tiny-types.gen.go