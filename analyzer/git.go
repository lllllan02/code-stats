package analyzer

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// GitStats 存储 Git 仓库的统计信息
type GitStats struct {
	// 基本统计
	CommitCount      int       // 提交总数
	ContributorCount int       // 贡献者数量
	FirstCommitDate  time.Time // 首次提交日期
	LastCommitDate   time.Time // 最后提交日期
	ActiveDays       int       // 活跃天数（有提交的天数）

	// 文件变更统计
	TotalAdditions   int // 添加的行数总计
	TotalDeletions   int // 删除的行数总计
	TotalFileChanges int // 文件变更总数

	// 贡献者统计
	TopContributors map[string]int // 贡献者提交次数统计

	// 分支统计
	BranchCount int             // 分支数量
	BranchList  map[string]bool // 分支列表
}

// AnalyzeGitRepo 分析 Git 仓库统计信息
func AnalyzeGitRepo(repoPath string) (*GitStats, error) {
	stats := &GitStats{
		TopContributors: make(map[string]int),
		BranchList:      make(map[string]bool),
	}

	// 检查是否是 Git 仓库
	if !isGitRepo(repoPath) {
		PrintWarning("目录不是 Git 仓库: %s", repoPath)
		return stats, nil
	}

	// 获取提交数量
	if count, err := getCommitCount(repoPath); err == nil {
		stats.CommitCount = count
	}

	// 获取贡献者数量和贡献者统计
	if contributors, err := getContributors(repoPath); err == nil {
		stats.ContributorCount = len(contributors)
		stats.TopContributors = contributors
	}

	// 获取提交时间范围
	if first, last, activeDays, err := getCommitTimeStats(repoPath); err == nil {
		stats.FirstCommitDate = first
		stats.LastCommitDate = last
		stats.ActiveDays = activeDays
	}

	// 获取文件变更统计
	if additions, deletions, fileChanges, err := getChangeStats(repoPath); err == nil {
		stats.TotalAdditions = additions
		stats.TotalDeletions = deletions
		stats.TotalFileChanges = fileChanges
	}

	// 获取分支统计
	if branches, count, err := getBranchStats(repoPath); err == nil {
		stats.BranchCount = count
		stats.BranchList = branches
	}

	return stats, nil
}

// 检查是否是 Git 仓库
func isGitRepo(path string) bool {
	cmd := exec.Command("git", "-C", path, "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// 获取提交数量
func getCommitCount(path string) (int, error) {
	cmd := exec.Command("git", "-C", path, "rev-list", "--count", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		PrintError("获取提交数量失败: %v", err)
		return 0, err
	}

	count, err := strconv.Atoi(strings.TrimSpace(out.String()))
	if err != nil {
		PrintError("解析提交数量失败: %v", err)
		return 0, err
	}

	return count, nil
}

// 获取贡献者数量和统计
func getContributors(path string) (map[string]int, error) {
	cmd := exec.Command("git", "-C", path, "shortlog", "-sn", "--all")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		PrintError("获取贡献者统计失败: %v", err)
		return nil, err
	}

	contributors := make(map[string]int)
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			count, err := strconv.Atoi(parts[0])
			if err != nil {
				PrintWarning("解析贡献者提交数量失败: %v", err)
				continue
			}
			name := strings.Join(parts[1:], " ")
			contributors[name] = count
		}
	}

	return contributors, nil
}

// 获取提交时间范围
func getCommitTimeStats(path string) (time.Time, time.Time, int, error) {
	// 获取首次提交日期
	firstCmd := exec.Command("git", "-C", path, "log", "--reverse", "--format=%at", "--max-count=1")
	var firstOut bytes.Buffer
	firstCmd.Stdout = &firstOut
	if err := firstCmd.Run(); err != nil {
		PrintError("获取首次提交日期失败: %v", err)
		return time.Time{}, time.Time{}, 0, err
	}

	// 获取最后提交日期
	lastCmd := exec.Command("git", "-C", path, "log", "--format=%at", "--max-count=1")
	var lastOut bytes.Buffer
	lastCmd.Stdout = &lastOut
	if err := lastCmd.Run(); err != nil {
		PrintError("获取最后提交日期失败: %v", err)
		return time.Time{}, time.Time{}, 0, err
	}

	// 获取活跃天数
	daysCmd := exec.Command("git", "-C", path, "log", "--format=%ad", "--date=short", "--all")
	var daysOut bytes.Buffer
	daysCmd.Stdout = &daysOut
	if err := daysCmd.Run(); err != nil {
		PrintError("获取活跃天数失败: %v", err)
		return time.Time{}, time.Time{}, 0, err
	}

	// 解析首次提交时间
	firstUnix, err := strconv.ParseInt(strings.TrimSpace(firstOut.String()), 10, 64)
	if err != nil {
		PrintError("解析首次提交时间失败: %v", err)
		return time.Time{}, time.Time{}, 0, err
	}
	firstCommit := time.Unix(firstUnix, 0)

	// 解析最后提交时间
	lastUnix, err := strconv.ParseInt(strings.TrimSpace(lastOut.String()), 10, 64)
	if err != nil {
		PrintError("解析最后提交时间失败: %v", err)
		return time.Time{}, time.Time{}, 0, err
	}
	lastCommit := time.Unix(lastUnix, 0)

	// 计算活跃天数
	activeDays := make(map[string]bool)
	lines := strings.Split(daysOut.String(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			activeDays[line] = true
		}
	}

	return firstCommit, lastCommit, len(activeDays), nil
}

// 获取变更统计
func getChangeStats(path string) (int, int, int, error) {
	cmd := exec.Command("git", "-C", path, "log", "--numstat", "--pretty=tformat:")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		PrintError("获取变更统计失败: %v", err)
		return 0, 0, 0, err
	}

	var additions, deletions, fileChanges int
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			// 跳过二进制文件
			if parts[0] == "-" || parts[1] == "-" {
				fileChanges++
				continue
			}

			a, err := strconv.Atoi(parts[0])
			if err != nil {
				PrintWarning("解析添加行数失败: %v", err)
				continue
			}

			d, err := strconv.Atoi(parts[1])
			if err != nil {
				PrintWarning("解析删除行数失败: %v", err)
				continue
			}

			additions += a
			deletions += d
			fileChanges++
		}
	}

	return additions, deletions, fileChanges, nil
}

// 获取分支统计
func getBranchStats(path string) (map[string]bool, int, error) {
	cmd := exec.Command("git", "-C", path, "branch", "-a")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		PrintError("获取分支统计失败: %v", err)
		return nil, 0, err
	}

	branches := make(map[string]bool)
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 移除前面的 * 或空格
		if strings.HasPrefix(line, "*") {
			line = strings.TrimSpace(line[1:])
		}

		// 过滤远程分支引用
		if strings.HasPrefix(line, "remotes/") {
			parts := strings.Split(line, "/")
			if len(parts) > 2 {
				branchName := strings.Join(parts[2:], "/")
				branches[branchName] = true
			}
		} else {
			branches[line] = true
		}
	}

	return branches, len(branches), nil
}
