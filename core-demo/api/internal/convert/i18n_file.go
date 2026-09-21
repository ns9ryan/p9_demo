package convert

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"unicode"

	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/coreclient"
)

// 限制上传大小 8MB
const I18nFileMaxBytes = 8 << 20

// 导出I18n文件名
func I18nExportFilename(lang string) string {
	lang = strings.TrimSpace(lang)
	if lang == "" {
		lang = "i18n"
	}
	var b strings.Builder
	b.WriteString("i18n-")
	for _, r := range lang {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.' {
			b.WriteRune(r)
		} else {
			b.WriteByte('_')
		}
	}
	b.WriteString(".json")
	return b.String()
}

// 将I18n导出响应转换为JSON
func MarshalI18nFile(in *types.ExportI18nResp) ([]byte, error) {
	if in == nil {
		in = &types.ExportI18nResp{}
	}
	if in.I18nItems == nil {
		in.I18nItems = []types.I18nFileItem{}
	}
	return json.MarshalIndent(in, "", "  ")
}

// 读取I18n上传文件
func ReadI18nUploadFile(r *http.Request) ([]byte, error) {
	if r == nil {
		return nil, xerr.BadRequest(i18n.I18nDataRequired)
	}
	if err := r.ParseMultipartForm(I18nFileMaxBytes); err != nil {
		return nil, xerr.BadRequest(i18n.InvalidParam)
	}
	f, _, err := r.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return nil, xerr.BadRequest(i18n.I18nDataRequired)
		}
		return nil, xerr.BadRequest(i18n.InvalidParam)
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, I18nFileMaxBytes+1))
	if err != nil {
		return nil, xerr.BadRequest(i18n.InvalidParam)
	}
	if int64(len(data)) > I18nFileMaxBytes {
		return nil, xerr.BadRequest(i18n.InvalidParam)
	}
	return data, nil
}

// 读取导入请求：文件 + 可选 form lang
func ReadI18nImportReq(r *http.Request) (*types.ImportI18nReq, error) {
	data, err := ReadI18nUploadFile(r)
	if err != nil {
		return nil, err
	}
	return &types.ImportI18nReq{
		File: data,
		Lang: strings.TrimSpace(r.FormValue("lang")),
	}, nil
}

// 校验请求 lang：空则用文件 lang；非空必须与文件 lang 相等
func CheckImportI18nLang(reqLang, fileLang string) error {
	reqLang = strings.TrimSpace(reqLang)
	if reqLang == "" {
		return nil
	}
	if reqLang != strings.TrimSpace(fileLang) {
		return xerr.BadRequest(i18n.I18nLangMismatch)
	}
	return nil
}

// 解析I18n上传文件
func UnmarshalI18nFile(data []byte) (*types.ExportI18nResp, error) {
	if len(data) == 0 {
		return nil, xerr.BadRequest(i18n.I18nDataRequired)
	}
	var out types.ExportI18nResp
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, xerr.BadRequest(i18n.InvalidParam)
	}
	if out.I18nItems == nil {
		out.I18nItems = []types.I18nFileItem{}
	}
	return &out, nil
}

// 转换I18n上传文件为RPC请求
func ImportI18nReq(in *types.ExportI18nResp) *coreclient.ImportI18NReq {
	if in == nil {
		return &coreclient.ImportI18NReq{}
	}
	items := make([]*coreclient.I18NFileItem, 0, len(in.I18nItems))
	for _, row := range in.I18nItems {
		items = append(items, &coreclient.I18NFileItem{
			I18NCode: row.I18nCode, I18NGroup: row.I18nGroup, I18NKey: row.I18nKey, I18NValue: row.I18nValue,
		})
	}
	return &coreclient.ImportI18NReq{Lang: in.Lang, I18NItems: items}
}
