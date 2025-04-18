package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/lllllan02/code-stats/analyzer"
)

var (
	// 分析的目录
	pathFlag = flag.String("path", ".", "Directory path to analyze")

	// 排除的目录
	excludeDirsFlag = flag.String("exclude-dirs", "", "Comma-separated list of directories to exclude")

	// 排除的文件扩展名
	excludeExtsFlag = flag.String("exclude-exts", "", "Comma-separated list of file extensions to exclude")

	// 最大并发数
	maxWorkersFlag = flag.Int("max-workers", 10, "Maximum number of concurrent workers")

	// 是否跟踪符号链接
	followLinksFlag = flag.Bool("follow-links", false, "Follow symbolic links")

	// 是否开启详细日志
	verboseFlag = flag.Bool("verbose", false, "Show verbose output")

	// 报告输出选项
	outputFlag = flag.String("output", "code-stats-report.html", "Output file path for the report")
	topNFlag   = flag.Int("top", 20, "Show top N files in report")

	// 帮助信息
	helpFlag = flag.Bool("help", false, "Show help message")
)

func init() {
	flag.Parse()
	analyzer.Verbose = *verboseFlag
}

// 打印使用说明
func printUsage() {
	fmt.Println("代码统计工具 - 分析代码行数、注释比例等指标")
	fmt.Println("\n使用方法:")
	fmt.Println("  code-stats [选项]")
	fmt.Println("\n选项:")
	flag.PrintDefaults()
	fmt.Println("\n示例:")
	fmt.Println("  # 分析当前目录")
	fmt.Println("  code-stats")
	fmt.Println("\n  # 分析指定目录并排除node_modules")
	fmt.Println("  code-stats -path=/path/to/code -exclude-dirs=node_modules,vendor")
	fmt.Println("\n  # 生成报告并保存到指定文件")
	fmt.Println("  code-stats -output=report.html")
	fmt.Println("\n  # 只显示前50个最大的文件")
	fmt.Println("  code-stats -top=50")
}

func main() {
	// 如果指定了帮助参数，则显示使用说明并退出
	if *helpFlag {
		printUsage()
		return
	}

	analyzer.PrintInfo("开始分析目录: %s", *pathFlag)

	options := analyzer.DefaultOptions()
	options.MaxWorkers = *maxWorkersFlag
	options.FollowLinks = *followLinksFlag
	if *excludeDirsFlag != "" {
		options.ExcludeDirs = strings.Split(*excludeDirsFlag, ",")
	}
	if *excludeExtsFlag != "" {
		options.ExcludeExt = strings.Split(*excludeExtsFlag, ",")
	}

	stats, err := analyzer.AnalyzeDirectory(*pathFlag, options)
	fmt.Println()
	if err != nil {
		analyzer.PrintError("分析失败: %v", err)
		return
	}
	analyzer.PrintInfo("分析完成!")

	// 设置报告数据 - 从默认值开始，然后覆盖需要的字段
	reportData := analyzer.DefaultReportData(stats)
	reportData.Options.OutputFile = *outputFlag
	reportData.TopN = *topNFlag

	// 生成 HTML 报告
	analyzer.PrintInfo("开始生成 HTML 报告...")
	report := analyzer.GenerateHTMLReport(reportData)

	// 保存报告到文件
	analyzer.PrintInfo("开始保存报告到文件...")
	if err := analyzer.SaveReportToFile(report, reportData.Options.OutputFile); err != nil {
		analyzer.PrintError("保存报告失败: %v", err)
		return
	}

	analyzer.PrintInfo("报告已生成: %s", reportData.Options.OutputFile)
}
