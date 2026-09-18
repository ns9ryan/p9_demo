// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"

	"oa.98ent.com/p9/platform-base/pkg/api/errorhandler"
	"oa.98ent.com/p9/platform-base/pkg/api/response"
	"oa.98ent.com/p9/platform-base/pkg/api/validate"
	"oa.98ent.com/p9/platform-game/api/internal/catalog"
	"oa.98ent.com/p9/platform-game/api/internal/config"
	"oa.98ent.com/p9/platform-game/api/internal/handler"
	gamemiddleware "oa.98ent.com/p9/platform-game/api/internal/middleware"
	"oa.98ent.com/p9/platform-game/api/internal/svc"
)

var configFile = flag.String("f", "etc/platform_game.yaml", "the config file")

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

	// 开发和测试环境返回调试信息
	debug := c.Mode == service.DevMode || c.Mode == service.TestMode

	// 创建服务上下文
	ctx := svc.NewServiceContext(c)

	// 注册全局错误响应处理器
	httpx.SetErrorHandlerCtx(
		errorhandler.New(ctx.Trans, debug).Handle,
	)

	// 注册全局成功响应处理器
	httpx.SetOkHandler(response.Ok)

	// 注册API语言中间件
	server.Use(ctx.Language)

	// 注册全局错误日志中间件
	server.Use(ctx.ErrorLog)

	// 注册全局中间件
	server.Use(gamemiddleware.LanguageMiddleware())

	// 注册菜单、API目录、多语言数据
	logx.Must(catalog.Register(ctx))
	handler.RegisterHandlers(server, ctx)
	registerSwagger(server)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

//go:embed swagger/platform_game.json
var swaggerSpec []byte

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <title>Platform Game API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.ui = SwaggerUIBundle({
      url: "/swagger.json",
      dom_id: "#swagger-ui",
      presets: [SwaggerUIBundle.presets.apis],
      layout: "BaseLayout"
    });
  </script>
</body>
</html>
`

func registerSwagger(server *rest.Server) {
	server.AddRoutes([]rest.Route{
		{
			Method: http.MethodGet,
			Path:   "/swagger.json",
			Handler: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				_, _ = w.Write(swaggerSpec)
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/swagger",
			Handler: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = w.Write([]byte(swaggerUIHTML))
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/swagger/index.html",
			Handler: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = w.Write([]byte(swaggerUIHTML))
			},
		},
	})
}
