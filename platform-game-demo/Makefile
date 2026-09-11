# Custom configuration | 独立配置
SERVICE=platform-game
SERVICE_STYLE=platform_game
SERVICE_LOWER=platform-game
SERVICE_SNAKE=platform_game
SERVICE_DASH=platform-game

# The project version, if you don't use git, you should set it manually | 项目版本，如果不使用git请手动设置
# VERSION=$(shell git describe --tags --always)
VERSION=1.0.0

PROJECT_STYLE=go_zero

GOCTL ?= goctl
GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")

STYLE := go_zero
MODULE := $(shell go list -m)

API_FILE := api/desc/main.api
API_DIR := api

RPC_PROTO ?= rpc/proto/platform_game.proto
RPC_PROTO_PATH := rpc/proto
RPC_DIR := rpc

# The arch of the build | 构建的架构
GOARCH=amd64
DOCKER_PLATFORM=linux/amd64

# The repository of docker | Docker 仓库地址
DOCKER_REPO=192.168.0.15:5000/p9

.PHONY: fmt api rpc gen swagger docker publish-docker build-win build-mac build-linux help

fmt: # Format the codes | 格式化代码
	$(GOFMT) -w $(GOFILES)

api: # Generate API files | 生成 API 的代码
	$(GOCTL) api go --api $(API_FILE) --dir $(API_DIR) --style=$(STYLE)

rpc: # Generate RPC files from proto | 生成 RPC 的代码
	cd $(RPC_PROTO_PATH) && $(GOCTL) rpc protoc platform_game.proto --go_out=../.. --go-grpc_out=../.. --zrpc_out=.. --go_opt=module=$(MODULE) --go-grpc_opt=module=$(MODULE) --module=$(MODULE) -I . --multiple --name-from-filename --style=$(STYLE)

rpc-vendor: # Generate RPC files from proto | 生成 RPC 的代码
	cd $(RPC_PROTO_PATH) && $(GOCTL) rpc protoc vendor.proto --go_out=../.. --go-grpc_out=../.. --zrpc_out=.. --go_opt=module=$(MODULE) --go-grpc_opt=module=$(MODULE) --module=$(MODULE) -I . --multiple --name-from-filename --style=$(STYLE)


ent: # 生成 Ent 代码
	goctl model pg datasource \
	--url="postgres://root:root@127.0.0.1:5432/app?sslmode=disable" \
	--table="game,game_category,game_provider,game_channel,game_currency,game_sync_checkpoint,sys_currency" \
	--dir ./rpc/ent \
	--style=$(STYLE)

swagger: # Generate Swagger API documentation | 生成 API Swagger 文档
	mkdir -p api/swagger
	$(GOCTL) api swagger --filename=$(SERVICE_STYLE) --dir=./api/swagger --api=$(API_FILE)
	@echo "Generate Swagger files successfully"

gen: api rpc swagger # Generate all | 生成所有代码

.PHONY: docker
docker: build-linux # Build the docker image | 构建 docker 镜像
	docker build --platform $(DOCKER_PLATFORM) -f Dockerfile-api -t $(DOCKER_REPO)/$(SERVICE_DASH)-api:${VERSION} .
	docker build --platform $(DOCKER_PLATFORM) -f Dockerfile-rpc -t $(DOCKER_REPO)/$(SERVICE_DASH)-rpc:${VERSION} .
	@echo "Build docker successfully"

.PHONY: publish-docker
publish-docker: # Publish docker image | 发布 docker 镜像
	docker push $(DOCKER_REPO)/$(SERVICE_DASH)-rpc:${VERSION}
	docker push $(DOCKER_REPO)/$(SERVICE_DASH)-api:${VERSION}
	@echo "Publish docker successfully"

.PHONY: build-win
build-win: # Build project for Windows | 构建Windows下的可执行文件
	env CGO_ENABLED=0 GOOS=windows GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -trimpath -o $(SERVICE_STYLE)_rpc.exe ./rpc/$(SERVICE_STYLE).go
	env CGO_ENABLED=0 GOOS=windows GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -trimpath -o $(SERVICE_STYLE)_api.exe ./api/platform_game.go
	@echo "Build project for Windows successfully"

.PHONY: build-mac
build-mac: # Build project for MacOS | 构建MacOS下的可执行文件
	env CGO_ENABLED=0 GOOS=darwin GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -trimpath -o $(SERVICE_STYLE)_rpc ./rpc/$(SERVICE_STYLE).go
	env CGO_ENABLED=0 GOOS=darwin GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -trimpath -o $(SERVICE_STYLE)_api ./api/platform_game.go
	@echo "Build project for MacOS successfully"

.PHONY: build-linux
build-linux: # Build project for Linux | 构建Linux下的可执行文件
	env CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -trimpath -o $(SERVICE_STYLE)_rpc ./rpc/$(SERVICE_STYLE).go
	env CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -trimpath -o $(SERVICE_STYLE)_api ./api/platform_game.go
	@echo "Build project for Linux successfully"

.PHONY: help
help: # Show help | 显示帮助
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done