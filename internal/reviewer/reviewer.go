package reviewer

import (
	"context"
	"fmt"
	"strings"

	"github.com/unclesam-ly/git-air/internal/llm"
)

type Reviewer struct {
	client       *llm.Client
	repoRoot     string
	customPrompt string
}

// NewReviewer 初始化评审引擎
// customPrompt 可传入用户自定义配置的 Prompt，传空字符串则使用默认 Prompt
func NewReviewer(client *llm.Client, repoRoot string, customPrompt string) *Reviewer {
	return &Reviewer{
		client:       client,
		repoRoot:     repoRoot,
		customPrompt: customPrompt,
	}
}

// Execute 执行代码评审，返回 (result, usage, error)
func (r *Reviewer) Execute(ctx context.Context, diff string, stream bool, onChunk func(chunk string)) (string, *llm.Usage, error) {
	if strings.TrimSpace(diff) == "" {
		return "", nil, fmt.Errorf("diff 内容为空，无需评审")
	}

	// 动态组装最终 System Prompt
	systemPrompt := BuildSystemPrompt(r.customPrompt, r.repoRoot)
	if stream && onChunk != nil {
		usage, err := r.client.ReviewStream(ctx, systemPrompt, diff, onChunk)
		return "", usage, err
	}

	result, err := r.client.Review(ctx, systemPrompt, diff)
	return result, nil, err
}

// GenerateCommitMsg 根据 Diff 智能生成符合 Conventional Commits 规范的提交信息
func (r *Reviewer) GenerateCommitMsg(ctx context.Context, diff string, lang string, customPrompt string) (string, error) {
	if strings.TrimSpace(diff) == "" {
		return "", fmt.Errorf("diff 内容为空，无法生成 Commit Message")
	}

	targetLang := DetectLanguage(lang)
	systemPrompt := BuildCommitMsgPromptWithCustom(targetLang, customPrompt)

	msg, err := r.client.Review(ctx, systemPrompt, diff)
	if err != nil {
		return "", err
	}

	// 清理多余的代码块标记或前后空白
	cleanMsg := strings.TrimSpace(msg)
	cleanMsg = strings.TrimPrefix(cleanMsg, "```gitcommit")
	cleanMsg = strings.TrimPrefix(cleanMsg, "```")
	cleanMsg = strings.TrimSuffix(cleanMsg, "```")
	return strings.TrimSpace(cleanMsg), nil
}
