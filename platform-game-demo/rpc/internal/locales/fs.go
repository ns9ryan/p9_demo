package locales

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/zeromicro/go-zero/core/logx"
)

func JSONToString(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func StringToJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

var syncTestFS embed.FS

func loadLocalJSON(fileName string, fileDir string, out any) error {
	fmt.Printf("Loading local JSON file: %s\n", filepath.ToSlash(filepath.Join(fileDir, fileName)))
	defer fmt.Printf("Finished attempting to load local JSON file: %s\n", filepath.ToSlash(filepath.Join(fileDir, fileName)))

	// 尝试从 embed.FS 读取
	if data, err := syncTestFS.ReadFile(filepath.ToSlash(filepath.Join(fileDir, fileName))); err == nil {
		return json.Unmarshal(data, out)
	}

	// 尝试多个候选路径
	candidates := []string{
		filepath.Join(".", "pkg", "sync", fileDir, fileName),
		filepath.Join(".", "sync", fileDir, fileName),
		filepath.Join(".", fileDir, fileName),
		filepath.Join(".", fileName),
	}

	if _, currentFile, _, ok := runtime.Caller(0); ok {
		candidates = append(candidates,
			filepath.Join(filepath.Dir(currentFile), fileDir, fileName),
			filepath.Join(filepath.Dir(currentFile), fileName),
		)
	}

	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			filepath.Join(wd, "pkg", "sync", fileDir, fileName),
			filepath.Join(wd, "sync", fileDir, fileName),
			filepath.Join(wd, fileDir, fileName),
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

func LoadLocalVendorRemoteJSON(fileName string, out any) error {
	return loadLocalJSON(fileName, "vendor_remote", out)
}

func MergeLocalJSON(fileName string, data map[string]string, biz string) error {
	logx.Infof("[名称映射] 开始加载并合并本地 JSON 文件: %s, 业务: %s", fileName, biz)

	existingData := make(map[string]string)
	_ = loadLocalJSON(fileName, filepath.Join("internal", "locales", "i18n"), &existingData)

	for k, v := range data {
		existingData[fmt.Sprintf("%s.%s", biz, k)] = v
	}

	dirPath := filepath.Join("internal", "locales", "i18n")

	// 创建目录（如果不存在）
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		logx.Errorf("[名称映射] 创建目录失败 %s, 业务: %s: %v", dirPath, biz, err)
		return fmt.Errorf("create directory failed: %w", err)
	}

	filePath := filepath.Join(dirPath, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		logx.Errorf("[名称映射] 创建本地 JSON 文件失败 %s: %v", filePath, err)
		return fmt.Errorf("create file failed: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	// 将合并后的数据写入本地 JSON 文件
	if err := encoder.Encode(existingData); err != nil {
		logx.Errorf("[名称映射] 写入 JSON 文件失败 %s: %v", filePath, err)
		return fmt.Errorf("encode JSON failed: %w", err)
	}

	logx.Infof("[名称映射] ✓ 加载并合并本地 JSON 文件成功: %s (共 %d 条记录)", filePath, len(existingData))
	return nil
}

func LoadLocalI18nNameMap(fileName string) (map[string]string, error) {
	data := make(map[string]string)
	err := loadLocalJSON(fileName, filepath.Join("internal", "locales", "i18n"), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
