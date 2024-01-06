.install-oapi-codegen: 
ifeq (, $(shell which oapi-codegen))
	go install -v github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.0.0
endif

gen:

generate: .install-oapi-codegen
	mkdir -p ./pkg/api/services/v1/tiny/
	oapi-codegen -generate client -package tiny ./api/services/v1/tiny/tiny-openapi.yaml > ./pkg/api/services/v1/tiny/tiny-client.gen.go
	oapi-codegen -generate server,spec -package tiny ./api/services/v1/tiny/tiny-openapi.yaml > ./pkg/api/services/v1/tiny/tiny-server.gen.go
	oapi-codegen -generate types -package tiny ./api/services/v1/tiny/tiny-openapi.yaml > ./pkg/api/services/v1/tiny/tiny-types.gen.go