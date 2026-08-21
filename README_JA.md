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
- 🧠 **主要モデルを完全網羅**：**Google Gemini 3.x、Anthropic Claude、xAI Grok、DeepSeek、Qwen、智譜 GLM、Kimi、ローカルOllama、OpenAI** などを設定なしで即座に切り替え可能。
- 🛡️ **スマートなノイズ除去**：`go.sum` や `package-lock.json`、自動生成された `*.pb.go` を自動除外。Token消費を抑え、思考ループを防ぎます。
- 🎯 **厳格なシニアエンジニア基準**：お世辞や無駄口は一切なし。並行処理の競合、デッドロック、SQLインジェクション、NULLポインタ、リソースリークを指摘し、具体的な修正コードを直接提示します。
- 📊 **Token消費量と費用のリアルタイム見積もり**：入出力Token数を正確に集計し、各社最新の公式料金表に基づいて実行コストを即座に計算（ローカル完全オフラインモデルは無料表示）。
- 📋 **チーム規約の読み込み（`.airules`）**：リポジトリ直下に `.airules` を配置するだけで、チーム固有のアーキテクチャ規約を最優先で適用。
- 🪝 **ワンクリックで Pre-commit フック化**：コマンド1発でGitフックとして登録。危険なコードのコミットを未然にブロックします。

---

## 🖥️ ターミナル実行イメージ

```text
$ git air

[git-air] コードレビュー中... (Engine: gemini / gemini-3.7-flash)
─────────────────────────────────────────────────────────────────
#### 変更概要
ユーザー認証ハンドラーに Redis キャッシュ層を追加し、Context のタイムアウト伝播を修正。

#### 詳細レビュー結果
- [BLOCKER] internal/service/user.go:45 - 重大なセキュリティ脆弱性: ループ内で生SQL文字列を直接結合しています。SQLインジェクションのリスクがあります。
  // 推奨修正コード:
  db.Where("username = ?", inputName).First(&user)

- [WARNING] internal/service/user.go:82 - 潜在的リスク: Redis 接続エラーを握りつぶしています。キャッシュ障害時にDBへ負荷が集中します。

- [WARNING] internal/service/user.go:103 - 並行処理の欠陥: 早期リターン時に Mutex の Unlock を忘れています（Mutex Leak）。デッドロックの危険性があります。

#### 結論
- 判定: [REJECT]
- スコア: 60 / 100

📊 Token: 入力 1,243 / 出力 412 | ≈ $0.0025
─────────────────────────────────────────────────────────────────
```

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

# 13. カスタム Token 料金設定 (任意：組み込み料金表を上書き、単位: 米ドル/1M Tokens)
git air config set --price-input 0.75 --price-output 3.75
```

現在の設定を確認（カスタム料金やマスク化Keyを含む）：
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
- **コミット自動ブロック機能**：レビューで `[BLOCKER]` 重大欠陥や `[REJECT]` 判定が検出された場合、`git-air` は非ゼロの終了コードで**コミット（`git commit`）を自動的に中断・ブロック**します！
- **厳格モード (`--strict`)**：`[WARNING]` 警告でもコミットをブロックしたい場合は `git air --strict` を使用；
- **一時的なブロック解除**：緊急でコミットを強制したい場合は、Git 標準の `git commit --no-verify` または `git air --no-block` を付与します。

フックを解除する場合：
```bash
git air hook uninstall
```

---

## ⚙️ 設定の優先順位

1. **コマンドライン引数** (`--key`, `--model`, `--provider`, `--prompt`, `--price-input`, `--price-output`, `--strict`, `--no-block`)
2. **環境変数** (`GIT_AIR_API_KEY`, `GIT_AIR_PROVIDER`, `GIT_AIR_MODEL`)
3. **プロジェクト個別設定** (リポジトリ直下の `./config.yaml` または `./.git-air.yaml`)
4. **グローバル設定** (`~/.git-air/config.yaml`)

### `config.yaml` 設定ファイル例：
```yaml
# プロバイダー: gemini, claude, grok, deepseek, qwen, zhipu, moonshot, siliconflow, ollama, openai, custom
provider: "gemini"
api_key: "YOUR_API_KEY_HERE"
model: "gemini-3.7-flash"

# カスタム API エンドポイント (任意)
# base_url: "https://api.deepseek.com/v1"

# カスタム System Prompt (任意)
# custom_prompt: ""

# カスタム Token 料金設定 (任意、単位: 米ドル / 1M Tokens)
# price_input: 0.75
# price_output: 3.75
```

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

## 🛡️ カスタム除外設定（`.airignore`）

組み込みの自動除外（ロックファイルやバイナリ）に加え、リポジトリ直下に `.airignore` を配置することで、任意のファイルやディレクトリをAIレビューから除外できます：

```gitignore
# ドキュメントや設計書を除外
docs/
*.md

# テストモックやダミーデータを除外
tests/mock/
testdata/

# 機密設定ファイルを除外
*.env
secrets.yaml
```

---

## 📄 ライセンス

本プロジェクトは [MIT ライセンス](LICENSE) のもとで公開されています。Issue や Pull Request も大歓迎です！
