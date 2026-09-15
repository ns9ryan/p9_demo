package game_sync

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"oa.98ent.com/p9/platform-game/common/logger"
)

func JSONToString(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func StringToJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
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

func mergeLocalJSON(fileName string, data map[string]string, biz string) error {
	logger.Infof("[名称映射] 开始加载并合并本地 JSON 文件: %s, 业务: %s", fileName, biz)

	existingData := make(map[string]string)
	_ = loadLocalJSON(filepath.Join("name_map", fileName), &existingData)

	for k, v := range data {
		existingData[fmt.Sprintf("%s.%s", biz, k)] = v
	}

	dirPath := "name_map"

	// 创建目录（如果不存在）
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		logger.Errorf("[名称映射] 创建目录失败 %s, 业务: %s: %v", dirPath, biz, err)
		return fmt.Errorf("create directory failed: %w", err)
	}

	filePath := filepath.Join(dirPath, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		logger.Errorf("[名称映射] 创建本地 JSON 文件失败 %s: %v", filePath, err)
		return fmt.Errorf("create file failed: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	// 将合并后的数据写入本地 JSON 文件
	if err := encoder.Encode(existingData); err != nil {
		logger.Errorf("[名称映射] 写入 JSON 文件失败 %s: %v", filePath, err)
		return fmt.Errorf("encode JSON failed: %w", err)
	}

	logger.Infof("[名称映射] ✓ 加载并合并本地 JSON 文件成功: %s (共 %d 条记录)", filePath, len(existingData))
	return nil
}

func LoadLocalI18nNameMap(fileName string) (map[string]string, error) {
	data := make(map[string]string)
	err := loadLocalJSON(filepath.Join("name_map", fileName), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
