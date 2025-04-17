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
)

func init() {
	flag.Parse()
	analyzer.Verbose = *verboseFlag
}

func main() {
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

	_, err := analyzer.AnalyzeDirectory(*pathFlag, options)
	fmt.Println()
	if err != nil {
		analyzer.PrintError("分析失败: %v", err)
		return
	}

	analyzer.PrintInfo("分析完成!")
}
