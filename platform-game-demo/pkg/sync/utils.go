package sync

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func JSONToString(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func StringToJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

func CompareName(i18nName, remoteName string) bool {
	var obj map[string]interface{}
	_ = StringToJSON(i18nName, &obj)
	if obj["default"] == nil {
		obj["default"] = ""
	}
	return obj["default"] != remoteName
}

//go:embed test/*.json
var syncTestFS embed.FS

func loadLocalJSON(fileName string, out any) error {
	if data, err := syncTestFS.ReadFile(filepath.ToSlash(filepath.Join("test", fileName))); err == nil {
		return json.Unmarshal(data, out)
	}

	candidates := []string{
		filepath.Join(".", "pkg", "sync", "test", fileName),
		filepath.Join(".", "sync", "test", fileName),
		filepath.Join(".", fileName),
	}

	if _, currentFile, _, ok := runtime.Caller(0); ok {
		candidates = append(candidates,
			filepath.Join(filepath.Dir(currentFile), "test", fileName),
			filepath.Join(filepath.Dir(currentFile), fileName),
		)
	}

	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			filepath.Join(wd, "pkg", "sync", "test", fileName),
			filepath.Join(wd, "sync", "test", fileName),
			filepath.Join(wd, fileName),
		)
	}

	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(data, out); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("load test JSON %s: file not found in embedded or local paths", fileName)
}
