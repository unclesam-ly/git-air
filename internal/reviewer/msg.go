package reviewer

import (
	"fmt"
	"os"
	"strings"
)

// DetectLanguage 智能解析目标自然语言
// 优先级：传入指定语言 -> 系统 LANG/LC_ALL 环境变量 -> 默认 "en"
func DetectLanguage(specifiedLang string) string {
	lang := strings.ToLower(strings.TrimSpace(specifiedLang))
	if lang != "" && lang != "auto" {
		return NormalizeLang(lang)
	}

	// 尝试从系统环境变量中探测
	envLang := os.Getenv("LC_ALL")
	if envLang == "" {
		envLang = os.Getenv("LANG")
	}
	if envLang == "" {
		envLang = os.Getenv("LC_MESSAGES")
	}

	envLang = strings.ToLower(envLang)
	if strings.HasPrefix(envLang, "zh") {
		return "zh"
	} else if strings.HasPrefix(envLang, "ja") {
		return "ja"
	} else if strings.HasPrefix(envLang, "ko") {
		return "ko"
	}

	return "en"
}

// NormalizeLang 规范化语言缩写
func NormalizeLang(lang string) string {
	switch strings.ToLower(lang) {
	case "zh", "cn", "chinese", "zh-cn", "zh-tw", "中文":
		return "zh"
	case "ja", "jp", "japanese", "ja-jp", "日本語", "日文":
		return "ja"
	case "ko", "kr", "korean", "ko-kr", "한국어", "韩文":
		return "ko"
	default:
		return "en"
	}
}

// BuildCommitMsgPrompt 构建用于生成 Conventional Commits 的专业提示词
func BuildCommitMsgPrompt(lang string) string {
	langName := "English"
	exampleSummary := "feat(auth): add Redis token caching and fix mutex deadlock"
	exampleBullet1 := "- Introduce Redis caching layer for user authentication"
	exampleBullet2 := "- Fix concurrency bug where mutex was not unlocked on early return"

	switch lang {
	case "zh":
		langName = "简体中文 (Simplified Chinese)"
		exampleSummary = "feat(auth): 引入 Redis 缓存并修复 Mutex 死锁隐患"
		exampleBullet1 = "- 在用户认证模块新增 Redis Token 缓存机制"
		exampleBullet2 = "- 修复提前 return 时未释放互斥锁的并发隐患"
	case "ja":
		langName = "日本語 (Japanese)"
		exampleSummary = "feat(auth): Redisキャッシュを追加し、Mutexのデッドロックを修正"
		exampleBullet1 = "- ユーザー認証ハンドラーに Redis トークンキャッシュ層を導入"
		exampleBullet2 = "- 早期リターン時に Mutex の Unlock が漏れる並行処理のバグを修正"
	case "ko":
		langName = "한국어 (Korean)"
		exampleSummary = "feat(auth): Redis 캐시 도입 및 Mutex 데드락 버그 수정"
		exampleBullet1 = "- 사용자 인증 핸들러에 Redis 토큰 캐시 계층 추가"
		exampleBullet2 = "- 조기 반환 시 Mutex가 해제되지 않는 동시성 문제 해결"
	}

	return fmt.Sprintf(`你是一位精通 Git 版本控制与 Conventional Commits 规范的资深架构师。
你的任务是根据开发者提供的 Git Diff 增量代码，生成一条极其规范、精炼、清晰的 Git Commit Message。

### 规则要求：
1. 必须使用 Conventional Commits 规范：<type>(<scope>): <subject>
   - type 可选：feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert
   - scope 为本次变动的核心模块或文件名（英文小写，如 auth, db, ui, config 等）
2. 【语言要求】：除 type(scope) 以外的主题与说明必须使用【%s】撰写！
3. 输出结构：
   - 第一行：精炼的主题行（不超过 60 个字符，不带句号）
   - 第二行：留空行
   - 第三行及之后（可选）：使用 1~3 个无序列表项 (-) 简述核心改动要点
4. 严禁任何 Markdown 代码块包裹，严禁任何客套话或多余解释，直接输出纯文本 Commit Message！

### 示例参考：
%s

%s
%s
`, langName, exampleSummary, exampleBullet1, exampleBullet2)
}
