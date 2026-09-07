package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"oa.98ent.com/p9/core/common/coreadapt"
	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/entdb"
	"oa.98ent.com/p9/core/common/entmixin"
	"oa.98ent.com/p9/core/common/middleware"
	"oa.98ent.com/p9/core/common/response"
	"oa.98ent.com/p9/core/example/promo-api/ent"
	"oa.98ent.com/p9/core/example/promo-api/ent/intercept"
	"oa.98ent.com/p9/core/example/promo-api/ent/migrate"
	_ "oa.98ent.com/p9/core/example/promo-api/ent/runtime"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	DB struct {
		Driver string
		DSN    string
	}
	CoreRpc zrpc.RpcClientConf
}

var configFile = flag.String("f", "etc/promo-api.yaml", "config file")

func main() {
	flag.Parse()
	var c Config
	conf.MustLoad(*configFile, &c)
	response.SetupHTTPX()

	promo, err := openPromo(c.DB.Driver, c.DB.DSN)
	logx.Must(err)
	logx.Must(promo.Schema.Create(context.Background(), migrate.WithDropIndex(true), migrate.WithDropColumn(true)))
	seedPromos(promo)

	cli := zrpc.MustNewClient(c.CoreRpc)
	coreCli := coreclient.NewCore(cli)
	auth := coreadapt.Auth(coreCli)
	logx.Must(registerPromoCatalog(coreCli))

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	server.Use(middleware.I18n)
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{middleware.JWT(auth), middleware.Authority(auth)},
			rest.Route{Method: http.MethodGet, Path: "/promo/list", Handler: promoList(promo)},
		),
		rest.WithPrefix("/admin"),
	)
	logx.Infof("promo-api listening on %s:%d", c.Host, c.Port)
	server.Start()
}

func openPromo(driver, dsn string) (*ent.Client, error) {
	drv, err := entdb.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	client := ent.NewClient(ent.Driver(drv), ent.Debug())
	client.Intercept(intercept.TraverseFunc(func(ctx context.Context, q intercept.Query) error {
		entmixin.FilterOperatorCode(ctx, q)
		return nil
	}))
	return client, nil
}

func promoList(promo *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := promo.Promotion.Query().All(r.Context())
		if err != nil {
			response.FailCtx(r.Context(), w, err)
			return
		}
		all, _ := promo.Promotion.Query().All(ctxdata.SkipTenant(r.Context()))
		code := ""
		if c := ctxdata.ClaimsFromCtx(r.Context()); c != nil {
			code = c.OperatorCode
		}
		out := make([]map[string]any, 0, len(rows))
		for _, row := range rows {
			out = append(out, map[string]any{
				"id": row.ID, "title": row.Title, "operator_code": row.OperatorCode,
				"created_at": row.CreatedAt, "updated_at": row.UpdatedAt,
			})
		}
		response.OkCtx(r.Context(), w, map[string]any{
			"operator_code": code, "list": out, "scoped": len(rows), "unscoped_total": len(all),
		})
	}
}

const (
	menuTypeDir  int32 = 0
	menuTypeMenu int32 = 1
)

func registerPromoCatalog(cli coreclient.Core) error {
	req := &coreclient.RegisterCatalogReq{
		Menus: []*coreclient.RegisterMenuReq{
			{Name: "PromoCenter", Title: "menu.route.promoCenter", MenuType: menuTypeDir, Path: "/promo", Sort: 20},
			{Name: "PromoActivityList", Title: "menu.route.promoActivityList", MenuType: menuTypeMenu, Path: "/promo/activity/list", Component: "promo/activity/list", ParentName: "PromoCenter", Sort: 21},
		},
		Apis: []*coreclient.CreateApiReq{
			{Description: "api.promoList", ApiGroup: "promo", Method: http.MethodGet, Path: "/admin/promo/list", ServiceName: "promo-api"},
		},
		I18N: []*coreclient.I18NItem{
			{I18NGroup: "menu", TransKey: "menu.route.promoCenter", Lang: "zh-CN", Value: "优惠中心"},
			{I18NGroup: "menu", TransKey: "menu.route.promoCenter", Lang: "en-US", Value: "Promotions"},
			{I18NGroup: "menu", TransKey: "menu.route.promoActivityList", Lang: "zh-CN", Value: "活动列表"},
			{I18NGroup: "menu", TransKey: "menu.route.promoActivityList", Lang: "en-US", Value: "Activities"},
			{I18NGroup: "api", TransKey: "api.promoList", Lang: "zh-CN", Value: "活动列表"},
			{I18NGroup: "api", TransKey: "api.promoList", Lang: "zh-HK", Value: "活動列表"},
			{I18NGroup: "api", TransKey: "api.promoList", Lang: "en-US", Value: "Promotion list"},
		},
	}
	var last error
	for range 20 {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		_, last = cli.RegisterCatalog(ctx, req)
		cancel()
		if last == nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("register promo catalog: %w", last)
}

func seedPromos(promo *ent.Client) {
	ctx := ctxdata.SkipTenant(context.Background())
	n, err := promo.Promotion.Query().Count(ctx)
	if err != nil || n > 0 {
		return
	}
	_, _ = promo.Promotion.Create().SetTitle("本厅活动").SetOperatorCode("demo").Save(ctx)
	_, _ = promo.Promotion.Create().SetTitle("其他厅活动").SetOperatorCode("other").Save(ctx)
}
