<p align="center">
  <img src="assets/logo.png" width="180" height="180" alt="git-air logo" />
</p>

<h1 align="center">git-air 🍃</h1>

<p align="center">
  <strong>纯 Go 编写、像空气一样轻盈极速的 Git 原生 AI 代码评审工具。</strong><br>
  <em>A lightweight, lightning-fast, Git-native AI code reviewer written in pure Go.</em>
</p>

<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?style=flat&logo=go" alt="Go Version" /></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License" /></a>
  <a href="https://github.com/unclesam-ly/git-air/pulls"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg" alt="PRs Welcome" /></a>
</p>

<p align="center">
  <a href="README.md"><strong>简体中文</strong></a> |
  <a href="README_EN.md">English</a> |
  <a href="README_JA.md">日本語</a>
</p>

---

## 🌟 为什么做 git-air？

### 痛点：AI 写代码一时爽，Review 起来火葬场

现在的开发工作流几乎离不开 AI 辅助（Cursor、Copilot、各种大模型插件）。几秒钟生成上百行代码，写起来确实爽快。

**但这也带来了一个全新的致命痛点：**
- **“看起来太完美”的视觉欺骗**：AI 生成的代码排版工整、注释漂亮、语法标准，极容易让人产生盲目信任；
- **生成 AI 的“自恋惯性”**：在同一个对话上下文里，你问写代码的那个 AI *“有没有 Bug”*，它往往带有思维惯性与自恋幻觉，大概率自信满满地回复 *“逻辑完全正确、非常健壮”*；
- **隐蔽暗坑极难肉眼捕捉**：比如在某个 `if err != nil` 提前返回时漏了解锁（Mutex Leak）、协程内部缺少 `recover`、循环内未释放的连接……这些隐患肉眼一扫极易漏掉。

这导致了一个尴尬的死循环：**AI 负责批量生产代码，人类程序员负责在半夜线上故障时擦屁股。**

---

### 解法：把审查能力极致前置到本地终端（Pre-Commit）

市面上有许多优秀的云端与 CI/CD 级别的代码审查平台，而 **`git-air` 的定位是：在代码正式 commit 或发起 PR 之前，作为你本地终端里完全独立的“第三方审查员”。**

它完全剥离聊天上下文，像一个**坐在你工位旁边、嘴有点损但眼神极毒的资深架构师同事**，只对纯粹的代码变更（`git diff`）进行冷酷挑刺：

- ⚡ **Git 原生极简体验**：无需改变任何工作习惯，在命令行敲 `git air` 即可秒级唤起。
- 🧠 **全模型自由接入**：开箱即用预置 **Google Gemini 3.x、Anthropic Claude、xAI Grok、DeepSeek、阿里通义千问、智谱 GLM、月之暗面 Kimi、本地 Ollama、OpenAI** 等全系列模型。
- 🛡️ **智能降噪过滤**：自动剥离 `go.sum`、`package-lock.json`、`*.pb.go`、`*.gen.go` 等锁文件与自动生成代码，省 Token、防死循环。
- 🎯 **严谨架构师准则**：直奔主题，聚焦并发竞态、死锁、SQL 注入、空指针、资源泄露，并直接提供修复参考代码。
- 📊 **Token 消耗与费用实时估算**：实时捕获并统计输入/输出 Token 消耗，内置 2026 最新官方定价表精准核算每次审查成本（本地离线模型自动标注免费）。
- ✨ **AI 智能 Commit Message**：一键根据代码变更生成规范的 Conventional Commits 提交说明，支持**中/英/日/韩**等多语言自适应。
- 📋 **团队规范适配（`.airules`）**：在项目根目录放置 `.airules`，AI 自动加载团队专属架构规范并最高优先级执行。
- 🪝 **一键 Pre-commit 钩子**：一行命令挂载 Git 提交前自动拦截，在代码上库前把好第一道质量关。

---

## 🖥️ 终端效果演示

