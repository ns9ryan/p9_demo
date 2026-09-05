# Custom configuration | 独立配置
SERVICE=Core
SERVICE_STYLE=core
SERVICE_LOWER=core
SERVICE_SNAKE=core
SERVICE_DASH=core

# The project version, if you don't use git, you should set it manually | 项目版本，如果不使用git请手动设置
# VERSION=$(shell git describe --tags --always)
VERSION=1.0.0

PROJECT_STYLE=go_zero

GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")

# The arch of the build | 构建的架构
GOARCH=amd64
DOCKER_PLATFORM=linux/amd64

# The repository of docker | Docker 仓库地址
DOCKER_REPO=192.168.0.15:5000/p9

.PHONY: fmt
fmt: # Format the codes | 格式化代码
	$(GOFMT) -w $(GOFILES)

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

.PHONY: gen-api
gen-api: # Generate API files | 生成 API 的代码
	goctls api go --api ./api/desc/all.api --dir ./api --style=$(PROJECT_STYLE)
	find ./api/internal/handler -name '*.go' -exec sed -i '' 's/httpx.Parse(r, \&req, true)/httpx.Parse(r, \&req)/g' {} +
	mkdir -p ./api/swagger
	goctls api swagger --filename=$(SERVICE_STYLE) --dir=./api/swagger --api=./api/desc/all.api --makefile=./Makefile
	@echo "Generate API files successfully"

.PHONY: gen-rpc
gen-rpc: # Generate RPC files from proto | 生成 RPC 的代码
	goctls rpc protoc ./rpc/$(SERVICE_STYLE).proto --style=$(PROJECT_STYLE) --go_out=./rpc/types --go-grpc_out=./rpc/types --zrpc_out=./rpc --style=$(PROJECT_STYLE)
	@echo "Generate RPC files successfully"

.PHONY: gen-ent
gen-ent: # Generate Ent codes | 生成 Ent 的代码
	cd rpc/ent && go generate
	@echo "Generate Ent files successfully"

.PHONY: build-win
build-win: # Build project for Windows | 构建Windows下的可执行文件
	env CGO_ENABLED=0 GOOS=windows GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -trimpath -o $(SERVICE_STYLE)_rpc.exe ./rpc/$(SERVICE_STYLE).go
	env CGO_ENABLED=0 GOOS=windows GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -trimpath -o $(SERVICE_STYLE)_api.exe ./api/$(SERVICE_STYLE).go
	@echo "Build project for Windows successfully"

.PHONY: build-mac
build-mac: # Build project for MacOS | 构建MacOS下的可执行文件
	env CGO_ENABLED=0 GOOS=darwin GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -trimpath -o $(SERVICE_STYLE)_rpc ./rpc/$(SERVICE_STYLE).go
	env CGO_ENABLED=0 GOOS=darwin GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -trimpath -o $(SERVICE_STYLE)_api ./api/$(SERVICE_STYLE).go
	@echo "Build project for MacOS successfully"

.PHONY: build-linux
build-linux: # Build project for Linux | 构建Linux下的可执行文件
	env CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -trimpath -o $(SERVICE_STYLE)_rpc ./rpc/$(SERVICE_STYLE).go
	env CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -trimpath -o $(SERVICE_STYLE)_api ./api/$(SERVICE_STYLE).go
	@echo "Build project for Linux successfully"

.PHONY: help
help: # Show help | 显示帮助
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done
