// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"oa.98ent.com/p9/common/i18n"
	"oa.98ent.com/p9/common/response"
	"oa.98ent.com/p9/common/validate"
	"oa.98ent.com/p9/platform-operator/api/internal/catalog"
	"oa.98ent.com/p9/platform-operator/api/internal/config"
	"oa.98ent.com/p9/platform-operator/api/internal/handler"
	"oa.98ent.com/p9/platform-operator/api/internal/svc"
)

var configFile = flag.String("f", "etc/platform_operator.yaml", "the config file")

func main() {
	flag.Parse()

	// 加载服务配置
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

	// 注册菜单和API目录
	logx.Must(catalog.Register(ctx.Core))

	// 设置HTTP响应格式
	response.SetupHTTPX(ctx.Trans, i18n.CodeOperator, c.IsDebug())

	// 注册Core国际化中间件
	server.Use(ctx.CoreI18n)

	// 注册全局错误日志中间件
	server.Use(ctx.ErrorLog)

	// 注册API路由
	handler.RegisterHandlers(server, ctx)

	// 启动API服务
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
