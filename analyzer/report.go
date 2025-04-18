package analyzer

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/template"
	"time"

	"github.com/spf13/cast"
)

// ReportData 包含报告所需的所有数据
type ReportData struct {
	TopN           int             // 显示的文件数量
	OutputFile     string          // 输出文件
	GenerationTime string          // 生成时间
	Stats          *DirectoryStats // 统计数据

	FileSizesLimit int
	FileLinesLimit int
	TopLanguages   []LanguageItem
	SortedExts     []ExtensionItem
	FilesBySize    []*FileStats
	FilesByLines   []*FileStats

	// Git 相关数据
	HasGitStats       bool              // 是否有 Git 统计信息
	TopContributors   []ContributorItem // 排名前N的贡献者
	ContributorsLimit int               // 贡献者数量限制
}

// ContributorItem 表示UI显示用的贡献者项
type ContributorItem struct {
	Name        string
	CommitCount int
}

// DefaultReportData 返回默认的报告数据
func DefaultReportData(stats *DirectoryStats) ReportData {
	return ReportData{
		Stats:          stats,
		OutputFile:     "code-stats-report.html",
		GenerationTime: time.Now().Format("2006-01-02 15:04:05"),
		TopN:           20,                    // 默认显示20个文件
		HasGitStats:    stats.GitStats != nil, // 是否有 Git 统计信息
	}
}

// LanguageItem 表示UI显示用的语言项
type LanguageItem struct {
	Name  string
	Stats *LanguageStats
}

// ExtensionItem 表示UI显示用的扩展名项
type ExtensionItem struct {
	Name  string
	Stats *ExtensionStats
}

//go:embed report.tpl
var htmlReportTemplate []byte

// GenerateHTMLReport 生成HTML格式的报告
func GenerateHTMLReport(data ReportData) string {
	// 创建模板函数映射
	funcMap := template.FuncMap{
		"divideBy": func(a, b interface{}) float64 {
			aFloat := cast.ToFloat64(a)
			bFloat := cast.ToFloat64(b)

			if bFloat == 0 {
				return 0
			}
			return aFloat / bFloat
		},
		"multiply": func(a, b float64) float64 {
			return a * b
		},
		"commentRatio": func(comments, code int) float64 {
			if code == 0 {
				return 0
			}
			return float64(comments) / float64(code)
		},
		"subtract": func(a, b int) int {
			return a - b
		},
		"add": func(a, b int) int {
			return a + b
		},
		"ext": filepath.Ext,
		"formatTime": func(t time.Time) string {
			if t.IsZero() {
				return "N/A"
			}
			return t.Format("2006-01-02 15:04:05")
		},
	}

	// 确保GenerationTime字段已设置
	if data.GenerationTime == "" {
		data.GenerationTime = time.Now().Format("2006-01-02 15:04:05")
	}

	stats := data.Stats
	// 处理语言数据
	if len(stats.LanguageStats) > 0 {
		langs := make([]LanguageItem, 0, len(stats.LanguageStats))
		for name, stat := range stats.LanguageStats {
			if name != "" {
				langs = append(langs, LanguageItem{name, stat})
			}
		}
		// 按代码行数排序
		sort.Slice(langs, func(i, j int) bool {
			return langs[i].Stats.CodeLines > langs[j].Stats.CodeLines
		})
		data.TopLanguages = langs
	}

	// 处理扩展名数据
	if len(stats.ExtensionStats) > 0 {
		exts := make([]ExtensionItem, 0, len(stats.ExtensionStats))
		for ext, stat := range stats.ExtensionStats {
			if ext == "" {
				ext = "(无扩展名)"
			}
			exts = append(exts, ExtensionItem{ext, stat})
		}
		// 按文件数排序
		sort.Slice(exts, func(i, j int) bool {
			return exts[i].Stats.TotalFiles > exts[j].Stats.TotalFiles
		})
		data.SortedExts = exts
	}

	// 处理文件数据
	if len(stats.FileStats) > 0 {
		// 使用用户指定的TopN值
		topN := data.TopN
		if topN <= 0 {
			topN = 20 // 默认值
		}

		// 按大小排序的文件
		filesBySize := make([]*FileStats, len(stats.FileStats))
		copy(filesBySize, stats.FileStats)
		sort.Slice(filesBySize, func(i, j int) bool {
			return filesBySize[i].TotalSize > filesBySize[j].TotalSize
		})
		limit := topN
		if limit > len(filesBySize) {
			limit = len(filesBySize)
		}
		data.FilesBySize = filesBySize[:limit]
		data.FileSizesLimit = limit

		// 按代码行排序的文件
		filesByLines := make([]*FileStats, len(stats.FileStats))
		copy(filesByLines, stats.FileStats)
		sort.Slice(filesByLines, func(i, j int) bool {
			return filesByLines[i].CodeLines > filesByLines[j].CodeLines
		})
		limit = topN
		if limit > len(filesByLines) {
			limit = len(filesByLines)
		}
		data.FilesByLines = filesByLines[:limit]
		data.FileLinesLimit = limit
	}

	// 处理 Git 数据
	if stats.GitStats != nil {
		data.HasGitStats = true

		// 处理贡献者数据
		if len(stats.GitStats.TopContributors) > 0 {
			contributors := make([]ContributorItem, 0, len(stats.GitStats.TopContributors))
			for name, count := range stats.GitStats.TopContributors {
				contributors = append(contributors, ContributorItem{name, count})
			}
			// 按提交数排序
			sort.Slice(contributors, func(i, j int) bool {
				return contributors[i].CommitCount > contributors[j].CommitCount
			})

			// 取前 N 个贡献者
			limit := data.TopN
			if limit > len(contributors) {
				limit = len(contributors)
			}
			data.TopContributors = contributors[:limit]
			data.ContributorsLimit = limit
		}
	}

	// 解析并执行模板
	tmpl, err := template.New("report").Funcs(funcMap).Parse(string(htmlReportTemplate))
	if err != nil {
		return fmt.Sprintf("模板解析错误: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Sprintf("模板执行错误: %v", err)
	}

	return buf.String()
}

// 保存报告到文件
func SaveReportToFile(content string, filePath string) error {
	// 创建目录（如果不存在）
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 写入文件
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}
