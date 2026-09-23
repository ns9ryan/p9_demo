package service

import (
	"context"
	"strings"

	"oa.98ent.com/p9/common/i18n"
	"oa.98ent.com/p9/common/xerr"
	coreI18n "oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/rpc/ent"
	enti18n "oa.98ent.com/p9/core/rpc/ent/i18n"
	"oa.98ent.com/p9/core/rpc/model"
)

type CreateI18nReq struct {
	I18nCode  string
	I18nGroup string
	TransKey  string
	Lang      string
	Value     string
}

type UpdateI18nReq struct {
	ID        int64
	I18nCode  *string
	I18nGroup *string
	TransKey  *string
	Lang      *string
	Value     *string
}

type I18nListReq struct {
	PageReq
	I18nCode  string
	I18nGroup string
	TransKey  string
	Lang      string
}

type I18nItem struct {
	I18nCode  string
	I18nGroup string
	TransKey  string
	Lang      string
	Value     string
}

type I18nFileItem struct {
	I18nCode  string
	I18nGroup string
	I18nKey   string
	I18nValue string
}

type ImportI18nResult struct {
	Created int32
	Updated int32
	Skipped int32
}

type UpdateI18nByKeyReq struct {
	I18nCode  string
	I18nGroup string
	TransKey  string
	Data      map[string]string
}

