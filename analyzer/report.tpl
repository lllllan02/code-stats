<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>代码统计报告</title>
    <link rel="stylesheet" type="text/css" href="https://cdn.datatables.net/1.11.5/css/jquery.dataTables.min.css">
    <script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
    <script src="https://cdn.datatables.net/1.11.5/js/jquery.dataTables.min.js"></script>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f9f9f9;
        }
        h1, h2, h3 {
            color: #2c3e50;
        }
        h1 {
            border-bottom: 2px solid #3498db;
            padding-bottom: 10px;
        }
        h2 {
            margin-top: 30px;
            border-bottom: 1px solid #bdc3c7;
            padding-bottom: 5px;
        }
        .summary {
            background-color: #fff;
            border-radius: 5px;
            padding: 20px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
            margin-bottom: 20px;
        }
        .summary-item {
            margin-bottom: 10px;
            display: flex;
        }
        .summary-label {
            font-weight: bold;
            width: 200px;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
            background-color: #fff;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
            border-radius: 5px;
            overflow: hidden;
        }
        th {
            background-color: #3498db;
            color: white;
            padding: 12px 15px;
            text-align: left;
        }
        td {
            padding: 10px 15px;
            border-bottom: 1px solid #ddd;
        }
        tr {
            background-color: #fff;
        }
        tr:nth-child(even) {
            background-color: #f8f9fa;
        }
        tr:hover {
            background-color: #e9f7fe;
        }
        tr:last-child td {
            border-bottom: none;
        }
        /* 覆盖DataTables默认样式 */
        .dataTables_wrapper .dataTables_length, 
        .dataTables_wrapper .dataTables_filter, 
        .dataTables_wrapper .dataTables_info, 
        .dataTables_wrapper .dataTables_processing, 
        .dataTables_wrapper .dataTables_paginate {
            color: #333;
        }
        .dataTables_wrapper .dataTables_paginate .paginate_button.current, 
        .dataTables_wrapper .dataTables_paginate .paginate_button.current:hover {
            background: #3498db;
            color: white !important;
            border: 1px solid #3498db;
        }
        table.dataTable.stripe tbody tr.odd, 
        table.dataTable.display tbody tr.odd {
            background-color: #fff;
        }
        table.dataTable.stripe tbody tr.even, 
        table.dataTable.display tbody tr.even {
            background-color: #f8f9fa;
        }
        table.dataTable.hover tbody tr:hover, 
        table.dataTable.display tbody tr:hover {
            background-color: #e9f7fe;
        }
        .chart-container {
            display: flex;
            justify-content: space-between;
            margin: 30px 0;
        }
        .chart {
            width: 48%;
            background-color: #fff;
            border-radius: 5px;
            padding: 20px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        .timestamp {
            color: #7f8c8d;
            font-style: italic;
            margin-bottom: 30px;
        }
        .path-info {
            color: #7f8c8d;
            font-style: italic;
            margin-bottom: 30px;
        }
        footer {
            margin-top: 50px;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #7f8c8d;
            font-size: 0.9em;
            text-align: center;
        }
        .file-browser-container {
            display: flex;
            flex-direction: row;
            background-color: #fff;
            border-radius: 5px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
            margin-bottom: 30px;
            min-height: 500px;
            max-height: 700px;
            overflow: hidden;
        }
        
        .directory-tree {
            width: 30%;
            min-width: 250px;
            padding: 15px;
            border-right: 1px solid #eee;
            overflow-y: auto;
            overflow-x: hidden;
            height: 100%;
            max-height: 700px;
            scrollbar-width: thin; /* Firefox */
            scrollbar-color: #ccc #f8f9fa; /* Firefox */
        }
        
        /* WebKit 滚动条样式 */
        .directory-tree::-webkit-scrollbar {
            width: 6px;
        }
        
        .directory-tree::-webkit-scrollbar-track {
            background: #f8f9fa;
        }
        
        .directory-tree::-webkit-scrollbar-thumb {
            background-color: #ccc;
            border-radius: 20px;
        }
        
        .file-details {
            width: 70%;
            padding: 20px;
            overflow: auto;
            height: 100%;
        }
        
        /* 目录树样式 */
        .treeview {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            position: relative;
            width: 100%;
        }
        
        .treeview ul {
            list-style: none;
            padding-left: 20px;
            margin: 0;
        }
        
        .treeview li {
            margin: 5px 0;
            position: relative;
        }
        
        .directory-item {
            position: relative;
            width: 100%;
        }
        
        .treeview .directory {
            cursor: pointer;
            font-weight: bold;
            color: #2c3e50;
            display: inline-block;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
            max-width: 90%;
        }
        
        .treeview .file {
            cursor: pointer;
            color: #3498db;
            display: inline-block;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
            max-width: 90%;
        }
        
        .treeview .file:hover {
            text-decoration: underline;
        }
        
        /* 选中文件的样式 */
        .treeview .file.selected {
            font-weight: bold;
            color: #2980b9;
            background-color: rgba(52, 152, 219, 0.1);
            border-radius: 3px;
            padding: 2px 5px;
            margin: -2px 0;
        }
        
        .treeview .collapsed > ul {
            display: none;
        }
        
        .treeview .expanded > ul {
            display: block;
        }
        
        .treeview .directory:before {
            content: "📁 ";
            margin-right: 3px;
        }
        
        .treeview .expanded > .directory:before {
            content: "📂 ";
            margin-right: 3px;
        }
        
        .treeview .file:before {
            content: "📄 ";
            margin-right: 3px;
        }
        
        /* 添加展开/折叠指示器 */
        .directory-item.collapsed:after {
            content: "+";
            position: absolute;
            right: 10px;
            top: 2px;
            font-size: 14px;
            color: #999;
            font-weight: bold;
        }
        
        .directory-item.expanded:after {
            content: "-";
            position: absolute;
            right: 10px;
            top: 2px;
            font-size: 14px;
            color: #999;
            font-weight: bold;
        }
        
        /* 文件详情样式 */
        .file-details h3 {
            border-bottom: 1px solid #eee;
            padding-bottom: 10px;
            margin-top: 0;
        }
        
        .file-details .info-group {
            margin-bottom: 15px;
        }
        
        .file-details .info-label {
            font-weight: bold;
            margin-right: 10px;
        }
        
        .file-details .metrics {
            display: flex;
            flex-wrap: wrap;
            margin-top: 20px;
        }
        
        .file-details .metric-box {
            background-color: #f8f9fa;
            border-radius: 4px;
            padding: 10px 15px;
            margin: 0 10px 10px 0;
            width: calc(33.33% - 10px);
            box-sizing: border-box;
        }
        
        .file-details .metric-value {
            font-size: 18px;
            font-weight: bold;
            color: #3498db;
        }
        
        .file-details .metric-name {
            font-size: 12px;
            color: #7f8c8d;
        }
        
        .file-mini-chart {
            margin-top: 20px;
            height: 200px;
        }
        
        .no-file-selected {
            color: #7f8c8d;
            font-style: italic;
            text-align: center;
            margin-top: 40px;
        }

        /* 导航栏样式 */
        .nav-container {
            position: sticky;
            top: 0;
            z-index: 100;
            background-color: #fff;
            border-radius: 5px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
            margin-bottom: 20px;
            display: flex;
            overflow-x: auto;
            white-space: nowrap;
            padding: 5px;
        }
        
        .nav-item {
            display: inline-block;
            padding: 10px 15px;
            cursor: pointer;
            border-radius: 5px;
            margin: 0 5px;
            font-weight: bold;
            color: #2c3e50;
            transition: all 0.3s ease;
        }
        
        .nav-item.active {
            background-color: #3498db;
            color: white;
        }
        
        .nav-item:hover:not(.active) {
            background-color: #e9f7fe;
        }
        
        .section {
            display: none;
        }
        
        .section.active {
            display: block;
        }
    </style>
</head>
<body>
    <div style="position: relative; margin-bottom: 0;">
        <h1 style="margin-bottom: 0; padding-bottom: 10px;">代码统计报告</h1>
        <div style="position: absolute; right: 0; bottom: 15px; text-align: right; font-size: 0.9em;">
            <span class="timestamp">生成时间: {{.GenerationTime}}</span> | 
            <span class="path-info">分析路径: {{.Stats.Path}}</span>
        </div>
    </div>

    <!-- 导航栏 -->
    <div class="nav-container" style="margin-top: 0; border-top: 2px solid #3498db;">
        <div class="nav-item active" data-target="section-summary">总体摘要</div>
        <div class="nav-item" data-target="section-languages">语言统计</div>
        <div class="nav-item" data-target="section-extensions">扩展名统计</div>
        <div class="nav-item" data-target="section-files-size">最大文件</div>
        <div class="nav-item" data-target="section-files-lines">最长文件</div>
        {{if .HasGitStats}}
        <div class="nav-item" data-target="section-git-stats">Git 统计</div>
        {{end}}
        <div class="nav-item" data-target="section-file-browser">文件浏览器</div>
    </div>

    <!-- 总体摘要区域 -->
    <div id="section-summary" class="section active">
        <div class="summary">
            <div class="summary-item"><span class="summary-label">总文件数:</span> {{.Stats.TotalFiles}} 个文件</div>
            <div class="summary-item"><span class="summary-label">总代码量:</span> {{.Stats.TotalLines}} 行 ({{printf "%.2f" (divideBy .Stats.TotalSize 1048576)}} MB)</div>
            <div class="summary-item"><span class="summary-label">代码行数:</span> {{.Stats.CodeLines}} 行 ({{printf "%.1f%%" (multiply .Stats.CodeDensity 100)}})</div>
            <div class="summary-item"><span class="summary-label">注释行数:</span> {{.Stats.CommentLines}} 行 ({{printf "%.1f%%" (multiply .Stats.CommentDensity 100)}})</div>
            <div class="summary-item"><span class="summary-label">空白行数:</span> {{.Stats.BlankLines}} 行 ({{printf "%.1f%%" (multiply .Stats.AvgBlankLines 100)}})</div>
            <div class="summary-item"><span class="summary-label">注释比例:</span> {{printf "%.2f" .Stats.CommentRatio}} (注释行/代码行)</div>
            <div class="summary-item"><span class="summary-label">平均文件大小:</span> {{printf "%.2f" (divideBy .Stats.AvgFileSize 1024)}} KB</div>
            <div class="summary-item"><span class="summary-label">平均行长度:</span> {{printf "%.1f" .Stats.AvgLineLength}} 字符/行</div>
        </div>

        <!-- 可视化图表 -->
        <div class="chart-container">
            <div class="chart">
                <h3>代码组成</h3>
                <div style="position: relative; height: 200px;">
                    <canvas id="compositionChart"></canvas>
                </div>
            </div>
            <div class="chart">
                <h3>语言分布</h3>
                <div style="position: relative; height: 200px;">
                    <canvas id="languageChart"></canvas>
                </div>
            </div>
        </div>
    </div>

    <!-- Git 统计区域 -->
    {{if .HasGitStats}}
    <div id="section-git-stats" class="section">
        <div class="summary">
            <h3>Git 仓库基本信息</h3>
            <div class="summary-item"><span class="summary-label">提交总数:</span> {{.Stats.GitStats.CommitCount}} 次提交</div>
            <div class="summary-item"><span class="summary-label">贡献者数量:</span> {{.Stats.GitStats.ContributorCount}} 人</div>
            <div class="summary-item"><span class="summary-label">首次提交时间:</span> {{formatTime .Stats.GitStats.FirstCommitDate}}</div>
            <div class="summary-item"><span class="summary-label">最后提交时间:</span> {{formatTime .Stats.GitStats.LastCommitDate}}</div>
            <div class="summary-item"><span class="summary-label">活跃天数:</span> {{.Stats.GitStats.ActiveDays}} 天</div>
            <div class="summary-item"><span class="summary-label">添加的行数总计:</span> {{.Stats.GitStats.TotalAdditions}} 行</div>
            <div class="summary-item"><span class="summary-label">删除的行数总计:</span> {{.Stats.GitStats.TotalDeletions}} 行</div>
            <div class="summary-item"><span class="summary-label">文件变更总数:</span> {{.Stats.GitStats.TotalFileChanges}} 个</div>
            <div class="summary-item"><span class="summary-label">分支数量:</span> {{.Stats.GitStats.BranchCount}} 个</div>
        </div>

        {{if .TopContributors}}
        <h3>贡献者排行 (前 {{.ContributorsLimit}} 名)</h3>
        <table id="contributors-table" class="display">
            <thead>
                <tr>
                    <th>排名</th>
                    <th>贡献者</th>
                    <th>提交数</th>
                    <th>占比</th>
                </tr>
            </thead>
            <tbody>
                {{range $index, $contributor := .TopContributors}}
                <tr>
                    <td>{{add $index 1}}</td>
                    <td>{{$contributor.Name}}</td>
                    <td>{{$contributor.CommitCount}}</td>
                    <td>{{printf "%.1f%%" (multiply (divideBy $contributor.CommitCount $.Stats.GitStats.CommitCount) 100)}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
        {{end}}

        <!-- Git 统计可视化 -->
        <div class="chart-container">
            <div class="chart">
                <h3>贡献者分布</h3>
                <div style="position: relative; height: 200px;">
                    <canvas id="contributorsChart"></canvas>
                </div>
            </div>
            <div class="chart">
                <h3>提交活跃度</h3>
                <div style="position: relative; height: 200px;">
                    <canvas id="commitActivityChart"></canvas>
                </div>
            </div>
        </div>
    </div>
    {{end}}

    <!-- 语言统计区域 -->
    <div id="section-languages" class="section">
        {{if gt (len .TopLanguages) 0}}
        <table id="language-table" class="display">
            <thead>
                <tr>
                    <th>语言</th>
                    <th>文件数</th>
                    <th>代码行</th>
                    <th>注释行</th>
                    <th>空白行</th>
                    <th>注释比例</th>
                    <th>平均行长度</th>
                </tr>
            </thead>
            <tbody>
                {{range .TopLanguages}}
                {{if ne .Name ""}}
                <tr>
                    <td>{{.Name}}</td>
                    <td>{{.Stats.TotalFiles}}</td>
                    <td>{{.Stats.CodeLines}}</td>
                    <td>{{.Stats.CommentLines}}</td>
                    <td>{{.Stats.BlankLines}}</td>
                    <td>{{printf "%.2f" .Stats.CommentRatio}}</td>
                    <td>{{printf "%.1f" .Stats.AvgLineLength}}</td>
                </tr>
                {{end}}
                {{end}}
            </tbody>
        </table>
        {{end}}
    </div>

    <!-- 扩展名统计区域 -->
    <div id="section-extensions" class="section">
        {{if gt (len .SortedExts) 0}}
        <table id="extension-table" class="display">
            <thead>
                <tr>
                    <th>扩展名</th>
                    <th>文件数</th>
                    <th>代码行</th>
                    <th>总行数</th>
                    <th>平均大小(KB)</th>
                </tr>
            </thead>
            <tbody>
                {{range .SortedExts}}
                <tr>
                    <td>{{.Name}}</td>
                    <td>{{.Stats.TotalFiles}}</td>
                    <td>{{.Stats.CodeLines}}</td>
                    <td>{{.Stats.TotalLines}}</td>
                    <td>{{printf "%.2f" (divideBy .Stats.AvgFileSize 1024)}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
        {{end}}
    </div>

    <!-- 按大小排序的文件区域 -->
    <div id="section-files-size" class="section">
        {{if gt (len .FilesBySize) 0}}
        <table id="files-by-size-table" class="display">
            <thead>
                <tr>
                    <th>文件路径</th>
                    <th>大小(KB)</th>
                    <th>总行数</th>
                    <th>代码行</th>
                    <th>注释行</th>
                </tr>
            </thead>
            <tbody>
                {{range .FilesBySize}}
                <tr>
                    <td>{{.Path}}</td>
                    <td>{{printf "%.2f" (divideBy .TotalSize 1024)}}</td>
                    <td>{{.TotalLines}}</td>
                    <td>{{.CodeLines}}</td>
                    <td>{{.CommentLines}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
        {{end}}
    </div>

    <!-- 按代码行数排序的文件区域 -->
    <div id="section-files-lines" class="section">
        {{if gt (len .FilesByLines) 0}}
        <table id="files-by-lines-table" class="display">
            <thead>
                <tr>
                    <th>文件路径</th>
                    <th>代码行</th>
                    <th>注释行</th>
                    <th>空白行</th>
                    <th>注释比例</th>
                </tr>
            </thead>
            <tbody>
                {{range .FilesByLines}}
                <tr>
                    <td>{{.Path}}</td>
                    <td>{{.CodeLines}}</td>
                    <td>{{.CommentLines}}</td>
                    <td>{{.BlankLines}}</td>
                    <td>{{printf "%.2f" (commentRatio .CommentLines .CodeLines)}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
        {{end}}
    </div>

    <!-- 文件浏览器区域 -->
    <div id="section-file-browser" class="section">
        <div class="summary">
            <h3>文件浏览器</h3>
            <p>点击目录树中的文件可查看详细信息</p>
        </div>
        
        <div class="file-browser-container">
            <div class="directory-tree">
                <h4 style="margin-top: 0; color: #3498db;">文件结构</h4>
                <div class="treeview" id="fileTree"></div>
            </div>
            <div class="file-details" id="fileDetails">
                <div class="no-file-selected">
                    <p>请从左侧目录树中选择一个文件查看详情</p>
                </div>
            </div>
        </div>
    </div>
    
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <script>
        // 初始化所有DataTable
        $(document).ready(function() {
            // 为所有table.display表格启用DataTables排序功能
            // 使用 $.fn.dataTable.isDataTable 检查表格是否已经初始化
            $('table.display:not(#files-dashboard)').each(function() {
                if (!$.fn.dataTable.isDataTable(this)) {
                    $(this).DataTable({
                        paging: false,
                        searching: false,
                        info: false,
                        order: [],  // 保持默认排序
                        stripeClasses: [],  // 禁用交替行类
                        rowCallback: function(row, data) {
                            // 为表格行添加统一样式
                            $(row).addClass('unified-row');
                        }
                    });
                }
            });
            
            // 导航栏切换功能
            $('.nav-item').click(function() {
                // 移除所有导航项的active类
                $('.nav-item').removeClass('active');
                // 为当前点击的导航项添加active类
                $(this).addClass('active');
                
                // 获取目标部分的ID
                const targetId = $(this).data('target');
                
                // 隐藏所有部分
                $('.section').removeClass('active');
                // 显示目标部分
                $('#' + targetId).addClass('active');
                
                // 如果切换到Git统计页面，重新初始化图表
                if (targetId === 'section-git-stats') {
                    setTimeout(function() {
                        initGitCharts();
                    }, 100);
                }
            });
        });
        
        // 代码组成图表
        const compositionCtx = document.getElementById('compositionChart').getContext('2d');
        new Chart(compositionCtx, {
            type: 'pie',
            data: {
                labels: ['代码行', '注释行', '空白行'],
                datasets: [{
                    data: [{{.Stats.CodeLines}}, {{.Stats.CommentLines}}, {{.Stats.BlankLines}}],
                    backgroundColor: [
                        'rgba(54, 162, 235, 0.7)',
                        'rgba(255, 205, 86, 0.7)',
                        'rgba(201, 203, 207, 0.7)'
                    ],
                    borderColor: [
                        'rgb(54, 162, 235)',
                        'rgb(255, 205, 86)',
                        'rgb(201, 203, 207)'
                    ],
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'bottom',
                    }
                }
            }
        });

        // 语言分布图表
        const languageCtx = document.getElementById('languageChart').getContext('2d');
        new Chart(languageCtx, {
            type: 'bar',
            data: {
                labels: [{{range $i, $lang := .TopLanguages}}{{if $i}}, {{end}}'{{$lang.Name}}'{{end}}],
                datasets: [{
                    label: '代码行数',
                    data: [{{range $i, $lang := .TopLanguages}}{{if $i}}, {{end}}{{$lang.Stats.CodeLines}}{{end}}],
                    backgroundColor: 'rgba(54, 162, 235, 0.7)',
                    borderColor: 'rgb(54, 162, 235)',
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    y: {
                        beginAtZero: true
                    }
                },
                plugins: {
                    legend: {
                        display: false
                    }
                }
            }
        });
        
        // 文件数据
        const fileData = {
            {{range $i, $file := .Stats.FileStats}}
            "{{$file.Path}}": {
                path: "{{$file.Path}}",
                language: "{{if $file.Language}}{{$file.Language}}{{else}}未识别{{end}}",
                extension: "{{if ext $file.Path}}{{ext $file.Path}}{{else}}(无扩展名){{end}}",
                size: {{printf "%.2f" (divideBy $file.TotalSize 1024)}},
                totalLines: {{$file.TotalLines}},
                codeLines: {{$file.CodeLines}},
                commentLines: {{$file.CommentLines}},
                blankLines: {{$file.BlankLines}},
                commentRatio: {{printf "%.2f" (commentRatio $file.CommentLines $file.CodeLines)}},
                avgLineLength: {{printf "%.1f" $file.AvgLineLength}}
            }{{if lt $i (subtract (len $.Stats.FileStats) 1)}},{{end}}
            {{end}}
        };
        
        // 构建目录树结构
        function buildDirectoryTree() {
            const root = { name: "根目录", isDirectory: true, children: {} };
            
            // 处理每个文件路径
            Object.keys(fileData).forEach(path => {
                // 统一路径分隔符，处理 Windows 和 UNIX 路径
                const normalizedPath = path.replace(/\\/g, '/');
                const parts = normalizedPath.split('/').filter(p => p.length > 0);
                
                let current = root;
                
                // 逐级创建目录结构
                for (let i = 0; i < parts.length; i++) {
                    const part = parts[i];
                    
                    // 如果是最后一部分，则为文件
                    if (i === parts.length - 1) {
                        if (!current.children[part]) {
                            current.children[part] = { 
                                name: part, 
                                isDirectory: false, 
                                path: path 
                            };
                        }
                    } else {
                        // 否则是目录
                        if (!current.children[part]) {
                            current.children[part] = { 
                                name: part, 
                                isDirectory: true, 
                                children: {} 
                            };
                        }
                        current = current.children[part];
                    }
                }
            });
            
            return root;
        }
        
        // 将目录树渲染为HTML
        function renderDirectoryTree(node) {
            if (!node.isDirectory) return "";
            
            let html = '<ul>';
            
            // 获取目录和文件并排序
            const items = Object.values(node.children);
            const directories = items.filter(item => item.isDirectory);
            const files = items.filter(item => !item.isDirectory);
            
            // 按名称排序
            directories.sort((a, b) => a.name.localeCompare(b.name));
            files.sort((a, b) => a.name.localeCompare(b.name));
            
            // 先渲染目录
            directories.forEach(dir => {
                // 添加 title 属性以显示完整名称
                html += '<li class="directory-item collapsed">' +
                        '<span class="directory" title="' + dir.name + '">' + dir.name + '</span>' +
                        renderDirectoryTree(dir) +
                        '</li>';
            });
            
            // 再渲染文件
            files.forEach(file => {
                // 同样添加 title 属性
                html += '<li><span class="file" data-path="' + file.path + '" title="' + file.name + '">' + 
                        file.name + '</span></li>';
            });
            
            html += '</ul>';
            return html;
        }
        
        // 显示文件详情
        function showFileDetails(path) {
            const file = fileData[path];
            if (!file) return;
            
            let html = '<h3>' + path + '</h3>';
            
            html += '<div class="info-group">' +
                    '<span class="info-label">语言:</span>' + file.language + 
                    '<span class="info-label" style="margin-left:20px;">扩展名:</span>' + file.extension + 
                    '</div>';
            
            html += '<div class="metrics">';
            
            // 文件大小指标
            html += '<div class="metric-box">' +
                    '<div class="metric-value">' + file.size + ' KB</div>' +
                    '<div class="metric-name">文件大小</div>' +
                    '</div>';
            
            // 总行数指标
            html += '<div class="metric-box">' +
                    '<div class="metric-value">' + file.totalLines + '</div>' +
                    '<div class="metric-name">总行数</div>' +
                    '</div>';
            
            // 代码行指标
            html += '<div class="metric-box">' +
                    '<div class="metric-value">' + file.codeLines + '</div>' +
                    '<div class="metric-name">代码行</div>' +
                    '</div>';
            
            // 注释行指标
            html += '<div class="metric-box">' +
                    '<div class="metric-value">' + file.commentLines + '</div>' +
                    '<div class="metric-name">注释行</div>' +
                    '</div>';
            
            // 空白行指标
            html += '<div class="metric-box">' +
                    '<div class="metric-value">' + file.blankLines + '</div>' +
                    '<div class="metric-name">空白行</div>' +
                    '</div>';
            
            // 注释比例指标
            html += '<div class="metric-box">' +
                    '<div class="metric-value">' + file.commentRatio + '</div>' +
                    '<div class="metric-name">注释比例</div>' +
                    '</div>';
            
            // 平均行长指标
            html += '<div class="metric-box">' +
                    '<div class="metric-value">' + file.avgLineLength + '</div>' +
                    '<div class="metric-name">平均行长度(字符)</div>' +
                    '</div>';
            
            html += '</div>';
            
            // 添加文件组成饼图
            html += '<div class="file-mini-chart">' +
                    '<canvas id="fileCompositionChart"></canvas>' +
                    '</div>';
            
            document.getElementById('fileDetails').innerHTML = html;
            
            // 绘制饼图前先检查并销毁已有图表
            setTimeout(function() {
                const chartEl = document.getElementById('fileCompositionChart');
                if (chartEl) {
                    // 销毁已存在的图表实例
                    const existingChart = Chart.getChart(chartEl);
                    if (existingChart) {
                        existingChart.destroy();
                    }
                    
                    // 绘制新饼图
                    const ctx = chartEl.getContext('2d');
                    new Chart(ctx, {
                        type: 'pie',
                        data: {
                            labels: ['代码行', '注释行', '空白行'],
                            datasets: [{
                                data: [file.codeLines, file.commentLines, file.blankLines],
                                backgroundColor: [
                                    'rgba(54, 162, 235, 0.7)',
                                    'rgba(255, 205, 86, 0.7)',
                                    'rgba(201, 203, 207, 0.7)'
                                ],
                                borderColor: [
                                    'rgb(54, 162, 235)',
                                    'rgb(255, 205, 86)',
                                    'rgb(201, 203, 207)'
                                ],
                                borderWidth: 1
                            }]
                        },
                        options: {
                            responsive: true,
                            maintainAspectRatio: false,
                            plugins: {
                                legend: {
                                    position: 'bottom',
                                },
                                title: {
                                    display: true,
                                    text: '文件组成'
                                }
                            }
                        }
                    });
                }
            }, 0);
        }
        
        // 当文档加载完成时初始化目录树
        $(document).ready(function() {
            // 初始化文件树
            const tree = buildDirectoryTree();
            document.getElementById('fileTree').innerHTML = renderDirectoryTree(tree);
            
            // 为目录添加点击事件 - 折叠/展开
            $(document).on('click', '.directory', function(e) {
                e.stopPropagation(); // 防止事件冒泡
                const li = $(this).parent();
                li.toggleClass('collapsed expanded');
                
                // 当展开目录时，确保目录树可以滚动
                if (li.hasClass('expanded')) {
                    setTimeout(function() {
                        // 滚动到可视区域
                        const container = $('.directory-tree');
                        const position = li.position().top;
                        if (position < 0 || position > container.height()) {
                            container.animate({
                                scrollTop: container.scrollTop() + position - 50
                            }, 200);
                        }
                    }, 100);
                }
            });
            
            // 为文件添加点击事件 - 显示详情
            $(document).on('click', '.file', function(e) {
                e.stopPropagation(); // 防止事件冒泡
                const path = $(this).data('path');
                showFileDetails(path);
                
                // 高亮当前选中的文件
                $('.file').removeClass('selected');
                $(this).addClass('selected');
            });
            
            // 默认展开根目录
            $('#fileTree > ul > li').addClass('expanded').removeClass('collapsed');
            
            // 初始化文件详情区域
            $('.file-details').html('<div class="no-file-selected"><p>请从左侧目录树中选择一个文件查看详情</p></div>');
            
            // 为贡献者表格初始化 DataTable - 避免重复初始化
            if (document.getElementById('contributors-table') && !$.fn.dataTable.isDataTable('#contributors-table')) {
                $('#contributors-table').DataTable({
                    paging: false,
                    searching: false,
                    info: false,
                    order: [],
                    stripeClasses: []
                });
            }
            
            // 初始化Git统计图表
            initGitCharts();
            
            // 页面加载完成后，再次尝试初始化文件浏览器
            // 这是为了解决某些浏览器中文件树不显示的问题
            setTimeout(function() {
                if ($('#fileTree').children().length === 0) {
                    console.log('重试初始化文件树...');
                    document.getElementById('fileTree').innerHTML = renderDirectoryTree(buildDirectoryTree());
                    $('#fileTree > ul > li').addClass('expanded').removeClass('collapsed');
                }
                
                // 确保文件浏览器区域在切换时正确显示
                $('.nav-item[data-target="section-file-browser"]').on('click', function() {
                    setTimeout(function() {
                        if ($('#fileTree').children().length === 0) {
                            document.getElementById('fileTree').innerHTML = renderDirectoryTree(buildDirectoryTree());
                            $('#fileTree > ul > li').addClass('expanded').removeClass('collapsed');
                        }
                    }, 100);
                });
            }, 500);
        });
        
        // 初始化Git统计图表
        function initGitCharts() {
            {{if .HasGitStats}}
            // 检查图表元素是否存在
            const contributorsChartEl = document.getElementById('contributorsChart');
            const commitActivityChartEl = document.getElementById('commitActivityChart');
            
            // 销毁已存在的图表实例，避免重复初始化
            if (contributorsChartEl) {
                const existingChart = Chart.getChart(contributorsChartEl);
                if (existingChart) {
                    existingChart.destroy();
                }
                
                // 贡献者分布图
                const contributorsCtx = contributorsChartEl.getContext('2d');
                new Chart(contributorsCtx, {
                    type: 'pie',
                    data: {
                        labels: [{{range $i, $contributor := .TopContributors}}{{if $i}}, {{end}}'{{$contributor.Name}}'{{end}}],
                        datasets: [{
                            data: [{{range $i, $contributor := .TopContributors}}{{if $i}}, {{end}}{{$contributor.CommitCount}}{{end}}],
                            backgroundColor: [
                                'rgba(54, 162, 235, 0.7)',
                                'rgba(255, 99, 132, 0.7)',
                                'rgba(255, 205, 86, 0.7)',
                                'rgba(75, 192, 192, 0.7)',
                                'rgba(153, 102, 255, 0.7)',
                                'rgba(255, 159, 64, 0.7)',
                                'rgba(199, 199, 199, 0.7)'
                            ],
                            borderWidth: 1
                        }]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        plugins: {
                            legend: {
                                position: 'right',
                                labels: {
                                    boxWidth: 12
                                }
                            }
                        }
                    }
                });
            }
            
            if (commitActivityChartEl) {
                const existingChart = Chart.getChart(commitActivityChartEl);
                if (existingChart) {
                    existingChart.destroy();
                }
                
                // 添加/删除行统计图
                const activityCtx = commitActivityChartEl.getContext('2d');
                new Chart(activityCtx, {
                    type: 'bar',
                    data: {
                        labels: ['代码变更'],
                        datasets: [
                            {
                                label: '添加行数',
                                data: [{{.Stats.GitStats.TotalAdditions}}],
                                backgroundColor: 'rgba(75, 192, 192, 0.7)',
                                borderColor: 'rgb(75, 192, 192)',
                                borderWidth: 1
                            },
                            {
                                label: '删除行数',
                                data: [{{.Stats.GitStats.TotalDeletions}}],
                                backgroundColor: 'rgba(255, 99, 132, 0.7)',
                                borderColor: 'rgb(255, 99, 132)',
                                borderWidth: 1
                            }
                        ]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        plugins: {
                            legend: {
                                position: 'top'
                            },
                            tooltip: {
                                callbacks: {
                                    label: function(context) {
                                        return context.dataset.label + ': ' + context.raw + ' 行';
                                    }
                                }
                            }
                        },
                        scales: {
                            y: {
                                beginAtZero: true
                            }
                        }
                    }
                });
            }
            {{end}}
        }
    </script>
</body>
</html>
