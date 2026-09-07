package service

import (
	"context"
	"strings"

	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/ent"
	enti18n "oa.98ent.com/p9/core/rpc/ent/i18n"
	"oa.98ent.com/p9/core/rpc/model"
)

type CreateI18nReq struct {
	I18nGroup string
	TransKey  string
	Lang      string
	Value     string
}

type UpdateI18nReq struct {
	ID        int64
	I18nGroup *string
	TransKey  *string
	Lang      *string
	Value     *string
}

type I18nListReq struct {
	PageReq
	I18nGroup string
	TransKey  string
	Lang      string
}

type I18nItem struct {
	I18nGroup string
	TransKey  string
	Lang      string
	Value     string
}

type UpdateI18nByKeyReq struct {
	TransKey string
	Data     map[string]string
}

func (d *Deps) CreateI18n(ctx context.Context, req CreateI18nReq) (*model.I18n, error) {
	norm, err := normalizeI18n(req.I18nGroup, req.TransKey, req.Lang, req.Value)
	if err != nil {
		return nil, err
	}
	if err := d.requireSupportedLang(ctx, norm.Lang); err != nil {
		return nil, err
	}
	taken, err := d.i18nTaken(ctx, norm.TransKey, norm.Lang, 0)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, xerr.BadRequest(i18n.I18nExists)
	}
	row, err := d.Client.I18n.Create().
		SetI18nGroup(norm.I18nGroup).SetTransKey(norm.TransKey).
		SetLang(norm.Lang).SetValue(norm.Value).
		Save(ctx)
	if err != nil {
		return nil, xerr.BadRequest(i18n.I18nCreateFailed)
	}
	out := i18nFromEnt(row)
	return &out, nil
}

func (d *Deps) UpdateI18n(ctx context.Context, req UpdateI18nReq) error {
	row, err := d.i18nByID(ctx, req.ID)
	if err != nil {
		return err
	}
	next := i18nFromUpdate(row, req)
	norm, err := normalizeI18n(next.I18nGroup, next.TransKey, next.Lang, next.Value)
	if err != nil {
		return err
	}
	if norm.Lang != row.Lang {
		if err := d.requireSupportedLang(ctx, norm.Lang); err != nil {
			return err
		}
	}
	taken, err := d.i18nTaken(ctx, norm.TransKey, norm.Lang, row.ID)
	if err != nil {
		return err
	}
	if taken {
		return xerr.BadRequest(i18n.I18nExists)
	}
	return d.Client.I18n.UpdateOneID(row.ID).
		SetI18nGroup(norm.I18nGroup).SetTransKey(norm.TransKey).
		SetLang(norm.Lang).SetValue(norm.Value).
		Exec(ctx)
}

func (d *Deps) DeleteI18ns(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		if _, err := d.i18nByID(ctx, id); err != nil {
			return err
		}
		if err := d.Client.I18n.DeleteOneID(id).Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deps) ListI18ns(ctx context.Context, req I18nListReq) ([]model.I18n, int64, error) {
	q := d.Client.I18n.Query()
	if s := strings.TrimSpace(req.I18nGroup); s != "" {
		q.Where(enti18n.I18nGroupContains(s))
	}
	if s := strings.TrimSpace(req.TransKey); s != "" {
		q.Where(enti18n.TransKeyContains(s))
	}
	if s := strings.TrimSpace(req.Lang); s != "" {
		q.Where(enti18n.LangEQ(s))
	}
	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	req.normalizeNoLimit(50)
	list, err := q.Order(ent.Asc(enti18n.FieldI18nGroup), ent.Asc(enti18n.FieldTransKey), ent.Asc(enti18n.FieldLang), ent.Asc(enti18n.FieldID)).
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		All(ctx)
	return i18nsFromEnt(list), int64(total), err
}

func (d *Deps) GetI18nDict(ctx context.Context, group, lang string) (map[string]string, error) {
	group = strings.TrimSpace(group)
	lang = strings.TrimSpace(lang)
	if group == "" || lang == "" {
		return map[string]string{}, nil
	}
	list, err := d.Client.I18n.Query().
		Where(enti18n.I18nGroupEQ(group), enti18n.LangEQ(lang)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(list))
	for _, row := range list {
		out[row.TransKey] = row.Value
	}
	return out, nil
}

func (d *Deps) UpsertI18n(ctx context.Context, item I18nItem) error {
	norm, err := normalizeI18n(item.I18nGroup, item.TransKey, item.Lang, item.Value)
	if err != nil {
		return err
	}
	exist, err := d.Client.I18n.Query().
		Where(enti18n.TransKeyEQ(norm.TransKey), enti18n.LangEQ(norm.Lang)).
		Only(ctx)
	if ent.IsNotFound(err) {
		if err := d.requireSupportedLang(ctx, norm.Lang); err != nil {
			return err
		}
		_, err = d.Client.I18n.Create().
			SetI18nGroup(norm.I18nGroup).SetTransKey(norm.TransKey).
			SetLang(norm.Lang).SetValue(norm.Value).
			Save(ctx)
		return err
	}
	if err != nil {
		return err
	}
	return d.Client.I18n.UpdateOne(exist).SetI18nGroup(norm.I18nGroup).SetValue(norm.Value).Exec(ctx)
}

