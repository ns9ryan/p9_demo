package main

import (
	"flag"
	"fmt"

	"oa.98ent.com/p9/operator-base/rpc/internal/config"
	operatorserviceServer "oa.98ent.com/p9/operator-base/rpc/internal/server/operatorservice"
	"oa.98ent.com/p9/operator-base/rpc/internal/svc"
	"oa.98ent.com/p9/operator-base/rpc/pb/operatorbaserpc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/operator_base.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	defer ctx.DB.Close()

	// 执行数据库自动迁移
	ctx.MustMigrate()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		// 厅 operator 服务
		operatorbaserpc.RegisterOperatorServiceServer(grpcServer, operatorserviceServer.NewOperatorServiceServer(ctx))

		// 开发和测试环境额外开启服务反射
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
