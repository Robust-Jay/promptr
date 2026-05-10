# promptr

AI 提示词管理器 — 命令行下的提示词搜索、管理和一键复制工具。不绑定任何 AI 平台，简洁实用。

## 功能

- **分类管理** — 扁平标签体系，一条提示词可属于多个分类
- **全文搜索** — 搜索 id、标题、分类、内容，快速定位
- **变量填充** — 支持 `{{变量}}` 占位符，复制前交互式填入
- **一键复制** — 填入变量后直接复制到系统剪贴板，粘贴即用
- **内建提示词库** — 二进制内嵌 13 条提示词，覆盖代码、写作、通用场景
- **用户自定义** — 可新增、修改、删除自己的提示词
- **跨平台** — Windows / macOS / Linux 均支持

## 安装

### 从源码编译

```bash
# 要求 Go 1.26+
git clone <repo-url>
cd toy
go build -o promptr ./cmd/promptr/

# 将 promptr 所在目录加入 PATH 即可全局使用
```

> Windows 下编译产物为 `promptr.exe`，Linux/macOS 为 `promptr`。

### 环境依赖

默认 Go 代理在该环境不可用，使用以下环境变量：

**PowerShell（临时）：**
```powershell
$env:GOPROXY="direct"; $env:GOINSECURE="*"; $env:GONOSUMCHECK="*"; $env:GONOSUMDB="*"
```

**持久化配置：**
```bash
go env -w GOPROXY=direct
go env -w GOINSECURE=*
go env -w GONOSUMCHECK=*
```

## 快速开始

### 查看内建提示词

```bash
# 列出所有提示词
promptr list

# 按分类筛选（显示所有 code 分类下的提示词）
promptr ls code

# 列出所有分类标签
promptr categories
```

### 搜索提示词

```bash
# 全文搜索
promptr search "单元测试"

# 输出示例：
#   Found 2 prompt(s) matching "单元测试":
#
#     unit-test     生成单元测试          [code, testing]
#     code-review   代码审查              [code, review]
```

### 查看提示词详情

```bash
promptr show unit-test
```

输出：
```
ID:       unit-test
Title:    生成单元测试
Category: code, testing
---
请为以下 {{language}} 代码编写完整的单元测试，使用 {{framework}} 框架。
要求覆盖：
1. 正常路径
2. 边界情况
3. 异常处理
4. 空值/null处理

代码：
{{code}}
```

### 复制提示词到剪贴板

```bash
promptr cp unit-test
```

交互式变量填入：
```
language: TypeScript
framework: Jest
code: function add(a, b) { return a + b; }
Copied to clipboard
```

> 粘贴到任何 AI 对话框中即得到完整的提示词文本。

### 添加自定义提示词

**交互式添加：**
```bash
promptr add
```
```
ID: my-review
Title: 我的代码审查
Category (comma-separated): code, review
Content (press Enter on empty line to finish):
> 请审查以下 {{language}} 代码的安全性：
> {{code}}
>
Created user/my-review.yaml
```

**参数式添加：**
```bash
promptr add \
  --id my-review \
  --title "我的代码审查" \
  --category "code,review" \
  --content "请审查以下 {{language}} 代码的安全性：\n{{code}}"
```

### 编辑提示词

```bash
# 编辑用户提示词（直接用 $EDITOR 打开）
promptr edit my-review

# 编辑内建提示词（自动拷贝到 user/ 目录再编辑，不修改内建版本）
promptr edit unit-test
```

### 删除提示词

```bash
# 只能删除用户提示词，内建提示词受保护
promptr rm my-review
```

## 命令参考

| 命令 | 别名 | 说明 |
|---|---|---|
| `search <query>` | `s` | 全文搜索提示词 |
| `list [category]` | `ls` | 列出提示词，可选按分类过滤 |
| `show <id>` | `cat` | 显示提示词完整内容 |
| `cp <id>` | `copy` | 复制到剪贴板（含变量填充） |
| `categories` | `tags` | 列出所有分类标签 |
| `add` | `new` | 新增提示词（交互式或参数式） |
| `edit <id>` | `e` | 编辑提示词（用 `$EDITOR`） |
| `rm <id>` | `delete` | 删除用户提示词 |
| `--version` | `-v` | 显示版本 |
| `--help` | `-h` | 显示帮助 |

## 内建提示词

以下 13 条提示词随二进制分发，首次运行自动解包到 `~/.promptr/builtin/`：

**代码类 (`code`)**
| ID | 标题 | 分类 |
|---|---|---|
| `unit-test` | 生成单元测试 | code, testing |
| `code-review` | 代码审查 | code, review |
| `refactor` | 重构代码 | code, refactoring |
| `explain-code` | 解释代码逻辑 | code, explanation |
| `debug` | 调试问题 | code, debugging |

**写作类 (`writing`)**
| ID | 标题 | 分类 |
|---|---|---|
| `blog-post` | 撰写博客文章 | writing, blog |
| `meeting-summary` | 会议纪要 | writing, meeting |
| `email-draft` | 撰写邮件 | writing, email |
| `documentation` | 生成文档 | writing, documentation |

**通用类 (`general`)**
| ID | 标题 | 分类 |
|---|---|---|
| `explain-concept` | 解释概念 | general, explanation |
| `pros-cons` | 权衡分析 | general, analysis |
| `summarize` | 总结内容 | general, summary |
| `brainstorming` | 头脑风暴 | general, ideation |

## 自定义提示词格式

提示词以 YAML 文件存储，格式如下：

```yaml
id: my-prompt
title: 我的提示词
category: [code, review]
content: |
  请分析以下 {{language}} 代码的性能问题：
  {{code}}
```

文件存放于 `~/.promptr/user/` 目录。也可以手动创建 `.yaml` 文件放入该目录。

### 变量说明

提示词 `content` 中可使用 `{{变量名}}` 定义占位符。执行 `promptr cp` 时系统会逐一提示填入变量值，然后复制组装后的完整文本。

- 变量名只能包含字母、数字和下划线（如 `{{language}}`、`{{code}}`）
- 相同变量只提示一次
- 空值会被替换为空字符串

## 目录结构

```
~/.promptr/
├── builtin/          ← 首次运行从二进制解包的内建提示词
│   ├── code/         ← 5 条代码类提示词
│   ├── writing/      ← 4 条写作类提示词
│   └── general/      ← 4 条通用类提示词
└── user/             ← 用户自建提示词（优先级高于 builtin）
    └── my-review.yaml
```

用户提示词优先级高于内建提示词：同 ID 的提示词，`user/` 会覆盖 `builtin/`。

## 运行测试

```bash
go test ./...

# 单个包
go test ./internal/prompt/ -v
go test ./internal/clipboard/ -v
```

## 项目结构

```
cmd/promptr/main.go          CLI 入口，命令分发
internal/
  prompt/prompt.go           Prompt 数据模型，YAML 读写
  prompt/store.go            存储层：搜索、列表、CRUD
  prompt/vars.go             变量检测与交互式填充
  builtin/builtin.go         embed.FS，首次运行解包
  builtin/{code,writing,general}/  内建提示词 YAML 文件
  clipboard/clipboard*.go    平台剪贴板实现
  clipboard/clipboard_test.go
  prompt/prompt_test.go
pkg/                         公开 API（预留）
```

## 许可证

MIT
