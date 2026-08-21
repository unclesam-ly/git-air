<p align="center">
  <img src="assets/logo.png" width="180" height="180" alt="git-air logo" />
</p>

<h1 align="center">git-air 🍃</h1>

<p align="center">
  <strong>A lightweight, lightning-fast, Git-native AI code reviewer written in pure Go.</strong>
</p>

<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?style=flat&logo=go" alt="Go Version" /></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License" /></a>
  <a href="https://github.com/unclesam-ly/git-air/pulls"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg" alt="PRs Welcome" /></a>
</p>

<p align="center">
  <a href="README.md">简体中文</a> |
  <a href="README_EN.md"><strong>English</strong></a> |
  <a href="README_JA.md">日本語</a>
</p>

---

## 🌟 Why git-air?

### The Pain Point: AI Writes Code Fast, But Reviewing It Is a Nightmare

Most developers today rely heavily on AI coding assistants (Cursor, Copilot, and LLM extensions). Generating hundreds of lines of code in seconds feels incredible.

**However, this introduces a critical new problem:**
- **The "Looks Perfect" Illusion**: AI-generated code features neat indentation, convincing naming, and flawless syntax, making it dangerously easy to trust blindly;
- **The Confirmation Bias of Generating Models**: In the same conversation thread, asking the very same AI *"Are there any bugs here?"* almost always yields a self-flattering response: *"The logic is completely solid and robust!"*;
- **Subtle Traps Missed by Human Eyes**: Overlooked unlocked mutexes on early `return` paths, missing `recover()` inside goroutines, unreleased database connections in loops... These critical traps easily slip past quick manual skimming.

This creates an absurd cycle: **AI mass-produces code, while human engineers spend late nights fixing production crashes.**

---

### The Solution: Shift Code Review Directly to Your Local Terminal (Pre-Commit)

While there are great cloud/CI-based PR review bots, **`git-air` is designed to be your independent, cold-blooded "third-party auditor" right inside your local terminal before you ever type `git commit`.**

Completely isolated from your chat history, it acts like a **grumpy but sharp senior architect sitting beside your desk**, ruthlessly inspecting only raw code changes (`git diff`):

- ⚡ **Git-Native Simplicity**: No workflow changes needed. Just run `git air` in your terminal for instant feedback.
- 🧠 **Universal Multi-Model Support**: Pre-configured support for **Google Gemini 3.x, Anthropic Claude, xAI Grok, DeepSeek, Qwen, Zhipu GLM, Moonshot Kimi, Local Ollama, OpenAI**, and any OpenAI-compatible endpoint.
- 🛡️ **Smart Noise Reduction**: Automatically strips out `go.sum`, `package-lock.json`, `*.pb.go`, and lockfiles to save tokens and eliminate review loops.
- 🎯 **Strict Senior Architect Persona**: Zero corporate fluff. Directly flags concurrency deadlocks, SQL injections, nil dereferences, and resource leaks with copy-paste code patches.
- 📊 **Real-time Token & Cost Estimation**: Accurately tracks prompt/completion token consumption with built-in official pricing rates (local offline models automatically marked free).
- 📋 **Team Rules Support (`.airules`)**: Drop an `.airules` file in your repo root, and the AI will prioritize your team-specific guidelines.
- 🪝 **One-Click Pre-commit Hook**: Mounts as a local Git hook in one command to block dangerous bugs before code is committed.

---

## 🖥️ Terminal Preview

```text
$ git air

[git-air] Reviewing code... (Engine: gemini / gemini-3.7-flash)
─────────────────────────────────────────────────────────────────
#### Change Summary
Added Redis caching layer to user authentication and refined context propagation.

#### Detailed Findings
- [BLOCKER] internal/service/user.go:45 - Critical Vulnerability: Direct SQL string concatenation inside loop. Potential SQL injection.
  // Suggested Fix:
  db.Where("username = ?", inputName).First(&user)

- [WARNING] internal/service/user.go:82 - Potential Risk: Silently discarding Redis error. May lead to cache stampede.

- [WARNING] internal/service/user.go:103 - Concurrency Issue: Mutex unlock skipped on early return branch. High risk of deadlock.

#### Verdict
- Status: [REJECT]
- Score: 60 / 100

📊 Token: Input 1,243 / Output 412 | ≈ $0.0025
─────────────────────────────────────────────────────────────────
```

---

## 📦 Installation

### Option 1: Via Go Install (Recommended)
```bash
go install github.com/unclesam-ly/git-air@latest
```

### Option 2: Build from Source
```bash
git clone https://github.com/unclesam-ly/git-air.git
cd git-air
go build -o git-air .
sudo mv git-air /usr/local/bin/
```

---

## 🚀 Quick Start

### 1. Configure Model & API Key

`git-air` comes with pre-configured endpoints and default models for all major AI providers. Simply pass `--provider`:

