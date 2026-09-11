// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"net/http"

	"oa.98ent.com/p9/platform-game/api/internal/config"
	"oa.98ent.com/p9/platform-game/api/internal/handler"
	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/common/response"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/platform_game.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	response.SetupHTTPX()

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
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
