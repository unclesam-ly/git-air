<p align="center">
  <img src="assets/logo.png" width="180" height="180" alt="git-air logo" />
</p>

<h1 align="center">git-air 🍃</h1>

<p align="center">
  <strong>Go言語で書かれた、空気のように軽快で高速なGitネイティブAIコードレビューツール。</strong>
</p>

<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.26%2B-00ADD8?style=flat&logo=go" alt="Go Version" /></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License" /></a>
  <a href="https://github.com/unclesam-ly/git-air/pulls"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg" alt="PRs Welcome" /></a>
</p>

<p align="center">
  <a href="README.md">简体中文</a> |
  <a href="README_EN.md">English</a> |
  <a href="README_JA.md"><strong>日本語</strong></a>
</p>

---

## 🌟 なぜ git-air を開発したのか？

### 課題：AIが書いたコードは気持ちいいが、レビューは悪夢

現代の開発では、AI（Cursor、Copilot、各種LLMプラグイン）を活用することが当たり前になりました。数秒で数十行・数百行のコードが生成されるため、開発スピードは劇的に向上します。

**しかし、そこには見過ごせない重大な課題が存在します：**
- **「見た目が完璧すぎる」という錯覚**：インデントやコメント、命名が綺麗なため、バグがないと盲信しやすい；
- **生成AI自身の「自己肯定バイアス」**：同じチャット内で生成AIに「このコードにバグはある？」と聞いても、「問題なく堅牢に動作します！」と返ってくることが多い；
- **人間の目で見落としやすい潜在的な罠**：`if err != nil` の早期リターン時にロック解除（Mutex Leak）を忘れたり、Goroutine 内で `recover` がなかったり、ループ内でコネクションが解放されなかったりするバグは、目視では非常に見落としやすい。

その結果、**「AIがコードを大量生産し、深夜に人間が本番障害の対応に追われる」**という本末転倒な状況が生まれます。

---

### 解決策：レビューをローカルターミナル（Pre-Commit）に完全前置

GitHubなどのCI/CD上で動作するボットも優秀ですが、**`git-air` はコードを `git commit` する前に、ローカル端末上で動作する完全に独立した「第三者監査役」として設計されています。**

チャット履歴のコンテキストを一切排除し、まるで**あなたの隣に座る口は悪いが目の鋭いシニアアーキテクト**のように、純粋な差分（`git diff`）だけを厳しくチェックします：

- ⚡ **Gitネイティブな極上の操作性**：いつもの作業フローを崩さず、ターミナルで `git air` と叩くだけで即座に起動。
- 🧠 **主要モデルを完全網羅**：**Google Gemini、Anthropic Claude、xAI Grok、DeepSeek、ローカルOllama、OpenAI** などを設定なしで即座に切り替え可能。
- 🛡️ **スマートなノイズ除去**：`go.sum` や `package-lock.json`、自動生成された `*.pb.go` を自動除外。Token消費を抑え、思考ループを防ぎます。
- 🎯 **厳格なシニアエンジニア基準**：お世辞や無駄口は一切なし。並行処理の競合、デッドロック、SQLインジェクション、NULLポインタ、リソースリークを指摘し、具体的な修正コードを直接提示します。
- 📋 **チーム規約の読み込み（`.airules`）**：リポジトリ直下に `.airules` を配置するだけで、チーム固有のアーキテクチャ規約を最優先で適用。
- 🪝 **ワンクリックで Pre-commit フック化**：コマンド1発でGitフックとして登録。危険なコードのコミットを未然にブロックします。

---

## 📦 インストール

### 方法 1: Go Install（推奨）
```bash
go install github.com/unclesam-ly/git-air@latest
```

### 方法 2: ソースコードからビルド
```bash
git clone https://github.com/unclesam-ly/git-air.git
cd git-air
go build -o git-air .
sudo mv git-air /usr/local/bin/
```

---

## 🚀 クイックスタート

### 1. モデルと API Key の設定

