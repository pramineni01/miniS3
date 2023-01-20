SRC_DIR=./protos
DST_DIR=./gen/minis3
SERVER_VERSION="1.0"
CLIENT_VERSION="1.0"
# override platform specific commands when os is windows
ifeq ($(OS),Windows_NT)
	SHELL := powershell.exe
	INTERNAL_PROTO_FILES:=$$(Get-ChildItem -Path .\internal  -Filter *.proto -Recurse -File -Name | ForEach-Object { ".\internal\$$_"})
	API_PROTO_FILES:=$$(Get-ChildItem -Path .\api  -Filter *.proto -Recurse -File -Name | ForEach-Object { ".\api\$$_"})
	VERSION:=$$(git describe --tags --always)
endif

.PHONY: gen server client

all: gen server client
# init env
init:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.6.0
	go install github.com/envoyproxy/protoc-gen-validate@v0.6.1
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	

default: gen

gen:
	@echo "Generating grpc stubs"
	@protoc -I=${SRC_DIR}  --go_out=${DST_DIR} --go-grpc_out=${DST_DIR} --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative ${SRC_DIR}/mini_s3.proto
	@protoc -I=${SRC_DIR} --grpc-gateway_out=logtostderr=true,allow_delete_body=true,grpc_api_configuration=${SRC_DIR}/mini_s3_rest.yaml:${DST_DIR}  ${SRC_DIR}/mini_s3.proto
	@protoc -I=protos --openapiv2_out ${DST_DIR} --openapiv2_opt logtostderr=true ${SRC_DIR}/mini_s3.proto

server:
	@echo "Building mini s3 server"
	@(cd server && CGO_ENABLED=0 && GOOS=darwin && GOARCH=amd64 && go build --ldflags "-s -w -X 'main.BuildVersion=${SERVER_VERSION}'" -o ../mini_s3_server)

client:
	@echo "Building mini s3 client"
	@(cd client && CGO_ENABLED=0 && GOOS=darwin && GOARCH=amd64 && go build --ldflags "-s -w -X 'main.BuildVersion=${CLIENT_VERSION}'" -o ../lc)