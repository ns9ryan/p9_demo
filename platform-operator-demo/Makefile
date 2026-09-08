# ---- 仓库独立配置，新建仓库时需要修改 ----

MODULE := oa.98ent.com/p9/platform-operator
RPC_PROTO := platform_operator.proto
BUILD_NAME := platform_operator

# 镜像版本
VERSION ?= 1.0.0

# ---- 下面的配置通常不需要修改 ----

PROJECT := $(notdir $(MODULE))

GO ?= go
GOCTL ?= goctl

STYLE := go_zero

# 构建
GOARCH ?= amd64
LDFLAGS := -s -w

# Docker
DOCKER_PLATFORM ?= linux/$(GOARCH)
DOCKER_REPO ?= 192.168.0.15:5000/p9

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

.PHONY: build-linux
build-linux: export CGO_ENABLED := 0
build-linux: export GOOS := linux
build-linux: export GOARCH := $(GOARCH)
build-linux: # 编译 Linux 可执行文件
	$(GO) build -trimpath -ldflags="$(LDFLAGS)" -o $(BUILD_NAME)_rpc ./rpc
	$(GO) build -trimpath -ldflags="$(LDFLAGS)" -o $(BUILD_NAME)_api ./api

.PHONY: docker
docker: build-linux # 编译 Linux 可执行文件并构建 Docker 镜像
	docker build --platform $(DOCKER_PLATFORM) -f Dockerfile-api -t $(DOCKER_REPO)/$(PROJECT)-api:$(VERSION) .
	docker build --platform $(DOCKER_PLATFORM) -f Dockerfile-rpc -t $(DOCKER_REPO)/$(PROJECT)-rpc:$(VERSION) .

.PHONY: publish-docker
publish-docker: # 发布 Docker 镜像
	docker push $(DOCKER_REPO)/$(PROJECT)-api:$(VERSION)
	docker push $(DOCKER_REPO)/$(PROJECT)-rpc:$(VERSION)