package analyzer

import (
	"fmt"
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

// PrintProgress 打印进度日志（彩色进度条）
func PrintProgress(current, total int) {
	percent := float64(current) / float64(total) * 100
	width := 30
	completed := int(float64(width) * float64(current) / float64(total))

	// 移除方括号边框，使用更干净的设计
	fmt.Printf("\r%s%s[进度]%s ",
		ColorBold, ColorLightBlue, ColorReset)

	// 使用单一色系的渐变，颜色更接近
	for i := 0; i < width; i++ {
		if i < completed {
			// 使用单一蓝色系，从深到浅
			switch {
			case i < width/5:
				fmt.Printf("\033[38;5;27m█") // 深蓝色
			case i < width*2/5:
				fmt.Printf("\033[38;5;33m█") // 蓝色
			case i < width*3/5:
				fmt.Printf("\033[38;5;39m█") // 中蓝色
			case i < width*4/5:
				fmt.Printf("\033[38;5;45m█") // 亮蓝色
			default:
				fmt.Printf("\033[38;5;51m█") // 浅蓝色
			}
		} else {
			fmt.Printf("%s▒", ColorGray) // 使用更清晰的字符
		}
	}

	// 在进度条后显示数字信息，不使用括号
	fmt.Printf(" %s%d/%d%s %s%.1f%%%s",
		ColorBold, current, total, ColorReset,
		"\033[38;5;51m", percent, ColorReset) // 使用与进度条末端相同的颜色
}
