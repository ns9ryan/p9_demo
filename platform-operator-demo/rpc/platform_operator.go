package main

import (
	"flag"
	"fmt"

	"oa.98ent.com/p9/platform-operator/rpc/internal/config"
	agentlineallocationserviceServer "oa.98ent.com/p9/platform-operator/rpc/internal/server/agentlineallocationservice"
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
		operator.RegisterPingServiceServer(grpcServer, pingserviceServer.NewPingServiceServer(ctx))
		operator.RegisterOperatorServiceServer(grpcServer, operatorserviceServer.NewOperatorServiceServer(ctx))
		operator.RegisterOperatorProfileServiceServer(grpcServer, operatorprofileserviceServer.NewOperatorProfileServiceServer(ctx))
		operator.RegisterDomainServiceServer(grpcServer, domainserviceServer.NewDomainServiceServer(ctx))
		operator.RegisterLanguageAllocationServiceServer(grpcServer, languageallocationserviceServer.NewLanguageAllocationServiceServer(ctx))
		operator.RegisterRegionAllocationServiceServer(grpcServer, regionallocationserviceServer.NewRegionAllocationServiceServer(ctx))
		operator.RegisterAgentLineAllocationServiceServer(grpcServer, agentlineallocationserviceServer.NewAgentLineAllocationServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