```text
$ git air

[git-air] 代码评审中... (Engine: gemini / gemini-3.7-flash)
─────────────────────────────────────────────────────────────────
#### 变更概述
在用户鉴权模块中引入了 Redis 缓存逻辑，并调整了上下文超时传递机制。

#### 详细审查意见
- [BLOCKER] internal/service/user.go:45 - 严重安全隐患: 循环体内直接拼接原生 SQL 字符串，存在 SQL 注入风险。
  // 建议修复方案:
  db.Where("username = ?", inputName).First(&user)

- [WARNING] internal/service/user.go:82 - 潜在风险: 静默丢弃了 Redis 连接异常，当缓存宕机时可能造成缓存击穿。

- [WARNING] internal/service/user.go:103 - 并发隐患: 提前 return 时漏释放互斥锁（Mutex Leak），高并发下极易引发死锁。

#### 评审结论
- 结论: [REJECT]
- 评分: 60 / 100

📊 Token: 输入 1,243 / 输出 412 | ≈ $0.0025
─────────────────────────────────────────────────────────────────
```

---

## 📦 安装方式

### 方式 1：通过 Go 一键安装（推荐）
```bash
go install github.com/unclesam-ly/git-air@latest
```

### 方式 2：源码编译
```bash
git clone https://github.com/unclesam-ly/git-air.git
cd git-air
go build -o git-air .
sudo mv git-air /usr/local/bin/
```

---

## 🚀 快速上手

### 1. 配置模型与 API Key

`git-air` 原生预置了国内外所有主流大模型的官方端点与推荐模型，只需指定 `--provider` 即可秒级切换：

```bash
# 1. Google Gemini (推荐：极速且经济)
git air config set --provider gemini --key "YOUR_KEY"

# 2. Anthropic Claude (Claude-3.7 / 3.5 Sonnet)
git air config set --provider claude --key "YOUR_KEY" --model anthropic/claude-3.7-sonnet

# 3. xAI Grok (Grok-2 / Grok-3)
git air config set --provider grok --key "YOUR_KEY" --model grok-2-latest

# 4. DeepSeek (高推理能力 / 深度思考)
git air config set --provider deepseek --key "YOUR_KEY" --model deepseek-chat

# 5. 硅基流动 (SiliconFlow)
git air config set --provider siliconflow --key "YOUR_KEY" --model deepseek-ai/DeepSeek-V3

# 6. 通义千问 (Qwen / 阿里百炼)
git air config set --provider qwen --key "YOUR_KEY" --model qwen-plus

# 7. 智谱 AI (GLM-4)
git air config set --provider zhipu --key "YOUR_KEY" --model glm-4-plus

# 8. 月之暗面 (Kimi / Moonshot)
git air config set --provider moonshot --key "YOUR_KEY"

# 9. Groq (极致硬件加速推理)
git air config set --provider groq --key "YOUR_KEY" --model llama-3.3-70b-versatile

# 10. OpenRouter (全聚合模型网关)
git air config set --provider openrouter --key "YOUR_KEY"

# 11. 本地 100% 离线隐私模型 (Ollama)
git air config set --provider ollama --model qwen2.5-coder

# 12. OpenAI 官方 (GPT-4o / o3-mini)
git air config set --provider openai --key "YOUR_KEY" --model gpt-4o-mini

# 13. 自定义 Token 计费费率 (可选：覆盖内置价格表，单位: 美元/1M Tokens)
git air config set --price-input 0.75 --price-output 3.75
```

查看当前配置状态（含自定义费率与脱敏 Key）：
```bash
git air config get
```

---

### 2. 常用评审命令

```bash
# 评审当前已暂存（git add）的代码（默认行为）
git air

# 评审上一次 Commit 的代码
git air HEAD~1

# 评审特性分支与主干分支之间的所有差异
git air main..feature-agent

# 只评审指定文件的改动
git air internal/service/chat.go

# 临时指定模型或 Key 进行评审
git air --provider deepseek --model deepseek-chat
```

---

### 3. AI 智能生成规范 Commit Message (`git air msg`)

写代码一时爽，写 Commit 抓耳挠腮？`git-air` 会根据你的代码改动，自动生成符合国际 Conventional Commits 规范的精炼提交说明，并支持**中/英/日/韩多语言智能自适应**与**自动暂存**：

```bash
# 1. 基础用法：分析暂存区代码并交互式确认提交
git air msg

# 2. 自动暂存所有改动并生成提交（无需手动 git add）
git air msg -a

# 3. 指定输出自然语言（支持 auto, zh, en, ja, ko）
git air msg -l en    # 生成英文说明（适合 GitHub 开源项目）
git air msg -l ja    # 生成日语说明
git air msg -l ko    # 生成韩语说明

# 4. 极速一条龙：自动暂存 + 生成 + 免确认直接提交
git air msg -a -c
```