`git-air` は国内外の主要LLMのエンドポイントと推奨モデルを標準プリセットしています。`--provider` を指定するだけで簡単に切り替えられます：

```bash
# 1. Google Gemini (推奨: 高速・低コスト)
git air config set --provider gemini --key "YOUR_KEY"

# 2. Anthropic Claude (Claude-3.7 / 3.5 Sonnet)
git air config set --provider claude --key "YOUR_KEY" --model anthropic/claude-3.7-sonnet

# 3. xAI Grok (Grok-2 / Grok-3)
git air config set --provider grok --key "YOUR_KEY" --model grok-2-latest

# 4. DeepSeek (高度な推論・思考モデル)
git air config set --provider deepseek --key "YOUR_KEY" --model deepseek-chat

# 5. SiliconFlow (DeepSeek-V3 高速ホスティング)
git air config set --provider siliconflow --key "YOUR_KEY" --model deepseek-ai/DeepSeek-V3

# 6. 通義千問 (Qwen / Alibaba Cloud)
git air config set --provider qwen --key "YOUR_KEY" --model qwen-plus

# 7. 智譜 AI (GLM-4)
git air config set --provider zhipu --key "YOUR_KEY" --model glm-4-plus

# 8. Moonshot AI (Kimi)
git air config set --provider moonshot --key "YOUR_KEY"

# 9. Groq (超高速推論)
git air config set --provider groq --key "YOUR_KEY" --model llama-3.3-70b-versatile

# 10. OpenRouter (マルチモデル統合ゲートウェイ)
git air config set --provider openrouter --key "YOUR_KEY"

# 11. ローカル完全オフライン・プライバシー重視 (Ollama)
git air config set --provider ollama --model qwen2.5-coder

# 12. OpenAI 公式 (GPT-4o / o3-mini)
git air config set --provider openai --key "YOUR_KEY" --model gpt-4o-mini
```

現在の設定を確認：
```bash
git air config get
```

---

### 2. よく使うレビューコマンド

```bash
# 現在ステージングされた（git add済み）変更をレビュー（デフォルト）
git air

# 直前のコミット内容をレビュー
git air HEAD~1

# ブランチ間の差分をレビュー
git air main..feature-agent

# 特定のファイルのみをレビュー
git air internal/service/chat.go

# 一時的にモデルやProviderを指定して実行
git air --provider deepseek --model deepseek-chat
```

---

### 3. ワンタッチで Pre-commit フックを登録

Git リポジトリ直下で以下を実行：
```bash
git air hook install
```
以降、`git commit` を実行するたびにターミナル上で自動レビューが走り、重大な欠陥があればコミット前に警告を出してくれます！

フックを解除する場合：
```bash
git air hook uninstall
```

---

## ⚙️ 設定の優先順位

1. **コマンドライン引数** (`--key`, `--model`, `--provider`, `--prompt`)
2. **環境変数** (`GIT_AIR_API_KEY`, `GIT_AIR_PROVIDER`, `GIT_AIR_MODEL`)
3. **プロジェクト個別設定** (リポジトリ直下の `./config.yaml` または `./.git-air.yaml`)
4. **グローバル設定** (`~/.git-air/config.yaml`)

---

## 📋 チーム規約のカスタマイズ（`.airules`）

リポジトリのルートディレクトリに `.airules` ファイルを作成すると、AI はそのルールを最優先事項として適用します：

```markdown
# チーム共通開発ルール
1. すべての SQL 処理は DAO 層に集約し、文字列結合による生 SQL は禁止とする；
2. 並行処理で共有変数を扱う場合は、sync.RWMutex または sync/atomic を使用すること；
3. オープンしたリソース（DB rows、HTTP レスポンス、ファイル）は defer Close() で確実に解放すること；
4. error をアンダースコア（_ = err）で無条件に無視することを禁止する。
```

---

## 📄 ライセンス

本プロジェクトは [MIT ライセンス](LICENSE) のもとで公開されています。Issue や Pull Request も大歓迎です！
