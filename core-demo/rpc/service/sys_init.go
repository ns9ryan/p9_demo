package service

import (
	"context"

	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/ent/sysinit"
)

const (
	InitKeyMenu      = "menu"
	InitKeyAPI       = "api"
	InitKeyI18nLangs = "i18n_langs"
	InitKeyI18nDict  = "i18n_dict"
)

func (d *Deps) hasInit(ctx context.Context, key string) (bool, error) {
	return d.Client.SysInit.Query().Where(sysinit.InitKeyEQ(key)).Exist(ctx)
}

func (d *Deps) markInit(ctx context.Context, key string) error {
	exist, err := d.hasInit(ctx, key)
	if err != nil || exist {
		return err
	}
	_, err = d.Client.SysInit.Create().SetInitKey(key).Save(ctx)
	if err != nil && ent.IsConstraintError(err) {
		return nil
	}
	return err
}
