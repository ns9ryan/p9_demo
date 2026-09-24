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
	logx.Infof("[loadLocalJSON] 📂 开始加载文件: %s (目录: %s)", fileName, fileDir)

	// 获取当前工作目录
	wd, _ := os.Getwd()
	logx.Infof("[loadLocalJSON] 📍 当前工作目录: %s", wd)

	fmt.Printf("Loading local JSON file: %s\n", filepath.ToSlash(filepath.Join(fileDir, fileName)))
	defer fmt.Printf("Finished attempting to load local JSON file: %s\n", filepath.ToSlash(filepath.Join(fileDir, fileName)))

	// 尝试从 embed.FS 读取
	embedPath := filepath.ToSlash(filepath.Join(fileDir, fileName))
	logx.Infof("[loadLocalJSON] 🔍 尝试从 embed.FS 读取: %s", embedPath)
	if data, err := syncTestFS.ReadFile(embedPath); err == nil {
		logx.Infof("[loadLocalJSON] ✅ 从 embed.FS 成功读取文件")
		return json.Unmarshal(data, out)
	} else {
		logx.Errorf("[loadLocalJSON] ⚠️  embed.FS 读取失败: %v", err)
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

	logx.Infof("[loadLocalJSON] 🔍 尝试本地路径搜索 (共 %d 个候选路径):", len(candidates))
	for i, path := range candidates {
		logx.Infof("[loadLocalJSON]    [%d/%d] %s", i+1, len(candidates), path)
		data, err := os.ReadFile(path)
		if err != nil {
			logx.Debugf("[loadLocalJSON]       ❌ 未找到: %v", err)
			continue
		}
		logx.Infof("[loadLocalJSON]       ✅ 文件已找到!")
		if err := json.Unmarshal(data, out); err != nil {
			logx.Errorf("[loadLocalJSON]       ❌ JSON 解析失败: %v", err)
			return err
		}
		logx.Infof("[loadLocalJSON] ✅ 成功加载文件: %s", path)
		return nil
	}

	logx.Errorf("[loadLocalJSON] ❌ 文件未找到: %s (目录: %s)", fileName, fileDir)
	logx.Errorf("[loadLocalJSON] 📋 已尝试的路径数: %d", len(candidates))
	return fmt.Errorf("load test JSON %s: file not found in embedded or local paths", fileName)
}

func LoadLocalVendorRemoteJSON(fileName string, out any) error {
	logx.Infof("[LoadLocalVendorRemoteJSON] 📥 加载 vendor_remote 文件: %s", fileName)
	// 注意：文件在 internal/locales/vendor_remote/ 目录中（Dockerfile 中复制的位置）
	err := loadLocalJSON(fileName, filepath.Join("internal", "locales", "vendor_remote"), out)
	if err != nil {
		logx.Errorf("[LoadLocalVendorRemoteJSON] ❌ 加载失败: %v", err)
	} else {
		logx.Infof("[LoadLocalVendorRemoteJSON] ✅ 加载成功")
	}
	return err
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
