package analyzer

import (
	"fmt"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

// 颜色代码常量
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m" // 成功信息
	ColorYellow = "\033[33m" // 警告信息
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m" // 普通信息
	ColorGray   = "\033[90m" // 次要信息
	ColorBold   = "\033[1m"  // 粗体
	ColorItalic = "\033[3m"  // 斜体

	// 亮色系列
	ColorLightRed   = "\033[91m" // 错误信息
	ColorLightGreen = "\033[92m" // 进度条
	ColorLightBlue  = "\033[94m" // 信息前缀
	ColorLightCyan  = "\033[96m" // 主要信息
)

// Verbose 是否开启详细日志
var Verbose = false

// PrintInfo 打印信息类日志（青色）
func PrintInfo(format string, args ...interface{}) {
	if !Verbose {
		return
	}

	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s%s[info]%s %s%s%s\n", ColorBold, ColorLightBlue, ColorReset, ColorLightCyan, message, ColorReset)
}

// PrintError 打印错误类日志（亮红色）
func PrintError(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s%s[error]%s %s%s%s\n", ColorBold, ColorLightRed, ColorReset, ColorRed, message, ColorReset)
}

// PrintWarning 打印警告类日志（黄色）
func PrintWarning(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("%s%s[warning]%s %s%s%s\n", ColorBold, ColorYellow, ColorReset, ColorYellow, message, ColorReset)
}

// ProgressManager 进度条管理器，用于维护进度条实例
type ProgressManager struct {
	bars map[string]*progressbar.ProgressBar
	mu   sync.Mutex
}

// NewProgressManager 创建新的进度条管理器
func NewProgressManager() *ProgressManager {
	return &ProgressManager{
		bars: make(map[string]*progressbar.ProgressBar),
	}
}

// UpdateProgress 更新指定ID的进度条
// 如果进度条不存在则创建新的
func (pm *ProgressManager) UpdateProgress(id string, current, total int, description string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	bar, exists := pm.bars[id]
	if !exists {
		// 创建新的进度条
		bar = progressbar.NewOptions(total,
			progressbar.OptionSetDescription(description),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowCount(),
			progressbar.OptionShowIts(),
			progressbar.OptionFullWidth(),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[purple]█[reset]",
				SaucerHead:    "[purple]█[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
			progressbar.OptionThrottle(65*time.Millisecond), // 限制更新频率
			progressbar.OptionShowBytes(false),              // 不显示字节数
			progressbar.OptionSetRenderBlankState(true),     // 确保初始状态也渲染
			progressbar.OptionSetPredictTime(false),         // 不预测剩余时间
		)
		pm.bars[id] = bar
	}

	_ = bar.Set(current)

	// 如果完成了，可以从map中删除
	if current >= total {
		delete(pm.bars, id)
	}
}

// 全局进度条实例，用于单个进度显示的情况
var (
	globalBar     *progressbar.ProgressBar
	globalBarOnce sync.Once
	globalBarMu   sync.Mutex
)

// GetGlobalProgressBar 返回全局进度条实例（单例模式）
func GetGlobalProgressBar(total int, description string) *progressbar.ProgressBar {
	globalBarMu.Lock()
	defer globalBarMu.Unlock()

	// 如果已存在且总数相同，直接返回
	if globalBar != nil && int(globalBar.GetMax()) == total {
		return globalBar
	}

	// 创建新的全局进度条
	globalBar = progressbar.NewOptions(total,
		progressbar.OptionSetDescription(description),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[cyan]=[reset]",
			SaucerHead:    "[cyan]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionFullWidth(),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetPredictTime(false),
	)

	return globalBar
}
