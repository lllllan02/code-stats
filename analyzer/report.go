package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ReportOptions 配置报告生成的选项
type ReportOptions struct {
	ShowSummary    bool   // 是否显示总体摘要
	ShowLanguages  bool   // 是否显示各语言统计
	ShowExtensions bool   // 是否显示各扩展名统计
	ShowFiles      bool   // 是否显示文件列表
	TopN           int    // 显示前N个文件(按大小或行数)
	OutputFormat   string // 输出格式: "text", "html"
	OutputFile     string // 输出文件路径，为空则不输出到文件
}

// DefaultReportOptions 返回默认报告选项
func DefaultReportOptions() ReportOptions {
	return ReportOptions{
		ShowSummary:    true,
		ShowLanguages:  true,
		ShowExtensions: true,
		ShowFiles:      false,
		TopN:           10,
		OutputFormat:   "html",
		OutputFile:     "code-stats-report.html",
	}
}

// GenerateHTMLReport 生成HTML格式的报告
func GenerateHTMLReport(stats *DirectoryStats, options ReportOptions) string {
	var sb strings.Builder

	// HTML头部
	sb.WriteString(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>代码统计报告</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f9f9f9;
        }
        h1, h2, h3 {
            color: #2c3e50;
        }
        h1 {
            border-bottom: 2px solid #3498db;
            padding-bottom: 10px;
        }
        h2 {
            margin-top: 30px;
            border-bottom: 1px solid #bdc3c7;
            padding-bottom: 5px;
        }
        .summary {
            background-color: #fff;
            border-radius: 5px;
            padding: 20px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
            margin-bottom: 20px;
        }
        .summary-item {
            margin-bottom: 10px;
            display: flex;
        }
        .summary-label {
            font-weight: bold;
            width: 200px;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
            background-color: #fff;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
            border-radius: 5px;
            overflow: hidden;
        }
        th {
            background-color: #3498db;
            color: white;
            padding: 12px 15px;
            text-align: left;
        }
        td {
            padding: 10px 15px;
            border-bottom: 1px solid #ddd;
        }
        tr:hover {
            background-color: #f5f5f5;
        }
        tr:last-child td {
            border-bottom: none;
        }
        .chart-container {
            display: flex;
            justify-content: space-between;
            margin: 30px 0;
        }
        .chart {
            width: 48%;
            background-color: #fff;
            border-radius: 5px;
            padding: 20px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        .timestamp {
            color: #7f8c8d;
            font-style: italic;
            margin-bottom: 30px;
        }
        footer {
            margin-top: 50px;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #7f8c8d;
            font-size: 0.9em;
            text-align: center;
        }
    </style>
</head>
<body>
    <h1>代码统计报告</h1>
    <div class="timestamp">生成时间: ` + time.Now().Format("2006-01-02 15:04:05") + `</div>
`)

	// 总体摘要
	if options.ShowSummary {
		sb.WriteString(`
    <h2>总体摘要</h2>
    <div class="summary">
        <div class="summary-item"><span class="summary-label">总文件数:</span> ` + fmt.Sprintf("%d 个文件", stats.TotalFiles) + `</div>
        <div class="summary-item"><span class="summary-label">总代码量:</span> ` + fmt.Sprintf("%d 行 (%.2f MB)", stats.TotalLines, float64(stats.TotalSize)/(1024*1024)) + `</div>
        <div class="summary-item"><span class="summary-label">代码行数:</span> ` + fmt.Sprintf("%d 行 (%.1f%%)", stats.CodeLines, stats.CodeDensity*100) + `</div>
        <div class="summary-item"><span class="summary-label">注释行数:</span> ` + fmt.Sprintf("%d 行 (%.1f%%)", stats.CommentLines, stats.CommentDensity*100) + `</div>
        <div class="summary-item"><span class="summary-label">空白行数:</span> ` + fmt.Sprintf("%d 行 (%.1f%%)", stats.BlankLines, stats.AvgBlankLines*100) + `</div>
        <div class="summary-item"><span class="summary-label">注释比例:</span> ` + fmt.Sprintf("%.2f (注释行/代码行)", stats.CommentRatio) + `</div>
        <div class="summary-item"><span class="summary-label">平均文件大小:</span> ` + fmt.Sprintf("%.2f KB", stats.AvgFileSize/1024) + `</div>
        <div class="summary-item"><span class="summary-label">平均行长度:</span> ` + fmt.Sprintf("%.1f 字符/行", stats.AvgLineLength) + `</div>
    </div>

    <!-- 可视化图表 -->
    <div class="chart-container">
        <div class="chart">
            <h3>代码组成</h3>
            <div style="position: relative; height: 200px;">
                <canvas id="compositionChart"></canvas>
            </div>
        </div>
        <div class="chart">
            <h3>语言分布</h3>
            <div style="position: relative; height: 200px;">
                <canvas id="languageChart"></canvas>
            </div>
        </div>
    </div>
`)
	}

	// 语言统计
	if options.ShowLanguages && len(stats.LanguageStats) > 0 {
		sb.WriteString(`
    <h2>各语言统计</h2>
    <table>
        <thead>
            <tr>
                <th>语言</th>
                <th>文件数</th>
                <th>代码行</th>
                <th>注释行</th>
                <th>空白行</th>
                <th>注释比例</th>
                <th>平均行长度</th>
            </tr>
        </thead>
        <tbody>
`)

		// 按代码行数排序语言
		type langStat struct {
			name string
			stat *LanguageStats
		}
		langs := make([]langStat, 0, len(stats.LanguageStats))
		for name, stat := range stats.LanguageStats {
			if name != "" {
				langs = append(langs, langStat{name, stat})
			}
		}
		sort.Slice(langs, func(i, j int) bool {
			return langs[i].stat.CodeLines > langs[j].stat.CodeLines
		})

		for _, l := range langs {
			if l.name == "" {
				continue // 跳过未识别语言
			}
			s := l.stat
			sb.WriteString(`
            <tr>
                <td>` + l.name + `</td>
                <td>` + fmt.Sprintf("%d", s.TotalFiles) + `</td>
                <td>` + fmt.Sprintf("%d", s.CodeLines) + `</td>
                <td>` + fmt.Sprintf("%d", s.CommentLines) + `</td>
                <td>` + fmt.Sprintf("%d", s.BlankLines) + `</td>
                <td>` + fmt.Sprintf("%.2f", s.CommentRatio) + `</td>
                <td>` + fmt.Sprintf("%.1f", s.AvgLineLength) + `</td>
            </tr>`)
		}

		sb.WriteString(`
        </tbody>
    </table>
`)
	}

	// 扩展名统计
	if options.ShowExtensions && len(stats.ExtensionStats) > 0 {
		sb.WriteString(`
    <h2>各扩展名统计</h2>
    <table>
        <thead>
            <tr>
                <th>扩展名</th>
                <th>文件数</th>
                <th>代码行</th>
                <th>总行数</th>
                <th>平均大小(KB)</th>
            </tr>
        </thead>
        <tbody>
`)

		// 按文件数排序扩展名
		type extStat struct {
			ext  string
			stat *ExtensionStats
		}
		exts := make([]extStat, 0, len(stats.ExtensionStats))
		for ext, stat := range stats.ExtensionStats {
			if ext == "" {
				ext = "(无扩展名)"
			}
			exts = append(exts, extStat{ext, stat})
		}
		sort.Slice(exts, func(i, j int) bool {
			return exts[i].stat.TotalFiles > exts[j].stat.TotalFiles
		})

		for _, e := range exts {
			s := e.stat
			sb.WriteString(`
            <tr>
                <td>` + e.ext + `</td>
                <td>` + fmt.Sprintf("%d", s.TotalFiles) + `</td>
                <td>` + fmt.Sprintf("%d", s.CodeLines) + `</td>
                <td>` + fmt.Sprintf("%d", s.TotalLines) + `</td>
                <td>` + fmt.Sprintf("%.2f", s.AvgFileSize/1024) + `</td>
            </tr>`)
		}

		sb.WriteString(`
        </tbody>
    </table>
`)
	}

	// 文件统计
	if options.ShowFiles && len(stats.FileStats) > 0 && options.TopN > 0 {
		// 按大小排序
		sb.WriteString(`
    <h2>文件统计（按大小排序前` + fmt.Sprintf("%d", options.TopN) + `）</h2>
    <table>
        <thead>
            <tr>
                <th>文件路径</th>
                <th>大小(KB)</th>
                <th>总行数</th>
                <th>代码行</th>
                <th>注释行</th>
            </tr>
        </thead>
        <tbody>
`)

		// 复制并按大小排序文件
		filesBySize := make([]*FileStats, len(stats.FileStats))
		copy(filesBySize, stats.FileStats)
		sort.Slice(filesBySize, func(i, j int) bool {
			return filesBySize[i].TotalSize > filesBySize[j].TotalSize
		})

		limit := options.TopN
		if limit > len(filesBySize) {
			limit = len(filesBySize)
		}

		for i := 0; i < limit; i++ {
			f := filesBySize[i]
			path := f.Path
			sb.WriteString(`
            <tr>
                <td>` + path + `</td>
                <td>` + fmt.Sprintf("%.2f", float64(f.TotalSize)/1024) + `</td>
                <td>` + fmt.Sprintf("%d", f.TotalLines) + `</td>
                <td>` + fmt.Sprintf("%d", f.CodeLines) + `</td>
                <td>` + fmt.Sprintf("%d", f.CommentLines) + `</td>
            </tr>`)
		}

		sb.WriteString(`
        </tbody>
    </table>

    <h2>文件统计（按代码行数排序前` + fmt.Sprintf("%d", options.TopN) + `）</h2>
    <table>
        <thead>
            <tr>
                <th>文件路径</th>
                <th>代码行</th>
                <th>注释行</th>
                <th>空白行</th>
                <th>注释比例</th>
            </tr>
        </thead>
        <tbody>
`)

		// 按行数排序
		filesByLines := make([]*FileStats, len(stats.FileStats))
		copy(filesByLines, stats.FileStats)
		sort.Slice(filesByLines, func(i, j int) bool {
			return filesByLines[i].CodeLines > filesByLines[j].CodeLines
		})

		limit = options.TopN
		if limit > len(filesByLines) {
			limit = len(filesByLines)
		}

		for i := 0; i < limit; i++ {
			f := filesByLines[i]
			path := f.Path
			ratio := float64(0)
			if f.CodeLines > 0 {
				ratio = float64(f.CommentLines) / float64(f.CodeLines)
			}
			sb.WriteString(`
            <tr>
                <td>` + path + `</td>
                <td>` + fmt.Sprintf("%d", f.CodeLines) + `</td>
                <td>` + fmt.Sprintf("%d", f.CommentLines) + `</td>
                <td>` + fmt.Sprintf("%d", f.BlankLines) + `</td>
                <td>` + fmt.Sprintf("%.2f", ratio) + `</td>
            </tr>`)
		}

		sb.WriteString(`
        </tbody>
    </table>
`)
	}

	// JavaScript 图表代码
	sb.WriteString(`
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <script>
        // 代码组成图表
        const compositionCtx = document.getElementById('compositionChart').getContext('2d');
        new Chart(compositionCtx, {
            type: 'pie',
            data: {
                labels: ['代码行', '注释行', '空白行'],
                datasets: [{
                    data: [` + fmt.Sprintf("%d, %d, %d", stats.CodeLines, stats.CommentLines, stats.BlankLines) + `],
                    backgroundColor: [
                        'rgba(54, 162, 235, 0.7)',
                        'rgba(255, 205, 86, 0.7)',
                        'rgba(201, 203, 207, 0.7)'
                    ],
                    borderColor: [
                        'rgb(54, 162, 235)',
                        'rgb(255, 205, 86)',
                        'rgb(201, 203, 207)'
                    ],
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'bottom',
                    }
                }
            }
        });

        // 语言分布图表
        const languageCtx = document.getElementById('languageChart').getContext('2d');
        new Chart(languageCtx, {
            type: 'bar',
            data: {
                labels: [`)

	// 添加语言标签
	type langStat struct {
		name string
		stat *LanguageStats
	}
	langs := make([]langStat, 0, len(stats.LanguageStats))
	for name, stat := range stats.LanguageStats {
		if name != "" {
			langs = append(langs, langStat{name, stat})
		}
	}
	sort.Slice(langs, func(i, j int) bool {
		return langs[i].stat.CodeLines > langs[j].stat.CodeLines
	})

	// 最多显示前8种语言
	displayLimit := 8
	if len(langs) < displayLimit {
		displayLimit = len(langs)
	}

	// 语言标签
	for i := 0; i < displayLimit; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("'" + langs[i].name + "'")
	}

	sb.WriteString(`],
                datasets: [{
                    label: '代码行数',
                    data: [`)

	// 语言数据
	for i := 0; i < displayLimit; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%d", langs[i].stat.CodeLines))
	}

	sb.WriteString(`],
                    backgroundColor: 'rgba(54, 162, 235, 0.7)',
                    borderColor: 'rgb(54, 162, 235)',
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    y: {
                        beginAtZero: true
                    }
                },
                plugins: {
                    legend: {
                        display: false
                    }
                }
            }
        });
    </script>
</body>
</html>`)

	return sb.String()
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

	PrintInfo("报告已保存至: %s", filePath)
	return nil
}

// GenerateReport 根据选项生成报告
func GenerateReport(stats *DirectoryStats, options ReportOptions) (string, error) {
	content := GenerateHTMLReport(stats, options)

	if err := SaveReportToFile(content, options.OutputFile); err != nil {
		return content, err
	}

	return content, nil
}
