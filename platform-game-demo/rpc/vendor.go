package main

// import (
// 	"flag"
// 	"fmt"

// 	"oa.98ent.com/p9/platform-game/rpc/internal/config"
// 	vendorgameserviceServer "oa.98ent.com/p9/platform-game/rpc/internal/server/vendorgameservice"
// 	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
// 	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"

// 	"github.com/zeromicro/go-zero/core/conf"
// 	"github.com/zeromicro/go-zero/core/service"
// 	"github.com/zeromicro/go-zero/zrpc"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/reflection"
// )

// var configFile = flag.String("f", "etc/vendor.yaml", "the config file")

// func main() {
// 	flag.Parse()

// 	var c config.Config
// 	conf.MustLoad(*configFile, &c)
// 	ctx := svc.NewServiceContext(c)

// 	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
// 		vendors.RegisterVendorGameServiceServer(grpcServer, vendorgameserviceServer.NewVendorGameServiceServer(ctx))

// 		if c.Mode == service.DevMode || c.Mode == service.TestMode {
// 			reflection.Register(grpcServer)
// 		}
// 	})
// 	defer s.Stop()

// 	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
// 	s.Start()
// }
