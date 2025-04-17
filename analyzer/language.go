package analyzer

import (
	"path/filepath"
	"strings"
)

// LanguageStats 存储每种语言的统计信息
type LanguageStats = Stat

// GetLanguageByExt 根据文件扩展名确定编程语言
func GetLanguageByExt(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if lang, ok := languageExt[ext]; ok {
		return lang
	}
	return "Unknown"
}

// 语言定义映射
var languageExt = map[string]string{
	".go":    "Go",
	".java":  "Java",
	".js":    "JavaScript",
	".ts":    "TypeScript",
	".jsx":   "React JSX",
	".tsx":   "React TSX",
	".py":    "Python",
	".rb":    "Ruby",
	".php":   "PHP",
	".c":     "C",
	".cpp":   "C++",
	".h":     "C/C++ Header",
	".hpp":   "C++ Header",
	".cs":    "C#",
	".swift": "Swift",
	".kt":    "Kotlin",
	".rs":    "Rust",
	".html":  "HTML",
	".css":   "CSS",
	".scss":  "SCSS",
	".less":  "LESS",
	".json":  "JSON",
	".xml":   "XML",
	".yaml":  "YAML",
	".yml":   "YAML",
	".md":    "Markdown",
	".txt":   "Text",
	".sh":    "Shell",
	".bat":   "Batch",
	".ps1":   "PowerShell",
	".sql":   "SQL",
	".r":     "R",
	".dart":  "Dart",
	".lua":   "Lua",
	".ex":    "Elixir",
	".exs":   "Elixir",
	".erl":   "Erlang",
	".hrl":   "Erlang",
	".clj":   "Clojure",
	".elm":   "Elm",
	".hs":    "Haskell",
	".pl":    "Perl",
	".pm":    "Perl",
}

// CommentPatterns 存储不同语言的注释样式
var CommentPatterns = map[string]struct {
	SingleLine []string
	MultiStart []string
	MultiEnd   []string
}{
	"Go": {
		SingleLine: []string{"//"},
		MultiStart: []string{"/*"},
		MultiEnd:   []string{"*/"},
	},
	"Java": {
		SingleLine: []string{"//"},
		MultiStart: []string{"/*"},
		MultiEnd:   []string{"*/"},
	},
	"JavaScript": {
		SingleLine: []string{"//"},
		MultiStart: []string{"/*"},
		MultiEnd:   []string{"*/"},
	},
	"TypeScript": {
		SingleLine: []string{"//"},
		MultiStart: []string{"/*"},
		MultiEnd:   []string{"*/"},
	},
	"Python": {
		SingleLine: []string{"#"},
		MultiStart: []string{"'''", "\"\"\""},
		MultiEnd:   []string{"'''", "\"\"\""},
	},
	"Ruby": {
		SingleLine: []string{"#"},
		MultiStart: []string{"=begin"},
		MultiEnd:   []string{"=end"},
	},
	"PHP": {
		SingleLine: []string{"//", "#"},
		MultiStart: []string{"/*"},
		MultiEnd:   []string{"*/"},
	},
	"C": {
		SingleLine: []string{"//"},
		MultiStart: []string{"/*"},
		MultiEnd:   []string{"*/"},
	},
	"C++": {
		SingleLine: []string{"//"},
		MultiStart: []string{"/*"},
		MultiEnd:   []string{"*/"},
	},
	"C#": {
		SingleLine: []string{"//"},
		MultiStart: []string{"/*"},
		MultiEnd:   []string{"*/"},
	},
	"Rust": {
		SingleLine: []string{"//"},
		MultiStart: []string{"/*"},
		MultiEnd:   []string{"*/"},
	},
	"HTML": {
		SingleLine: []string{},
		MultiStart: []string{"<!--"},
		MultiEnd:   []string{"-->"},
	},
	"CSS": {
		SingleLine: []string{},
		MultiStart: []string{"/*"},
		MultiEnd:   []string{"*/"},
	},
	"Shell": {
		SingleLine: []string{"#"},
		MultiStart: []string{},
		MultiEnd:   []string{},
	},
	"SQL": {
		SingleLine: []string{"--"},
		MultiStart: []string{"/*"},
		MultiEnd:   []string{"*/"},
	},
}
