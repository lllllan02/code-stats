package analyzer

import (
	"bufio"
	"os"
	"slices"
	"strings"
)

type FileStats struct {
	*Stat

	Path     string // 文件路径
	Language string // 语言
}

func AnalyzeFile(path string) (*FileStats, error) {
	res := &FileStats{
		Stat:     &Stat{TotalFiles: 1},
		Path:     path,
		Language: GetLanguageByExt(path),
	}

	// 分析文件大小
	if err := res.analyzeSize(path); err != nil {
		return res, err
	}

	// 分析文件内容
	if err := res.analyzeFile(path); err != nil {
		return res, err
	}

	return res, nil
}

func (f *FileStats) analyzeSize(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		PrintError("无法获取文件大小: %s (%v)", path, err)
		return err
	}
	f.TotalSize = info.Size()
	return nil
}

func (f *FileStats) analyzeFile(path string) error {
	// 获取当前语言的注释标记
	commentStyle, hasCommentStyle := CommentPatterns[f.Language]

	file, err := os.Open(path)
	if err != nil {
		PrintError("无法打开文件: %s (%v)", path, err)
		return err
	}
	defer file.Close()

	var inMultilineComment bool   // 是否在多行注释中
	var multilineEndMarker string // 多行注释结束标记
	for scanner := bufio.NewScanner(file); scanner.Scan(); {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		lineLength := len(line)
		f.TotalChars += lineLength

		// 统计行数
		f.TotalLines++

		// 统计空白行数
		if trimmedLine == "" {
			f.BlankLines++
			continue
		}

		// 如果找不到当前语言的注释样式，则全部视为代码行
		if !hasCommentStyle {
			f.CodeLines++
			continue
		}

		// 检查是否在多行注释中
		if inMultilineComment {
			f.CommentLines++
			// 检查多行注释结束
			for _, endMarker := range commentStyle.MultiEnd {
				if strings.Contains(trimmedLine, endMarker) && endMarker == multilineEndMarker {
					inMultilineComment = false
					multilineEndMarker = ""
					break
				}
			}
			continue
		}

		// 检查是否是单行注释
		if slices.ContainsFunc(commentStyle.SingleLine, func(marker string) bool {
			return strings.TrimSpace(strings.Split(trimmedLine, marker)[0]) == ""
		}) {
			f.CommentLines++
			continue
		}

		// 检查多行注释开始
		var isComment bool
		for i, startMarker := range commentStyle.MultiStart {
			if strings.Contains(trimmedLine, startMarker) {
				f.CommentLines++
				inMultilineComment = true

				// 确保有匹配的结束标记
				if i < len(commentStyle.MultiEnd) {
					multilineEndMarker = commentStyle.MultiEnd[i]
				}

				// 检查多行注释是否在同一行结束
				if strings.Contains(trimmedLine, multilineEndMarker) {
					inMultilineComment = false
					multilineEndMarker = ""
					break
				}

				isComment = true
				break
			}
		}

		// 如果既不是单行注释也不是多行注释，则视为代码行
		if !isComment {
			f.CodeLines++
		}
	}

	f.CalculateAvg()
	return nil
}
