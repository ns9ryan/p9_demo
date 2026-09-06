package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/rest"
	"oa.98ent.com/p9/platform-game/api/internal/config"
)

// RegisterSwaggerHandler 使用 go-zero 官方方式注册 Swagger 路由
func RegisterSwaggerHandler(server *rest.Server, cfg *config.Config) {
	// 服务 Swagger JSON spec
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/swagger.json",
		Handler: handleSwaggerJSON(cfg),
	})

	// 服务 Swagger UI HTML
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/swagger/index.html",
		Handler: handleSwaggerUI(),
	})

	// 重定向 /swagger 到 /swagger/index.html
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/swagger",
		Handler: handleSwaggerRedirect(),
	})
}

// handleSwaggerJSON 返回 Swagger 2.0 格式的 API spec
func handleSwaggerJSON(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取 swagger.json 文件路径
		swaggerFile := getSwaggerFilePath(cfg)

		content, err := os.ReadFile(swaggerFile)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "swagger.json not found: " + err.Error(),
			})
			return
		}

		// 验证 JSON 格式
		var spec map[string]interface{}
		if err := json.Unmarshal(content, &spec); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid swagger.json: " + err.Error(),
			})
			return
		}

		// 直接使用预生成的文件内容（已包含正确的 schemes）
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}
}

// handleSwaggerUI 返回 Swagger UI HTML 页面
func handleSwaggerUI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Platform Game API - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@4.15.5/swagger-ui.min.css">
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *,
        *:before,
        *:after {
            box-sizing: inherit;
        }
        body {
            margin: 0;
            padding: 0;
            background: #fafafa;
        }
        .swagger-ui .topbar {
            background-color: #2d5aa3;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    
    <!-- 使用官方预构建的 Swagger UI Bundle，避免 webpack 依赖问题 -->
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@4.15.5/swagger-ui-bundle.min.js" charset="UTF-8"></script>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@4.15.5/swagger-ui-standalone-preset.min.js" charset="UTF-8"></script>
    
    <script>
        // 等待 Swagger UI 库加载完成
        setTimeout(function() {
            if (typeof window.SwaggerUIBundle === 'undefined') {
                console.error('SwaggerUIBundle is not defined - CDN loading failed');
                document.getElementById('swagger-ui').innerHTML = '<div style="padding: 20px; color: red;">Failed to load Swagger UI. Please check CDN connectivity.</div>';
                return;
            }
            
            try {
                const ui = window.SwaggerUIBundle({
                    url: "/swagger.json",
                    dom_id: '#swagger-ui',
                    deepLinking: true,
                    presets: [
                        window.SwaggerUIBundle.presets.apis,
                        window.SwaggerUIStandalonePreset
                    ],
                    plugins: [
                        window.SwaggerUIBundle.plugins.DownloadUrl
                    ],
                    layout: "StandaloneLayout"
                });
                window.ui = ui;
                console.log("Swagger UI loaded successfully");
            } catch (error) {
                console.error("Error initializing Swagger UI:", error);
                document.getElementById('swagger-ui').innerHTML = '<div style="padding: 20px; color: red;">Error: ' + error.message + '</div>';
            }
        }, 500);  // 延迟 500ms 确保脚本加载完成
    </script>
</body>
</html>`

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}
}

// handleSwaggerRedirect 重定向 /swagger 到 /swagger/index.html
func handleSwaggerRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	}
}

// getSwaggerFilePath 获取 swagger.json 文件路径
// 根据 swagger_protocol 配置选择使用 http.json (development) 或 https.json (production)
func getSwaggerFilePath(cfg *config.Config) string {
	protocol := cfg.Swagger.Protocol
	if protocol == "" {
		protocol = "https" // 默认生产环境
	}
	fmt.Println("Current swagger_protocol:", protocol)

	// 选择文件名：http 使用 http.json，https 使用 https.json
	filename := "https.json"
	if protocol == "http" {
		filename = "http.json"
	}

	// 尝试多个可能的路径
	paths := []string{
		// 当前工作目录
		filepath.Join("api/internal/swagger", filename),
		// 相对于启动目录
		filepath.Join("./api/internal/swagger", filename),
		filepath.Join("../api/internal/swagger", filename),
		// deployment 常见位置: etc/swagger (Dockerfile 复制到 ./etc/swagger)
		filepath.Join("etc", "swagger", filename),
		filepath.Join("./etc", "swagger", filename),
		// 绝对路径方式
		filepath.Join(os.Getenv("PWD"), "api/internal/swagger", filename),
	}

	// 另外尝试以可执行文件目录为基准的路径（在容器中工作目录可能不同）
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		paths = append(paths,
			filepath.Join(exeDir, "api", "internal", "swagger", filename),
			filepath.Join(exeDir, "etc", "swagger", filename),
		)
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// 如果找不到，返回默认路径
	return filepath.Join("api/internal/swagger", filename)
}
