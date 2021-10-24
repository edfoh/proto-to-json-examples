.PHONY: proto-gen
## proto-gen: compliles and generates example protobuf
proto-gen:
	@protoc --go_out=. ./examples/*.proto

.PHONY: proto2json-templates
## proto2json-templates: runs proto2json-templates
proto2json-templates:
	@go run ./cmd/proto2json-templates/main.go

.PHONY: setup
## setup: installs protoc gen tool
setup: 
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

.PHONY: deps
## deps: ensure go mod dependencies
deps: 
	@go mod tidy
	@go mod vendor

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
