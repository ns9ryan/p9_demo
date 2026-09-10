package coreadapt

import (
	"context"

	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/rpc/coreclient"
)

// DictLoader 通过 Core RPC 加载 sys_i18n 词条。
func DictLoader(cli coreclient.Core) i18n.DictLoader {
	return func(ctx context.Context, code, group, lang string) (map[string]string, error) {
		resp, err := cli.GetI18NDict(ctx, &coreclient.GetI18NDictReq{I18NCode: code, I18NGroup: group, Lang: lang})
		if err != nil {
			return nil, err
		}
		items := resp.GetItems()
		if items == nil {
			items = map[string]string{}
		}
		return items, nil
	}
}

// SetDictLoader 将 Core RPC 词条加载器注册到 i18n。
func SetDictLoader(cli coreclient.Core) {
	i18n.SetDictLoader(DictLoader(cli))
}
