package main

import (
	"flag"
	"fmt"

	"oa.98ent.com/p9/platform-game/rpc/internal/config"
	categoryserver "oa.98ent.com/p9/platform-game/rpc/internal/server/gamecategoryservice"
	channelserver "oa.98ent.com/p9/platform-game/rpc/internal/server/gamechannelservice"
	currencyserver "oa.98ent.com/p9/platform-game/rpc/internal/server/gamecurrencyservice"
	providerserver "oa.98ent.com/p9/platform-game/rpc/internal/server/gameproviderservice"
	gameserver "oa.98ent.com/p9/platform-game/rpc/internal/server/gameservice"
	checkpointserver "oa.98ent.com/p9/platform-game/rpc/internal/server/gamesynccheckpointservice"
	pingserver "oa.98ent.com/p9/platform-game/rpc/internal/server/pingservice"
	syncserver "oa.98ent.com/p9/platform-game/rpc/internal/server/syncservice"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platform_game_pb "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

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
		// Register all individual services
		platform_game_pb.RegisterPingServiceServer(grpcServer, pingserver.NewPingServiceServer(ctx))
		platform_game_pb.RegisterGameServiceServer(grpcServer, gameserver.NewGameServiceServer(ctx))
		platform_game_pb.RegisterGameCategoryServiceServer(grpcServer, categoryserver.NewGameCategoryServiceServer(ctx))
		platform_game_pb.RegisterGameProviderServiceServer(grpcServer, providerserver.NewGameProviderServiceServer(ctx))
		platform_game_pb.RegisterGameChannelServiceServer(grpcServer, channelserver.NewGameChannelServiceServer(ctx))
		platform_game_pb.RegisterGameCurrencyServiceServer(grpcServer, currencyserver.NewGameCurrencyServiceServer(ctx))
		platform_game_pb.RegisterGameSyncCheckpointServiceServer(grpcServer, checkpointserver.NewGameSyncCheckpointServiceServer(ctx))
		platform_game_pb.RegisterSyncServiceServer(grpcServer, syncserver.NewSyncServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
