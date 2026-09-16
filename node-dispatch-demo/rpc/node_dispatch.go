package main

import (
	"flag"
	"fmt"

	"oa.98ent.com/p9/node-dispatch/rpc/internal/config"
	dispatchserviceServer "oa.98ent.com/p9/node-dispatch/rpc/internal/server/dispatchservice"
	nodeserviceServer "oa.98ent.com/p9/node-dispatch/rpc/internal/server/nodeservice"
	pingserviceServer "oa.98ent.com/p9/node-dispatch/rpc/internal/server/pingservice"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
	"oa.98ent.com/p9/node-dispatch/rpc/pb/nodedispatchrpc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/node_dispatch.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	defer ctx.DB.Close()

	// 执行数据库自动迁移
	ctx.MustMigrate()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		// Ping 服务
		nodedispatchrpc.RegisterPingServiceServer(grpcServer, pingserviceServer.NewPingServiceServer(ctx))

		// 节点服务
		nodedispatchrpc.RegisterNodeServiceServer(grpcServer, nodeserviceServer.NewNodeServiceServer(ctx))

		// 调度服务
		nodedispatchrpc.RegisterDispatchServiceServer(grpcServer, dispatchserviceServer.NewDispatchServiceServer(ctx))

		// 开发和测试环境额外开启服务反射
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
