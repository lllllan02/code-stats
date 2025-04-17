package analyzer

type Stat struct {
	// 文件总数
	TotalFiles int // 文件总数

	// 文件大小
	TotalSize   int64   // 总大小
	AvgFileSize float64 // 平均文件大小

	// 字符数
	TotalChars int     // 总字符数
	AvgChars   float64 // 平均字符数

	// 行数
	TotalLines int     // 总行数
	AvgLines   float64 // 平均行数

	// 代码行数
	CodeLines    int     // 代码行数
	AvgCodeLines float64 // 平均代码行数

	// 注释行数
	CommentLines    int     // 注释行数
	AvgCommentLines float64 // 平均注释行数

	// 空白行数
	BlankLines    int     // 空白行数
	AvgBlankLines float64 // 平均空白行数

	// 代码密度
	CodeDensity    float64 // 代码密度: 代码行数/总行数
	CommentDensity float64 // 注释密度: 注释行数/总行数
	CommentRatio   float64 // 注释比例: 注释行数/代码行数

	// 平均每行字符数
	AvgLineLength float64 // 平均每行字符数
}

// 合并统计信息
func (s *Stat) Merge(other *Stat) {
	s.TotalFiles += other.TotalFiles
	s.TotalSize += other.TotalSize
	s.TotalChars += other.TotalChars
	s.TotalLines += other.TotalLines
	s.CodeLines += other.CodeLines
	s.CommentLines += other.CommentLines
	s.BlankLines += other.BlankLines
}

// 计算平均值
func (s *Stat) CalculateAvg() {
	s.AvgFileSize = float64(s.TotalSize) / float64(s.TotalFiles)
	s.AvgChars = float64(s.TotalChars) / float64(s.TotalLines)
	s.AvgLines = float64(s.TotalLines) / float64(s.TotalFiles)
	s.AvgCodeLines = float64(s.CodeLines) / float64(s.TotalLines)
	s.AvgCommentLines = float64(s.CommentLines) / float64(s.TotalLines)
	s.AvgBlankLines = float64(s.BlankLines) / float64(s.TotalLines)
	s.CodeDensity = float64(s.CodeLines) / float64(s.TotalLines)
	s.CommentDensity = float64(s.CommentLines) / float64(s.TotalLines)
	s.CommentRatio = float64(s.CommentLines) / float64(s.CodeLines)
	s.AvgLineLength = float64(s.TotalChars) / float64(s.TotalLines)
}