func (d *Deps) CreateI18n(ctx context.Context, req CreateI18nReq) (*model.I18n, error) {
	norm, err := normalizeI18n(req.I18nCode, req.I18nGroup, req.TransKey, req.Lang, req.Value)
	if err != nil {
		return nil, err
	}
	if err := d.requireSupportedLang(ctx, norm.Lang); err != nil {
		return nil, err
	}
	taken, err := d.i18nTaken(ctx, norm.I18nCode, norm.TransKey, norm.Lang, 0)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, xerr.BadRequest(coreI18n.I18nExists)
	}
	row, err := d.Client.I18n.Create().
		SetI18nCode(norm.I18nCode).SetI18nGroup(norm.I18nGroup).SetTransKey(norm.TransKey).
		SetLang(norm.Lang).SetValue(norm.Value).
		Save(ctx)
	if err != nil {
		return nil, xerr.BadRequest(coreI18n.I18nCreateFailed)
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
	norm, err := normalizeI18n(next.I18nCode, next.I18nGroup, next.TransKey, next.Lang, next.Value)
	if err != nil {
		return err
	}
	if norm.Lang != row.Lang {
		if err := d.requireSupportedLang(ctx, norm.Lang); err != nil {
			return err
		}
	}
	taken, err := d.i18nTaken(ctx, norm.I18nCode, norm.TransKey, norm.Lang, row.ID)
	if err != nil {
		return err
	}
	if taken {
		return xerr.BadRequest(coreI18n.I18nExists)
	}
	return d.Client.I18n.UpdateOneID(row.ID).
		SetI18nCode(norm.I18nCode).SetI18nGroup(norm.I18nGroup).SetTransKey(norm.TransKey).
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

func (d *Deps) DeleteI18nsByKey(ctx context.Context, code, group, transKey string) error {
	key := strings.TrimSpace(transKey)
	if key == "" {
		return xerr.BadRequest(coreI18n.I18nTransKeyRequired)
	}
	group = strings.TrimSpace(group)
	if group == "" {
		group = inferI18nGroup(key)
	}
	if group == "" {
		return xerr.BadRequest(coreI18n.I18nGroupKeyLangRequired)
	}
	code = normalizeI18nCode(code)
	key = concatTransKey(group, key)
	n, err := d.Client.I18n.Delete().
		Where(enti18n.I18nCodeEQ(code), enti18n.I18nGroupEQ(group), enti18n.TransKeyEQ(key)).
		Exec(ctx)
	if err != nil {
		return err
	}
	if n == 0 {
		return xerr.NotFound(coreI18n.I18nNotFound)
	}
	return nil
}

func (d *Deps) ListI18ns(ctx context.Context, req I18nListReq) ([]model.I18n, int64, error) {
	q := d.Client.I18n.Query()
	if s := strings.TrimSpace(req.I18nCode); s != "" {
		q.Where(enti18n.I18nCodeEQ(s))
	}
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
	req.normalize(50)
	list, err := q.Order(ent.Asc(enti18n.FieldI18nCode), ent.Asc(enti18n.FieldI18nGroup), ent.Asc(enti18n.FieldTransKey), ent.Asc(enti18n.FieldLang), ent.Asc(enti18n.FieldID)).
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		All(ctx)
	return i18nsFromEnt(list), int64(total), err
}

func (d *Deps) GetI18nDict(ctx context.Context, code, group, lang string) (map[string]string, error) {
	code = strings.TrimSpace(code)
	group = strings.TrimSpace(group)
	lang = strings.TrimSpace(lang)
	if lang == "" {
		return map[string]string{}, nil
	}
	q := d.Client.I18n.Query().Where(enti18n.LangEQ(lang))
	if group != "" {
		q.Where(enti18n.I18nGroupEQ(group))
	}
	if code != "" {
		q.Where(enti18n.I18nCodeEQ(code))
	}
	list, err := q.Order(ent.Asc(enti18n.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(list))
	for _, row := range list {
		out[row.TransKey] = row.Value
	}
	return out, nil
}

func (d *Deps) ExportI18ns(ctx context.Context, lang, code, group string) ([]I18nFileItem, error) {
	lang = strings.TrimSpace(lang)
	if lang == "" {
		return nil, xerr.BadRequest(coreI18n.I18nLangRequired)
	}
	code = strings.TrimSpace(code)
	group = strings.TrimSpace(group)
	q := d.Client.I18n.Query().Where(enti18n.LangEQ(lang))
	if code != "" {
		q.Where(enti18n.I18nCodeEQ(code))
	}
	if group != "" {
		q.Where(enti18n.I18nGroupEQ(group))
	}
	list, err := q.Order(ent.Asc(enti18n.FieldI18nCode), ent.Asc(enti18n.FieldI18nGroup), ent.Asc(enti18n.FieldTransKey), ent.Asc(enti18n.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]I18nFileItem, 0, len(list))
	for _, row := range list {
		out = append(out, I18nFileItem{
			I18nCode:  row.I18nCode,
			I18nGroup: row.I18nGroup,
			I18nKey:   row.TransKey,
			I18nValue: row.Value,
		})
	}
	return out, nil
}

func (d *Deps) ImportI18ns(ctx context.Context, lang string, items []I18nFileItem) (ImportI18nResult, error) {
	lang = strings.TrimSpace(lang)
	if lang == "" {
		return ImportI18nResult{}, xerr.BadRequest(coreI18n.I18nLangRequired)
	}
	if len(items) == 0 {
		return ImportI18nResult{}, xerr.BadRequest(coreI18n.I18nDataRequired)
	}
	if err := d.requireSupportedLang(ctx, lang); err != nil {
		return ImportI18nResult{}, err
	}
	var out ImportI18nResult
	err := d.withTx(ctx, func(tx *Deps) error {
		valid := 0
		for _, it := range items {
			key := strings.TrimSpace(it.I18nKey)
			if key == "" {
				out.Skipped++
				continue
			}
			group := strings.TrimSpace(it.I18nGroup)
			if group == "" {
				group = inferI18nGroup(key)
			}
			if group == "" {
				out.Skipped++
				continue
			}
			code := normalizeI18nCode(it.I18nCode)
			key = concatTransKey(group, key)
			created, err := tx.upsertI18nLangCounted(ctx, code, group, key, lang, strings.TrimSpace(it.I18nValue))
			if err != nil {
				return err
			}
			valid++
			if created {
				out.Created++
			} else {
				out.Updated++
			}
		}
		if valid == 0 {
			return xerr.BadRequest(coreI18n.I18nDataRequired)
		}
		return nil
	})
	if err != nil {
		return ImportI18nResult{}, err
	}
	return out, nil
}

func (d *Deps) UpsertI18n(ctx context.Context, item I18nItem) error {
	return d.upsertI18n(ctx, item, false)
}

func (d *Deps) upsertI18n(ctx context.Context, item I18nItem, insertOnly bool) error {
	norm, err := normalizeI18n(item.I18nCode, item.I18nGroup, item.TransKey, item.Lang, item.Value)
	if err != nil {
		return err
	}
	exist, err := d.Client.I18n.Query().
		Where(enti18n.I18nCodeEQ(norm.I18nCode), enti18n.TransKeyEQ(norm.TransKey), enti18n.LangEQ(norm.Lang)).
		Only(ctx)
	if ent.IsNotFound(err) {
		if err := d.requireSupportedLang(ctx, norm.Lang); err != nil {
			return err
		}
		_, err = d.Client.I18n.Create().
			SetI18nCode(norm.I18nCode).SetI18nGroup(norm.I18nGroup).SetTransKey(norm.TransKey).
			SetLang(norm.Lang).SetValue(norm.Value).
			Save(ctx)
		return err
	}
	if err != nil {
		return err
	}
	if insertOnly {
		return nil
	}
	return d.Client.I18n.UpdateOne(exist).SetI18nGroup(norm.I18nGroup).SetValue(norm.Value).Exec(ctx)
}

func (d *Deps) UpdateI18nByKey(ctx context.Context, req UpdateI18nByKeyReq) error {
	key := strings.TrimSpace(req.TransKey)
	if key == "" {
		return xerr.BadRequest(coreI18n.I18nTransKeyRequired)
	}
	if len(req.Data) == 0 {
		return xerr.BadRequest(coreI18n.I18nDataRequired)
	}
	code := strings.TrimSpace(req.I18nCode)
	group := strings.TrimSpace(req.I18nGroup)
	targets, err := d.i18nByKeyTargets(ctx, code, group, key)
	if err != nil {
		return err
	}
	wrote := false
	for lang, value := range req.Data {
		lang = strings.TrimSpace(lang)
		if lang == "" {
			continue
		}
		wrote = true
		val := strings.TrimSpace(value)
		for _, t := range targets {
			if err := d.upsertI18nLang(ctx, t.code, t.group, key, lang, val); err != nil {
				return err
			}
		}
	}
	if !wrote {
		return xerr.BadRequest(coreI18n.I18nDataRequired)
	}
	return nil
}

type i18nKeyTarget struct {
	code  string
	group string
}

func (d *Deps) i18nByKeyTargets(ctx context.Context, code, group, key string) ([]i18nKeyTarget, error) {
	q := d.Client.I18n.Query().Where(enti18n.TransKeyEQ(key))
	if code != "" {
		q.Where(enti18n.I18nCodeEQ(code))
	}
	if group != "" {
		q.Where(enti18n.I18nGroupEQ(group))
	}
	list, err := q.Order(ent.Asc(enti18n.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	var out []i18nKeyTarget
	seen := map[string]struct{}{}
	for _, row := range list {
		if _, ok := seen[row.I18nCode]; ok {
			continue
		}
		seen[row.I18nCode] = struct{}{}
		g := row.I18nGroup
		if group != "" {
			g = group
		}
		out = append(out, i18nKeyTarget{code: row.I18nCode, group: g})
	}
	if len(out) > 0 {
		return out, nil
	}
	createCode := code
	if createCode == "" {
		createCode = i18n.CodePlatform
	}
	createGroup := group
	if createGroup == "" {
		createGroup = inferI18nGroup(key)
	}
	return []i18nKeyTarget{{code: createCode, group: createGroup}}, nil
}

func (d *Deps) upsertI18nLang(ctx context.Context, code, group, key, lang, value string) error {
	_, err := d.upsertI18nLangCounted(ctx, code, group, key, lang, value)
	return err
}

func (d *Deps) upsertI18nLangCounted(ctx context.Context, code, group, key, lang, value string) (bool, error) {
	exist, err := d.Client.I18n.Query().
		Where(enti18n.I18nCodeEQ(code), enti18n.TransKeyEQ(key), enti18n.LangEQ(lang)).
		Only(ctx)
	if ent.IsNotFound(err) {
		if err := d.requireSupportedLang(ctx, lang); err != nil {
			return false, err
		}
		_, err = d.Client.I18n.Create().
			SetI18nCode(code).SetI18nGroup(group).SetTransKey(key).
			SetLang(lang).SetValue(value).
			Save(ctx)
		return true, err
	}
	if err != nil {
		return false, err
	}
	return false, d.Client.I18n.UpdateOne(exist).SetValue(value).Exec(ctx)
}

func (d *Deps) i18nByID(ctx context.Context, id int64) (model.I18n, error) {
	row, err := d.Client.I18n.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return model.I18n{}, xerr.NotFound(coreI18n.I18nNotFound)
		}
		return model.I18n{}, err
	}
	return i18nFromEnt(row), nil
}

func (d *Deps) i18nTaken(ctx context.Context, code, key, lang string, exceptID int64) (bool, error) {
	q := d.Client.I18n.Query().Where(enti18n.I18nCodeEQ(code), enti18n.TransKeyEQ(key), enti18n.LangEQ(lang))
	if exceptID > 0 {
		q = q.Where(enti18n.IDNEQ(exceptID))
	}
	return q.Exist(ctx)
}

type i18nNorm struct {
	I18nCode  string
	I18nGroup string
	TransKey  string
	Lang      string
	Value     string
}

func normalizeI18n(code, group, key, lang, value string) (i18nNorm, error) {
	out := i18nNorm{
		I18nCode:  normalizeI18nCode(code),
		I18nGroup: strings.TrimSpace(group),
		TransKey:  strings.TrimSpace(key),
		Lang:      strings.TrimSpace(lang),
		Value:     strings.TrimSpace(value),
	}
	out.TransKey = concatTransKey(out.I18nGroup, out.TransKey)
	if out.I18nGroup == "" || out.TransKey == "" || out.Lang == "" {
		return i18nNorm{}, xerr.BadRequest(coreI18n.I18nGroupKeyLangRequired)
	}
	return out, nil
}

func normalizeI18nCode(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return i18n.CodePlatform
	}
	return code
}

func concatTransKey(group, key string) string {
	if group == "" || key == "" {
		return key
	}
	prefix := group + "."
	if strings.HasPrefix(key, prefix) || group == i18n.GroupError {
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
	if req.I18nCode != nil {
		row.I18nCode = *req.I18nCode
	}
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
