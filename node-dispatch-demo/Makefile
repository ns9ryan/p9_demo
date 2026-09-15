GOCTL ?= goctl

STYLE := go_zero
MODULE := $(shell go list -m)

API_FILE := api/desc/main.api
API_DIR := api

RPC_PROTO ?= rpc/proto/ping.proto
RPC_PROTO_PATH := rpc/proto
RPC_DIR := rpc

.PHONY: api rpc gen

api:
	$(GOCTL) api go --api $(API_FILE) --dir $(API_DIR) --style=$(STYLE)

rpc:
	$(GOCTL) rpc protoc $(RPC_PROTO) --go_out=. --go-grpc_out=. --zrpc_out=$(RPC_DIR) --go_opt=module=$(MODULE) --go-grpc_opt=module=$(MODULE) --module=$(MODULE) -I $(RPC_PROTO_PATH) --multiple --style=$(STYLE)

gen: api rpc