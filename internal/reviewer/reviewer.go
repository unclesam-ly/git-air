package reviewer

import (
	"context"
	"fmt"
	"github.com/unclesam-ly/git-air/internal/llm"
	"strings"
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

// Execute 执行代码评审
func (r *Reviewer) Execute(ctx context.Context, diff string, stream bool, onChunk func(chunk string)) (string, error) {
	if strings.TrimSpace(diff) == "" {
		return "", fmt.Errorf("diff 内容为空，无需评审")
	}

	// 动态组装最终 System Prompt
	systemPrompt := BuildSystemPrompt(r.customPrompt, r.repoRoot)
	if stream && onChunk != nil {
		err := r.client.ReviewStream(ctx, systemPrompt, diff, onChunk)
		return "", err
	}

	return r.client.Review(ctx, systemPrompt, diff)
}
