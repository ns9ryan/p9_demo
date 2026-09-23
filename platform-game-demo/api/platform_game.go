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
	"github.com/zeromicro/go-zero/rest"

	"oa.98ent.com/p9/platform-game/api/internal/catalog"
	"oa.98ent.com/p9/platform-game/api/internal/config"
	"oa.98ent.com/p9/platform-game/api/internal/handler"
	"oa.98ent.com/p9/platform-game/api/internal/svc"

	"oa.98ent.com/p9/common/response"
	"oa.98ent.com/p9/core/common/middleware"
)

var configFile = flag.String("f", "etc/platform_game.yaml", "the config file")

func main() {
	flag.Parse()

	// 加载服务配置
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 创建API服务
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 注册I18n中间件
	server.Use(middleware.I18n)
	// 注册客户端 IP 中间件
	server.Use(middleware.ClientIP)

	// 创建服务上下文
	ctx := svc.NewServiceContext(c)

	// 注册错误日志中间件
	server.Use(ctx.ErrorLog)

	// 注册HTTPX的OK和Error处理函数
	response.SetupHTTPX(ctx.Trans, c.GetI18nCode(), c.IsDebug())

	// 注册菜单、API目录、多语言数据
	logx.Must(catalog.Register(ctx))
	// 注册API路由
	handler.RegisterHandlers(server, ctx)
	// 调试模式下注册swagger接口文档路由
	if c.IsDebug() {
		registerSwagger(server)
	}

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
