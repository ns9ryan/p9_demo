package svc

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/config"
	"oa.98ent.com/p9/core/common/coreadapt"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/middleware"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	Core      coreclient.Core
	Authority rest.Middleware
	Jwt       rest.Middleware
	ActionLog rest.Middleware
	ErrorLog  rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	cli := zrpc.MustNewClient(c.CoreRpc)
	coreCli := coreclient.NewCore(cli)
	auth := coreadapt.Auth(coreCli)
	i18n.SetDictLoader(func(ctx context.Context, group, lang string) (map[string]string, error) {
		resp, err := coreCli.GetI18NDict(ctx, &coreclient.GetI18NDictReq{I18NGroup: group, Lang: lang})
		if err != nil {
			return nil, err
		}
		items := resp.GetItems()
		if items == nil {
			items = map[string]string{}
		}
		return items, nil
	})
	return &ServiceContext{
		Config:    c,
		Core:      coreCli,
		Jwt:       middleware.JWT(auth),
		Authority: middleware.Authority(auth),
		ActionLog: middleware.ActionLog(coreadapt.ActionRecorder(coreCli)),
		ErrorLog:  middleware.ErrorLog(c.Name, coreadapt.ErrorRecorder(coreCli)),
	}
}