```bash
# 1. Google Gemini (Recommended: fast & cost-effective)
git air config set --provider gemini --key "YOUR_KEY"

# 2. Anthropic Claude (Claude-3.7 / 3.5 Sonnet)
git air config set --provider claude --key "YOUR_KEY" --model anthropic/claude-3.7-sonnet

# 3. xAI Grok (Grok-2 / Grok-3)
git air config set --provider grok --key "YOUR_KEY" --model grok-2-latest

# 4. DeepSeek (Deep reasoning)
git air config set --provider deepseek --key "YOUR_KEY" --model deepseek-chat

# 5. SiliconFlow (DeepSeek-V3 hosted)
git air config set --provider siliconflow --key "YOUR_KEY" --model deepseek-ai/DeepSeek-V3

# 6. Qwen / DashScope (Alibaba Cloud)
git air config set --provider qwen --key "YOUR_KEY" --model qwen-plus

# 7. Zhipu AI (GLM-4)
git air config set --provider zhipu --key "YOUR_KEY" --model glm-4-plus

# 8. Moonshot AI (Kimi)
git air config set --provider moonshot --key "YOUR_KEY"

# 9. Groq (Ultra-fast inference)
git air config set --provider groq --key "YOUR_KEY" --model llama-3.3-70b-versatile

# 10. OpenRouter (Universal AI Gateway)
git air config set --provider openrouter --key "YOUR_KEY"

# 11. 100% Offline Local Privacy (Ollama)
git air config set --provider ollama --model qwen2.5-coder

# 12. OpenAI Official (GPT-4o / o3-mini)
git air config set --provider openai --key "YOUR_KEY" --model gpt-4o-mini

# 13. Custom Token Pricing (Optional: override default rates, USD/1M tokens)
git air config set --price-input 0.75 --price-output 3.75
```

Check current configuration (including masked keys and custom rates):
```bash
git air config get
```

---

### 2. Everyday Commands

```bash
# Review currently staged changes (default)
git air

# Review the previous commit
git air HEAD~1

# Review differences between branches
git air main..feature-agent

# Review a specific file
git air internal/service/chat.go

# Temporarily override provider or model
git air --provider deepseek --model deepseek-chat
```

---

### 3. One-Click Git Pre-commit Hook

Install the hook inside any Git repository:
```bash
git air hook install
```
- **Automated Commit Gate**: When `git-air` flags any `[BLOCKER]` issue or outputs a `[REJECT]` verdict, it automatically exits with a non-zero status code, **directly intercepting and blocking `git commit`**!
- **Strict Mode (`--strict`)**: To block commits even on `[WARNING]` issues, use `git air --strict`;
- **Bypass Interception**: If you urgently need to bypass the check, use standard Git `git commit --no-verify` or pass `git air --no-block`.

To uninstall:
```bash
git air hook uninstall
```

---

## ⚙️ Configuration Hierarchy

`git-air` searches for configurations in the following priority order:
1. **CLI Flags** (`--key`, `--model`, `--provider`, `--prompt`, `--price-input`, `--price-output`, `--strict`, `--no-block`)
2. **Environment Variables** (`GIT_AIR_API_KEY`, `GIT_AIR_PROVIDER`, `GIT_AIR_MODEL`)
3. **Project Config** (`./config.yaml` or `./.git-air.yaml` in repo root)
4. **Global User Config** (`~/.git-air/config.yaml`)

### `config.yaml` Example Template:
```yaml
# Provider: gemini, claude, grok, deepseek, qwen, zhipu, moonshot, siliconflow, ollama, openai, custom
provider: "gemini"
api_key: "YOUR_API_KEY_HERE"
model: "gemini-3.7-flash"

# Custom Base URL (Optional)
# base_url: "https://api.deepseek.com/v1"

# Custom System Prompt (Optional)
# custom_prompt: ""

# Custom Token Pricing (Optional, in USD per 1M tokens)
# price_input: 0.75
# price_output: 3.75
```

---

## 📋 Custom Team Rules (`.airules`)

Create an `.airules` file in your repository root. `git-air` will automatically inject these rules into the model with highest priority:

```markdown
# Team Coding Guidelines
1. All SQL operations must reside in the DAO/repository layer. No raw SQL string concatenation;
2. Concurrent shared variable access must use sync.RWMutex or sync/atomic;
3. Opened resource handles (DB rows, HTTP resp.Body, files) must be immediately closed via defer;
4. Never discard errors with blank identifiers (_ = err).
```

---

## 🛡️ Custom Ignore Rules (`.airignore`)

In addition to built-in ignored lockfiles and binaries, create an `.airignore` file in your repository root to exclude specific paths from AI review:

```gitignore
# Ignore documentation & design specs
docs/
*.md

# Ignore test mocks and fixtures
tests/mock/
testdata/

# Ignore sensitive config files
*.env
secrets.yaml
```

---

## 📄 License

This project is licensed under the [MIT License](LICENSE). Contributions, PRs, and issues are welcome!

