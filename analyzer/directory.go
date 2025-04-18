package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/samber/lo"
)

// DirectoryAnalyzerOptions 配置目录分析器的选项
type DirectoryAnalyzerOptions struct {
	ExcludeDirs []string // 排除的目录
	ExcludeExt  []string // 排除的文件扩展名
	MaxWorkers  int      // 最大并发数
	FollowLinks bool     // 是否跟踪符号链接
}

type DirectoryStats struct {
	*Stat

	Path           string
	FileStats      []*FileStats
	LanguageStats  map[string]*LanguageStats
	ExtensionStats map[string]*ExtensionStats
}

func AnalyzeDirectory(path string, options DirectoryAnalyzerOptions) (*DirectoryStats, error) {
	res := &DirectoryStats{
		Stat:           &Stat{},
		Path:           path,
		FileStats:      make([]*FileStats, 0),
		LanguageStats:  make(map[string]*LanguageStats),
		ExtensionStats: make(map[string]*ExtensionStats),
	}

	// 检查目录是否存在
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return res, fmt.Errorf("目录不存在: %s", path)
	} else if !info.IsDir() {
		return res, fmt.Errorf("不是目录: %s", path)
	}

	// 遍历目录
	var filePaths []string
	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			PrintError("无法访问: %s (%v)", path, err)
			return nil // 继续执行
		}

		// 跳过目录
		if info.IsDir() {
			baseName := filepath.Base(path)
			if slices.Contains(options.ExcludeDirs, baseName) {
				PrintInfo("已跳过目录: %s", path)
				return filepath.SkipDir
			}
			return nil
		}

		// 跳过符号链接
		if info.Mode()&os.ModeSymlink != 0 && !options.FollowLinks {
			PrintInfo("已跳过链接: %s", path)
			return nil
		}

		// 跳过指定扩展名
		ext := strings.ToLower(filepath.Ext(path))
		if slices.Contains(options.ExcludeExt, ext) {
			PrintInfo("已跳过文件: %s", path)
			return nil
		}

		filePaths = append(filePaths, path)
		return nil
	}); err != nil {
		PrintError("目录遍历失败: %v", err)
		return res, err
	}

	// 计算文件数量
	totalFiles := len(filePaths)
	if totalFiles == 0 {
		PrintWarning("目录中没有找到符合条件的文件")
		return res, nil
	}

	// 创建工作池
	var (
		wg           sync.WaitGroup
		mutex        sync.Mutex
		maxWorkers   = lo.Ternary(options.MaxWorkers > 0, options.MaxWorkers, 4)
		fileChan     = make(chan string, totalFiles)
		processedCnt = 0
	)

	// 获取全局进度条
	bar := GetGlobalProgressBar(totalFiles, "分析文件")

	// 启动工作池
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range fileChan {
				stats, err := AnalyzeFile(path)
				if err != nil {
					PrintError("分析失败: %s (%v)", path, err)
					continue
				}

				mutex.Lock()
				res.FileStats = append(res.FileStats, stats)
				processedCnt++
				// 更新进度条
				_ = bar.Set(processedCnt)
				mutex.Unlock()
			}
		}()
	}

	// 发送文件到通道
	for _, path := range filePaths {
		fileChan <- path
	}
	close(fileChan)

	// 等待所有工作完成
	wg.Wait()

	for _, fs := range res.FileStats {
		// 汇总统计
		res.Merge(fs.Stat)

		// 语言统计
		lang := fs.Language
		if _, exists := res.LanguageStats[lang]; !exists {
			res.LanguageStats[lang] = &Stat{}
		}
		res.LanguageStats[lang].Merge(fs.Stat)

		// 文件扩展名统计
		ext := strings.ToLower(filepath.Ext(fs.Path))
		if _, exists := res.ExtensionStats[ext]; !exists {
			res.ExtensionStats[ext] = &Stat{}
		}
		res.ExtensionStats[ext].Merge(fs.Stat)
	}

	res.CalculateAvg()
	for _, lang := range res.LanguageStats {
		lang.CalculateAvg()
	}
	for _, ext := range res.ExtensionStats {
		ext.CalculateAvg()
	}
	return res, nil
}

// 默认排除的目录
var defaultExcludeDirs = []string{
	".git", "node_modules", "vendor", "dist", "build",
	"bin", "obj", "target", "out", "tmp", "temp",
	".idea", ".vscode", ".vs", ".github", ".gitlab",
}

// 默认排除的文件扩展名
var defaultExcludeExt = []string{
	".exe", ".dll", ".so", ".dylib", ".o", ".obj", ".a", ".lib",
	".jar", ".war", ".ear", ".zip", ".tar", ".gz", ".rar",
	".jpg", ".jpeg", ".png", ".gif", ".bmp", ".ico", ".svg",
	".mp3", ".mp4", ".avi", ".mkv", ".wav", ".flac",
	".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
	".lock", ".sum", ".mod", ".toml", ".map",
}

func DefaultOptions() DirectoryAnalyzerOptions {
	return DirectoryAnalyzerOptions{
		ExcludeDirs: defaultExcludeDirs,
		ExcludeExt:  defaultExcludeExt,
		MaxWorkers:  4,
		FollowLinks: false,
	}
}
