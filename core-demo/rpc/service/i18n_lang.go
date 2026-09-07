package service

import (
	"context"
	"strings"

	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/ent"
	enti18n "oa.98ent.com/p9/core/rpc/ent/i18n"
	enti18nlang "oa.98ent.com/p9/core/rpc/ent/i18nlang"
	"oa.98ent.com/p9/core/rpc/model"
)

type CreateI18nLangReq struct {
	Lang      string
	Name      string
	Disabled  int16
	IsDefault int16
}

type UpdateI18nLangReq struct {
	ID        int64
	Lang      *string
	Name      *string
	Disabled  *int16
	IsDefault *int16
}

type I18nLangListReq struct {
	PageReq
	Lang     string
	Disabled *int16
}

func (d *Deps) EnsureI18nLangs(ctx context.Context, seeds []CreateI18nLangReq) error {
	for _, s := range seeds {
		row, err := d.Client.I18nLang.Query().Where(enti18nlang.LangEQ(s.Lang)).Only(ctx)
		if err != nil {
			if !ent.IsNotFound(err) {
				return err
			}
			if _, err := d.Client.I18nLang.Create().
				SetLang(s.Lang).SetName(s.Name).SetDisabled(s.Disabled).SetIsDefault(s.IsDefault).
				Save(ctx); err != nil {
				return err
			}
			continue
		}
		if strings.TrimSpace(row.Name) == "" && s.Name != "" {
			if err := d.Client.I18nLang.UpdateOneID(row.ID).SetName(s.Name).Exec(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *Deps) CreateI18nLang(ctx context.Context, req CreateI18nLangReq) (*model.I18nLang, error) {
	lang, err := normalizeLangCode(req.Lang)
	if err != nil {
		return nil, err
	}
	name, err := normalizeLangName(req.Name, true)
	if err != nil {
		return nil, err
	}
	if err := validFlag01(req.Disabled); err != nil {
		return nil, err
	}
	if err := validFlag01(req.IsDefault); err != nil {
		return nil, err
	}
	if req.IsDefault == 1 && req.Disabled == 1 {
		return nil, xerr.BadRequest(i18n.I18nCannotDisableDefault)
	}
	taken, err := d.i18nLangTaken(ctx, lang, 0)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, xerr.BadRequest(i18n.I18nLangExists)
	}
	var out *model.I18nLang
	err = d.withTx(ctx, func(tx *Deps) error {
		if req.IsDefault == 1 {
			if err := tx.clearDefaultLang(ctx, 0); err != nil {
				return err
			}
		}
		row, err := tx.Client.I18nLang.Create().
			SetLang(lang).SetName(name).SetDisabled(req.Disabled).SetIsDefault(req.IsDefault).
			Save(ctx)
		if err != nil {
			return xerr.BadRequest(i18n.I18nLangCreateFailed)
		}
		m := i18nLangFromEnt(row)
		out = &m
		return nil
	})
	return out, err
}

func (d *Deps) UpdateI18nLang(ctx context.Context, req UpdateI18nLangReq) error {
	row, err := d.i18nLangByID(ctx, req.ID)
	if err != nil {
		return err
	}
	next := row
	if req.Lang != nil {
		lang, err := normalizeLangCode(*req.Lang)
		if err != nil {
			return err
		}
		if lang != row.Lang {
			used, err := d.i18nHasEntries(ctx, row.Lang)
			if err != nil {
				return err
			}
			if used {
				return xerr.BadRequest(i18n.I18nCannotChangeLangWithEntries)
			}
		}
		next.Lang = lang
	}
	if req.Name != nil {
		name, err := normalizeLangName(*req.Name, true)
		if err != nil {
			return err
		}
		next.Name = name
	}
	if req.Disabled != nil {
		if err := validFlag01(*req.Disabled); err != nil {
			return err
		}
		next.Disabled = *req.Disabled
	}
	if req.IsDefault != nil {
		if err := validFlag01(*req.IsDefault); err != nil {
			return err
		}
		next.IsDefault = *req.IsDefault
	}
	if next.IsDefault == 1 && next.Disabled == 1 {
		return xerr.BadRequest(i18n.I18nCannotDisableDefault)
	}
	if row.IsDefault == 1 && next.IsDefault == 0 {
		return xerr.BadRequest(i18n.I18nCannotDisableDefault)
	}
	if row.IsDefault == 1 && next.Disabled == 1 {
		return xerr.BadRequest(i18n.I18nCannotDisableDefault)
	}
	taken, err := d.i18nLangTaken(ctx, next.Lang, row.ID)
	if err != nil {
		return err
	}
	if taken {
		return xerr.BadRequest(i18n.I18nLangExists)
	}
	return d.withTx(ctx, func(tx *Deps) error {
		if next.IsDefault == 1 {
			if err := tx.clearDefaultLang(ctx, row.ID); err != nil {
				return err
			}
		}
		return tx.Client.I18nLang.UpdateOneID(row.ID).
			SetLang(next.Lang).SetName(next.Name).SetDisabled(next.Disabled).SetIsDefault(next.IsDefault).
			Exec(ctx)
	})
}

func (d *Deps) DeleteI18nLangs(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		row, err := d.i18nLangByID(ctx, id)
		if err != nil {
			return err
		}
		if row.IsDefault == 1 {
			return xerr.BadRequest(i18n.I18nCannotDeleteDefault)
		}
		used, err := d.i18nHasEntries(ctx, row.Lang)
		if err != nil {
			return err
		}
		if used {
			return xerr.BadRequest(i18n.I18nCannotDeleteLangWithEntries)
		}
		if err := d.Client.I18nLang.DeleteOneID(id).Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deps) ListI18nLangs(ctx context.Context, req I18nLangListReq) ([]model.I18nLang, int64, error) {
	q := d.Client.I18nLang.Query()
	if s := strings.TrimSpace(req.Lang); s != "" {
		q.Where(enti18nlang.LangContains(s))
	}
	if req.Disabled != nil {
		q.Where(enti18nlang.DisabledEQ(*req.Disabled))
	}
	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	req.normalizeNoLimit(50)
	list, err := q.Order(ent.Desc(enti18nlang.FieldIsDefault), ent.Asc(enti18nlang.FieldLang), ent.Asc(enti18nlang.FieldID)).
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		All(ctx)
	return i18nLangsFromEnt(list), int64(total), err
}

func (d *Deps) ListEnabledI18nLangs(ctx context.Context) ([]model.I18nLang, error) {
	list, err := d.Client.I18nLang.Query().
		Where(enti18nlang.DisabledEQ(0)).
		Order(ent.Desc(enti18nlang.FieldIsDefault), ent.Asc(enti18nlang.FieldLang), ent.Asc(enti18nlang.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return i18nLangsFromEnt(list), nil
}

func (d *Deps) i18nLangByID(ctx context.Context, id int64) (model.I18nLang, error) {
	row, err := d.Client.I18nLang.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return model.I18nLang{}, xerr.NotFound(i18n.I18nLangNotFound)
		}
		return model.I18nLang{}, err
	}
	return i18nLangFromEnt(row), nil
}

func (d *Deps) i18nLangTaken(ctx context.Context, lang string, exceptID int64) (bool, error) {
	q := d.Client.I18nLang.Query().Where(enti18nlang.LangEQ(lang))
	if exceptID > 0 {
		q = q.Where(enti18nlang.IDNEQ(exceptID))
	}
	return q.Exist(ctx)
}

func (d *Deps) i18nHasEntries(ctx context.Context, lang string) (bool, error) {
	return d.Client.I18n.Query().Where(enti18n.LangEQ(lang)).Exist(ctx)
}

func (d *Deps) requireSupportedLang(ctx context.Context, lang string) error {
	exist, err := d.Client.I18nLang.Query().Where(enti18nlang.LangEQ(lang)).Exist(ctx)
	if err != nil {
		return err
	}
	if !exist {
		return xerr.BadRequest(i18n.I18nLangNotSupported)
	}
	return nil
}

func (d *Deps) clearDefaultLang(ctx context.Context, exceptID int64) error {
	q := d.Client.I18nLang.Update().Where(enti18nlang.IsDefaultEQ(1))
	if exceptID > 0 {
		q = q.Where(enti18nlang.IDNEQ(exceptID))
	}
	return q.SetIsDefault(0).Exec(ctx)
}

func (d *Deps) withTx(ctx context.Context, fn func(*Deps) error) error {
	tx, err := d.Client.Tx(ctx)
	if err != nil {
		return err
	}
	nd := *d
	nd.Client = tx.Client()
	if err := fn(&nd); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func normalizeLangCode(lang string) (string, error) {
	lang = strings.TrimSpace(lang)
	if lang == "" {
		return "", xerr.BadRequest(i18n.I18nLangRequired)
	}
	return lang, nil
}

func normalizeLangName(name string, required bool) (string, error) {
	name = strings.TrimSpace(name)
	if required && name == "" {
		return "", xerr.BadRequest(i18n.I18nNameRequired)
	}
	return name, nil
}

func validFlag01(v int16) error {
	if v != 0 && v != 1 {
		return xerr.BadRequest(i18n.InvalidParam)
	}
	return nil
}
