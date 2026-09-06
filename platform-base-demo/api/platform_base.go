// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package main

import (
	"flag"
	"fmt"

	"oa.98ent.com/p9/platform-base/api/internal/config"
	"oa.98ent.com/p9/platform-base/api/internal/handler"
	"oa.98ent.com/p9/platform-base/api/internal/svc"
	"oa.98ent.com/p9/platform-base/pkg/api/errorhandler"
	"oa.98ent.com/p9/platform-base/pkg/api/response"
	"oa.98ent.com/p9/platform-base/pkg/api/validate"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/platform_base.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 创建参数校验器
	v, err := validate.New(c.I18n.DefaultLanguage)
	logx.Must(err)

	// 注册全局参数校验器
	httpx.SetValidator(v)

	// 创建API服务
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 创建服务上下文
	ctx := svc.NewServiceContext(c)

	// 开发和测试环境返回调试信息
	debug := c.Mode == service.DevMode || c.Mode == service.TestMode

	// 注册全局错误响应处理器
	httpx.SetErrorHandlerCtx(
		errorhandler.New(ctx.Trans, debug).Handle,
	)

	// 注册全局成功响应处理器
	httpx.SetOkHandler(response.Ok)

	// 注册Core国际化中间件
	server.Use(ctx.CoreI18n)

	// 注册API语言中间件
	server.Use(ctx.Language)

	// 注册全局错误日志中间件
	server.Use(ctx.ErrorLog)

	// 注册API路由
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
