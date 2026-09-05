package main

import (
	"flag"
	"fmt"

	"oa.98ent.com/p9/platform-base/rpc/internal/config"
	currencyserviceServer "oa.98ent.com/p9/platform-base/rpc/internal/server/currencyservice"
	languageserviceServer "oa.98ent.com/p9/platform-base/rpc/internal/server/languageservice"
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
		base.RegisterPingServiceServer(grpcServer, pingserviceServer.NewPingServiceServer(ctx))
		base.RegisterLanguageServiceServer(grpcServer, languageserviceServer.NewLanguageServiceServer(ctx))
		base.RegisterTimezoneServiceServer(grpcServer, timezoneserviceServer.NewTimezoneServiceServer(ctx))
		base.RegisterCurrencyServiceServer(grpcServer, currencyserviceServer.NewCurrencyServiceServer(ctx))
		base.RegisterRegionServiceServer(grpcServer, regionserviceServer.NewRegionServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
