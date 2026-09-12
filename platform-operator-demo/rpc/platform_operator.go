package main

import (
	"flag"
	"fmt"

	"oa.98ent.com/p9/platform-operator/rpc/internal/config"
	agentlineallocationserviceServer "oa.98ent.com/p9/platform-operator/rpc/internal/server/agentlineallocationservice"
	basicresourceallocationserviceServer "oa.98ent.com/p9/platform-operator/rpc/internal/server/basicresourceallocationservice"
	domainserviceServer "oa.98ent.com/p9/platform-operator/rpc/internal/server/domainservice"
	languageallocationserviceServer "oa.98ent.com/p9/platform-operator/rpc/internal/server/languageallocationservice"
	operatorprofileserviceServer "oa.98ent.com/p9/platform-operator/rpc/internal/server/operatorprofileservice"
	operatorserviceServer "oa.98ent.com/p9/platform-operator/rpc/internal/server/operatorservice"
	pingserviceServer "oa.98ent.com/p9/platform-operator/rpc/internal/server/pingservice"
	regionallocationserviceServer "oa.98ent.com/p9/platform-operator/rpc/internal/server/regionallocationservice"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/platform_operator.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		// Ping 服务
		operator.RegisterPingServiceServer(grpcServer, pingserviceServer.NewPingServiceServer(ctx))
		// 分站服务
		operator.RegisterOperatorServiceServer(grpcServer, operatorserviceServer.NewOperatorServiceServer(ctx))
		// 分站档案服务
		operator.RegisterOperatorProfileServiceServer(grpcServer, operatorprofileserviceServer.NewOperatorProfileServiceServer(ctx))
		// 分站域名服务
		operator.RegisterDomainServiceServer(grpcServer, domainserviceServer.NewDomainServiceServer(ctx))
		// 分站基础资源分配服务
		operator.RegisterBasicResourceAllocationServiceServer(grpcServer, basicresourceallocationserviceServer.NewBasicResourceAllocationServiceServer(ctx))
		// 分站语言分配服务
		operator.RegisterLanguageAllocationServiceServer(grpcServer, languageallocationserviceServer.NewLanguageAllocationServiceServer(ctx))
		// 分站经营地区分配服务
		operator.RegisterRegionAllocationServiceServer(grpcServer, regionallocationserviceServer.NewRegionAllocationServiceServer(ctx))
		// 分站代理子线路分配服务
		operator.RegisterAgentLineAllocationServiceServer(grpcServer, agentlineallocationserviceServer.NewAgentLineAllocationServiceServer(ctx))

		// 开发和测试环境额外开启服务反射
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
