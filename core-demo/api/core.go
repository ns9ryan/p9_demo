// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"net/http"

	"oa.98ent.com/p9/core/api/internal/catalog"
	"oa.98ent.com/p9/core/api/internal/config"
	"oa.98ent.com/p9/core/api/internal/handler"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/common/middleware"
	"oa.98ent.com/p9/core/common/response"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/core.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	response.SetupHTTPX()

	server := rest.MustNewServer(c.RestConf, rest.WithCors(c.CROSConf.Address))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	// 注册I18n中间件
	server.Use(middleware.I18n)
	// 注册错误日志中间件
	server.Use(ctx.ErrorLog)

	// 注册菜单、API目录、多语言数据
	logx.Must(catalog.Register(ctx.Core))
	// 注册API路由
	handler.RegisterHandlers(server, ctx)
	// 开发环境或测试环境注册swagger接口文档路由
	if c.Mode == service.DevMode || c.Mode == service.TestMode {
		registerSwagger(server)
	}

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

//go:embed swagger/core.json
var swaggerSpec []byte

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <title>Core API</title>
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

// registerSwagger 注册swagger接口文档路由
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
	})
}
