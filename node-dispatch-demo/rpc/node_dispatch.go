package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"

	"oa.98ent.com/p9/node-dispatch/rpc/internal/config"
	dispatchserviceServer "oa.98ent.com/p9/node-dispatch/rpc/internal/server/dispatchservice"
	nodeserviceServer "oa.98ent.com/p9/node-dispatch/rpc/internal/server/nodeservice"
	pingserviceServer "oa.98ent.com/p9/node-dispatch/rpc/internal/server/pingservice"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/websocket"
	"oa.98ent.com/p9/node-dispatch/rpc/pb/nodedispatchrpc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/node_dispatch.yaml", "the config file")

func main() {
	// 解析启动参数
	flag.Parse()

	// 加载服务配置
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 创建服务上下文
	ctx := svc.NewServiceContext(c)
	defer ctx.DB.Close()

	// 执行数据库自动迁移
	ctx.MustMigrate()

	// 创建RPC服务
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

	// 创建WebSocket服务
	webSocketServer := websocket.NewServer(ctx)
	defer webSocketServer.Stop()

	// 启动WebSocket服务
	go func() {
		fmt.Printf("Starting websocket server at %s...\n", c.WebSocketListenOn)

		err := webSocketServer.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logx.Must(err)
		}
	}()

	// 启动RPC服务
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
