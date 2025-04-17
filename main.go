package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/lllllan02/code-stats/analyzer"
)

var (
	pathFlag        = flag.String("path", ".", "Directory path to analyze")
	excludeDirsFlag = flag.String("exclude-dirs", "", "Comma-separated list of directories to exclude")
	excludeExtsFlag = flag.String("exclude-exts", "", "Comma-separated list of file extensions to exclude")
	maxWorkersFlag  = flag.Int("max-workers", 10, "Maximum number of concurrent workers")
	followLinksFlag = flag.Bool("follow-links", false, "Follow symbolic links")
	verboseFlag     = flag.Bool("verbose", false, "Show verbose output")
)

func main() {
	flag.Parse()

	options := analyzer.DefaultOptions()
	options.MaxWorkers = *maxWorkersFlag
	options.FollowLinks = *followLinksFlag
	if *excludeDirsFlag != "" {
		options.ExcludeDirs = strings.Split(*excludeDirsFlag, ",")
	}
	if *excludeExtsFlag != "" {
		options.ExcludeExt = strings.Split(*excludeExtsFlag, ",")
	}
	if *verboseFlag {
		options.ProgressFunc = func(current, total int) {
			fmt.Printf("进度: %d/%d (%.2f%%)\n", current, total, float64(current)/float64(total)*100)
		}
	}

	stats, err := analyzer.AnalyzeDirectory(*pathFlag, options)
	if err != nil {
		fmt.Printf("分析目录失败: %v\n", err)
		return
	}

	fmt.Printf("stats.Stat: %+v\n", stats.Stat)

	fmt.Printf("分析目录成功: %+v\n", stats)
}
