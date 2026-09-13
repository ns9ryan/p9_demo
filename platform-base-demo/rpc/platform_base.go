package main

import (
	"flag"
	"fmt"

	"oa.98ent.com/p9/platform-base/rpc/internal/config"
	currencyserviceServer "oa.98ent.com/p9/platform-base/rpc/internal/server/currencyservice"
	pingserviceServer "oa.98ent.com/p9/platform-base/rpc/internal/server/pingservice"
	regionserviceServer "oa.98ent.com/p9/platform-base/rpc/internal/server/regionservice"
	timezoneserviceServer "oa.98ent.com/p9/platform-base/rpc/internal/server/timezoneservice"
	"oa.98ent.com/p9/platform-base/rpc/internal/svc"
	"oa.98ent.com/p9/platform-base/rpc/pb/base"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/platform_base.yaml", "the config file")

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
		base.RegisterPingServiceServer(grpcServer, pingserviceServer.NewPingServiceServer(ctx))
		// 时区服务
		base.RegisterTimezoneServiceServer(grpcServer, timezoneserviceServer.NewTimezoneServiceServer(ctx))
		// 货币服务
		base.RegisterCurrencyServiceServer(grpcServer, currencyserviceServer.NewCurrencyServiceServer(ctx))
		// 国家地区服务
		base.RegisterRegionServiceServer(grpcServer, regionserviceServer.NewRegionServiceServer(ctx))

		// 开发和测试环境额外开启服务反射
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
