# 代码统计工具 (Code Stats)

这是一个强大的代码统计分析工具，可以帮助开发者分析代码仓库的各种统计指标，包括代码行数、注释比例、文件大小分布等。

## 特性

- 多线程并行处理，快速分析大型代码库
- 支持排除指定目录和文件类型
- 按语言、文件扩展名自动分类统计
- 提供详细的统计报告，包括代码行数、注释比例等
- 内置可视化图表，直观展示代码组成和语言分布
- 彩色终端输出，直观展示分析进度

## 安装

```bash
go install github.com/lllllan02/code-stats@latest
```

或者从源码编译:

```bash
git clone https://github.com/lllllan02/code-stats.git
cd code-stats
go build
```

## 使用方法

最简单的用法是分析当前目录:

```bash
code-stats
```

### 常用选项

```
-path          指定分析的目录路径
-exclude-dirs  排除的目录，逗号分隔
-exclude-exts  排除的文件扩展名，逗号分隔
-max-workers   最大并发工作数量
-follow-links  是否跟踪符号链接
-verbose       显示详细日志
-show-files    在报告中显示文件详情
-top           在报告中显示前N个文件
-output        报告输出文件路径（默认为code-stats-report.html）
-help          显示帮助信息
```

### 示例

分析指定目录，排除node_modules和.git:

```bash
code-stats -path=/path/to/project -exclude-dirs=node_modules,.git,vendor
```

指定报告的输出文件名:

```bash
code-stats -output=我的项目报告.html -show-files
```

只显示前20个最大的文件:

```bash
code-stats -show-files -top=20
```

## 统计指标与可视化

工具会生成包含以下内容的详细报告:

- **图表可视化**
  - 代码组成比例饼图（代码行、注释行、空白行）
  - 主要编程语言分布柱状图

- **总体统计**
  - 文件总数、总代码量
  - 代码行数及占比
  - 注释行数及占比
  - 空白行数及占比
  - 注释比例(注释行/代码行)
  - 平均文件大小、平均行长度

- **语言统计**
  - 各编程语言的文件数、代码行数、注释行数
  - 各语言的注释比例和平均行长度

- **扩展名统计**
  - 各文件扩展名的文件数、代码行数、总行数
  - 各扩展名的平均文件大小

- **文件统计**
  - 按大小排序的文件列表
  - 按代码行数排序的文件列表

## 许可证

MIT 