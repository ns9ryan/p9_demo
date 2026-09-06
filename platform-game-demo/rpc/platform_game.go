package main

import (
	"flag"
	"fmt"

	"oa.98ent.com/p9/platform-game/rpc/internal/config"
	"oa.98ent.com/p9/platform-game/rpc/internal/server"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/platform_game.yaml", "the config file")

func main() {
	flag.Parse()

	c, err := config.LoadConfig(*configFile)
	if err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(*c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		platformgame.RegisterPlatformGameServiceServer(grpcServer, server.NewPlatformGameServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