func (d *Deps) UpdateI18nByKey(ctx context.Context, req UpdateI18nByKeyReq) error {
	key := strings.TrimSpace(req.TransKey)
	if key == "" {
		return xerr.BadRequest(i18n.I18nTransKeyRequired)
	}
	if len(req.Data) == 0 {
		return xerr.BadRequest(i18n.I18nDataRequired)
	}
	group, err := d.i18nGroupByKey(ctx, key)
	if err != nil {
		return err
	}
	if group == "" {
		group = inferI18nGroup(key)
	}
	wrote := false
	for lang, value := range req.Data {
		lang = strings.TrimSpace(lang)
		if lang == "" {
			continue
		}
		wrote = true
		if err := d.upsertI18nLang(ctx, group, key, lang, strings.TrimSpace(value)); err != nil {
			return err
		}
	}
	if !wrote {
		return xerr.BadRequest(i18n.I18nDataRequired)
	}
	return nil
}

func (d *Deps) upsertI18nLang(ctx context.Context, group, key, lang, value string) error {
	exist, err := d.Client.I18n.Query().
		Where(enti18n.TransKeyEQ(key), enti18n.LangEQ(lang)).
		Only(ctx)
	if ent.IsNotFound(err) {
		if err := d.requireSupportedLang(ctx, lang); err != nil {
			return err
		}
		_, err = d.Client.I18n.Create().
			SetI18nGroup(group).SetTransKey(key).
			SetLang(lang).SetValue(value).
			Save(ctx)
		return err
	}
	if err != nil {
		return err
	}
	upd := d.Client.I18n.UpdateOne(exist).SetValue(value)
	if group != "" && exist.I18nGroup == "" {
		upd.SetI18nGroup(group)
	}
	return upd.Exec(ctx)
}

func (d *Deps) i18nGroupByKey(ctx context.Context, key string) (string, error) {
	row, err := d.Client.I18n.Query().
		Where(enti18n.TransKeyEQ(key), enti18n.I18nGroupNEQ("")).
		First(ctx)
	if ent.IsNotFound(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return row.I18nGroup, nil
}

func (d *Deps) i18nByID(ctx context.Context, id int64) (model.I18n, error) {
	row, err := d.Client.I18n.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return model.I18n{}, xerr.NotFound(i18n.I18nNotFound)
		}
		return model.I18n{}, err
	}
	return i18nFromEnt(row), nil
}

func (d *Deps) i18nTaken(ctx context.Context, key, lang string, exceptID int64) (bool, error) {
	q := d.Client.I18n.Query().Where(enti18n.TransKeyEQ(key), enti18n.LangEQ(lang))
	if exceptID > 0 {
		q = q.Where(enti18n.IDNEQ(exceptID))
	}
	return q.Exist(ctx)
}

type i18nNorm struct {
	I18nGroup string
	TransKey  string
	Lang      string
	Value     string
}

func normalizeI18n(group, key, lang, value string) (i18nNorm, error) {
	out := i18nNorm{
		I18nGroup: strings.TrimSpace(group),
		TransKey:  strings.TrimSpace(key),
		Lang:      strings.TrimSpace(lang),
		Value:     strings.TrimSpace(value),
	}
	out.TransKey = concatTransKey(out.I18nGroup, out.TransKey)
	if out.I18nGroup == "" || out.TransKey == "" || out.Lang == "" {
		return i18nNorm{}, xerr.BadRequest(i18n.I18nGroupKeyLangRequired)
	}
	return out, nil
}

func concatTransKey(group, key string) string {
	if group == "" || key == "" {
		return key
	}
	prefix := group + "."
	if strings.HasPrefix(key, prefix) {
		return key
	}
	return prefix + key
}

func inferI18nGroup(transKey string) string {
	i := strings.IndexByte(transKey, '.')
	if i <= 0 {
		return ""
	}
	return transKey[:i]
}

func i18nFromUpdate(row model.I18n, req UpdateI18nReq) model.I18n {
	if req.I18nGroup != nil {
		row.I18nGroup = *req.I18nGroup
	}
	if req.TransKey != nil {
		row.TransKey = *req.TransKey
	}
	if req.Lang != nil {
		row.Lang = *req.Lang
	}
	if req.Value != nil {
		row.Value = *req.Value
	}
	return row
}