**效果演示：**
```text
$ git air msg -l zh
[git-air] 正在生成 Conventional Commit Message... (Language: 简体中文 / Engine: gemini-3.7-flash)
─────────────────────────────────────────────────────────────────
feat(auth): 引入 Redis 缓存并修复 Mutex 死锁隐患

- 在用户认证模块新增 Redis Token 缓存机制
- 修复提前 return 时未释放互斥锁的并发隐患
─────────────────────────────────────────────────────────────────
? 是否直接以此信息执行 git commit？[Y/n]: y
[SUCCESS] 代码提交成功！🎉
```

---

### 4. 一键安装 Git Pre-commit 拦截钩子

在你的 Git 项目根目录下执行：
```bash
git air hook install
```
- **自动门禁阻断**：当审查发现 `[BLOCKER]` 严重缺陷或 `[REJECT]` 评审结论时，`git-air` 会自动以非零退出码**直接拦截并阻断本次 `git commit`**，死守代码质量底线！
- **严格模式 (`--strict`)**：若希望发现 `[WARNING]` 也强制阻断提交，可使用 `git air --strict`；
- **临时跳过阻断**：如确需强行提交，可使用 Git 原生命令 `git commit --no-verify` 或附加 `git-air --no-block`。

如需卸载钩子：
```bash
git air hook uninstall
```

---

### 5. 版本查看与一键自动升级 (`git air update`)

`git-air` 内置了**轻量级 24 小时静默更新检测**，当 GitHub 发布新版本时会在终端底部自动提醒；你也可以随时手动检查与一键升级：

```bash
# 查看当前版本并在线检查更新
git air version

# 一键自动拉取并升级到最新版本
git air update
```

---

## ⚙️ 配置文件说明

`git-air` 支持多层配置覆盖，优先级从高到低为：
1. **命令行 Flags**（如 `--key`, `--model`, `--provider`, `--lang`, `--strict`, `--no-block`）
2. **环境变量**（`GIT_AIR_API_KEY`, `GIT_AIR_PROVIDER`, `GIT_AIR_MODEL`, `GIT_AIR_COMMIT_LANG`）
3. **项目级配置**（当前 Git 项目根目录下的 `config.yaml` 或 `.git-air.yaml`）
4. **全局用户配置**（`~/.git-air/config.yaml`）

### `config.yaml` 示例模板：
```yaml
# 模型提供商: gemini, deepseek, ollama, openai, custom
provider: "gemini"

# API 密钥 (本地 ollama 可留空)
api_key: "YOUR_API_KEY_HERE"

# 模型名称
model: "gemini-3.7-flash"

# 自定义端点 (可选)
# base_url: "https://api.deepseek.com/v1"

# 自定义基础提示词 (可选，留空则使用默认资深架构师准则)
# custom_prompt: ""

# 自定义模型单价 (可选，单位: 美元/1M Tokens)
# price_input: 0.75
# price_output: 3.75

# AI Commit Message 生成语言 (可选: auto, zh, en, ja, ko)
# commit_lang: "auto"
```

---

## 📋 团队规则定制（`.airules`）

在任何 Git 仓库的根目录下创建一个 `.airules` 文件，`git-air` 在评审时会自动将其作为**团队最高优先级铁律**注入大模型：

```markdown
# 团队专属开发铁律
1. 所有 SQL 操作必须下沉至 DAO 层，严禁裸拼 SQL 字符串；
2. 并发操作共享变量必须加锁或使用原子包（sync/atomic）；
3. 打开的连接/文件句柄必须紧跟 defer Close() 释放；
4. 严禁使用下划线 _ 盲目丢弃未处理的 error。
```

---

## 🛡️ 自定义忽略审查（`.airignore`）

除了内置自动忽略的锁文件与二进制外，你还可以在仓库根目录创建 `.airignore` 文件，自定义排除不需要 AI 评审的文件或目录：

```gitignore
# 忽略所有文档与原型设计
docs/
*.md

# 忽略测试 Mock 与假数据
tests/mock/
testdata/

# 忽略私密配置文件
*.env
secrets.yaml
```

---

## 📄 开源许可证

本项目采用 [MIT 许可证](LICENSE)。欢迎提交 Issue 与 Pull Request！

