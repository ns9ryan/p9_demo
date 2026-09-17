package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"oa.98ent.com/p9/node-agent/internal/config"
	"oa.98ent.com/p9/node-agent/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/node_agent.yaml", "the config file")

func main() {
	// 解析启动参数
	flag.Parse()

	// 加载服务配置
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 创建服务上下文
	svcCtx := svc.NewServiceContext(c)

	// 监听服务停止信号
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 启动WebSocket连接
	err := svcCtx.WebSocket.Run(ctx)
	if err != nil && ctx.Err() == nil {
		logx.Must(err)
	}
}
