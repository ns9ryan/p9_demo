# ---- 仓库独立配置，新建仓库时需要修改 ----

# 仓库独立配置, 新建仓库时需要修改
MODULE := oa.98ent.com/p9/platform-base
RPC_PROTO := platform_base.proto

# ---- 下面的配置通常不需要修改 ----

GO ?= go
GOCTL ?= goctl

STYLE := go_zero

# API
API_FILE := api/desc/main.api
API_DIR := api

# RPC
RPC_PROTO_PATH := rpc/proto

# Ent
ENT_DIR := rpc/ent

.PHONY: api
api: # 生成 API 代码
	$(GOCTL) api go --api $(API_FILE) --dir $(API_DIR) --style=$(STYLE)

.PHONY: rpc
rpc: # 生成 RPC 代码
	cd $(RPC_PROTO_PATH) && $(GOCTL) rpc protoc $(RPC_PROTO) --go_out=../.. --go-grpc_out=../.. --zrpc_out=.. --go_opt=module=$(MODULE) --go-grpc_opt=module=$(MODULE) --module=$(MODULE) -I . --multiple --style=$(STYLE)

.PHONY: ent-new
ent-new: # 创建新的 Ent Schema，需要指定 name 参数
ifndef name
	$(error name is required, example: make ent-new name=Language)
endif
	cd rpc && $(GO) run entgo.io/ent/cmd/ent new $(name)

.PHONY: ent
ent: # 生成 Ent 代码
	$(GO) generate ./$(ENT_DIR)

.PHONY: gen
gen: api rpc ent # 生成 API、RPC 和 Ent 代码