package llm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// Config 大模型统一配置
type Config struct {
	Provider string // "gemini", "deepseek", "openai", "ollama", "custom"
	APIKey   string // 密钥 (Ollama 本地可填任意字符或留空)
	BaseURL  string // 自定义接口地址 (非必填，会自动根据 Provider 设置默认值)
	Model    string // 模型名称 (非必填，会自动设置默认主力模型)
}

// Client 统一大模型客户端
type Client struct {
	client *openai.Client
	model  string
}

// NewClient 初始化客户端
func NewClient(cfg Config) (*Client, error) {
	provider := strings.ToLower(strings.TrimSpace(cfg.Provider))
	if provider == "" {
		provider = "gemini"
	}

	apiKey := cfg.APIKey
	baseURL := cfg.BaseURL
	model := cfg.Model

	// 智能适配各大主流模型的官方地址与主力模型
	switch provider {
	case "gemini", "google":
		if baseURL == "" {
			baseURL = "https://generativelanguage.googleapis.com/v1beta/openai/"
		}
		if model == "" {
			model = "gemini-3.7-flash"
		}
		if apiKey == "" {
			return nil, fmt.Errorf("gemini 必须配置 api_key")
		}

	case "deepseek":
		if baseURL == "" {
			baseURL = "https://api.deepseek.com/v1"
		}
		if model == "" {
			model = "deepseek-chat"
		}
		if apiKey == "" {
			return nil, fmt.Errorf("deepseek 必须配置 api_key")
		}

	case "siliconflow", "silicon":
		if baseURL == "" {
			baseURL = "https://api.siliconflow.cn/v1"
		}
		if model == "" {
			model = "deepseek-ai/DeepSeek-V3"
		}
		if apiKey == "" {
			return nil, fmt.Errorf("siliconflow 必须配置 api_key")
		}

	case "qwen", "dashscope", "aliyun":
		if baseURL == "" {
			baseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
		}
		if model == "" {
			model = "qwen-plus"
		}
		if apiKey == "" {
			return nil, fmt.Errorf("qwen/dashscope 必须配置 api_key")
		}

	case "zhipu", "glm":
		if baseURL == "" {
			baseURL = "https://open.bigmodel.cn/api/paas/v4/"
		}
		if model == "" {
			model = "glm-4-plus"
		}
		if apiKey == "" {
			return nil, fmt.Errorf("zhipu 必须配置 api_key")
		}

	case "moonshot", "kimi":
		if baseURL == "" {
			baseURL = "https://api.moonshot.cn/v1"
		}
		if model == "" {
			model = "moonshot-v1-8k"
		}
		if apiKey == "" {
			return nil, fmt.Errorf("moonshot/kimi 必须配置 api_key")
		}

	case "doubao", "volcengine":
		if baseURL == "" {
			baseURL = "https://ark.cn-beijing.volces.com/api/v3"
		}
		if model == "" {
			model = "doubao-pro-32k"
		}
		if apiKey == "" {
			return nil, fmt.Errorf("doubao 必须配置 api_key (并在 model 中填写接入点 Endpoint ID)")
		}

	case "minimax":
		if baseURL == "" {
			baseURL = "https://api.minimax.chat/v1"
		}
		if model == "" {
			model = "abab6.5s-chat"
		}
		if apiKey == "" {
			return nil, fmt.Errorf("minimax 必须配置 api_key")
		}

	case "yi", "lingyi":
		if baseURL == "" {
			baseURL = "https://api.lingyiwanwu.com/v1"
		}
		if model == "" {
			model = "yi-lightning"
		}
		if apiKey == "" {
			return nil, fmt.Errorf("yi 必须配置 api_key")
		}

	case "groq":
		if baseURL == "" {
			baseURL = "https://api.groq.com/openai/v1"
		}
		if model == "" {
			model = "llama-3.3-70b-versatile"
		}
		if apiKey == "" {
			return nil, fmt.Errorf("groq 必须配置 api_key")
		}

	case "grok", "xai":
		if baseURL == "" {
			baseURL = "https://api.x.ai/v1"
		}
		if model == "" {
			model = "grok-2-latest"
		}
		if apiKey == "" {
			return nil, fmt.Errorf("grok/xai 必须配置 api_key")
		}

	case "claude", "anthropic":
		if baseURL == "" {
			baseURL = "https://openrouter.ai/api/v1"
		}
		if model == "" {
			model = "anthropic/claude-3.7-sonnet"
		}
		if apiKey == "" {
			return nil, fmt.Errorf("claude 必须配置 api_key (默认走 OpenRouter 统一网关，或通过 base_url 指定自建代理端点)")
		}

	case "openrouter":
		if baseURL == "" {
			baseURL = "https://openrouter.ai/api/v1"
		}
		if model == "" {
			model = "anthropic/claude-3.7-sonnet"
		}
		if apiKey == "" {
			return nil, fmt.Errorf("openrouter 必须配置 api_key")
		}

	case "ollama":
		if baseURL == "" {
			baseURL = "http://localhost:11434/v1"
		}
		if model == "" {
			model = "qwen2.5-coder"
		}
		if apiKey == "" {
			apiKey = "ollama" // 本地 ollama 传非空占位符即可
		}

	case "lmstudio", "vllm":
		if baseURL == "" {
			baseURL = "http://localhost:1234/v1"
		}
		if model == "" {
			model = "local-model"
		}
		if apiKey == "" {
			apiKey = "local"
		}

	case "openai":
		if baseURL == "" {
			baseURL = "https://api.openai.com/v1"
		}
		if model == "" {
			model = "gpt-4o-mini"
		}
		if apiKey == "" {
			return nil, fmt.Errorf("openai 必须配置 api_key")
		}

	default:
		// custom 自定义端点（只要提供了 baseURL 就能直连）
		if baseURL == "" {
			return nil, fmt.Errorf("未知 Provider: '%s'，如需自定义请提供 base_url", provider)
		}
		if model == "" {
			model = "gpt-4o"
		}
	}

	clientConfig := openai.DefaultConfig(apiKey)
	clientConfig.BaseURL = baseURL

	return &Client{
		client: openai.NewClientWithConfig(clientConfig),
		model:  model,
	}, nil
}

// Review 一次性完整返回代码评审结果
func (c *Client) Review(ctx context.Context, systemPrompt, diff string) (string, error) {
	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: diff},
		},
		Temperature: 0.2, // 代码评审低随机性
	})

	if err != nil {
		return "", fmt.Errorf("代码评审调用失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("模型未返回任何有效评审结果")
	}

	return resp.Choices[0].Message.Content, nil
}

// ReviewStream 流式逐字返回评审结果（打字机效果）
func (c *Client) ReviewStream(ctx context.Context, systemPrompt, diff string, onChunk func(chunk string)) error {
	stream, err := c.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: diff},
		},
		Temperature: 0.2,
		Stream:      true,
	})
	if err != nil {
		return fmt.Errorf("创建流式请求失败: %w", err)
	}
	if stream != nil {
		defer stream.Close()
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("流式读取异常: %w", err)
		}
		if len(response.Choices) > 0 {
			chunk := response.Choices[0].Delta.Content
			if chunk != "" {
				onChunk(chunk)
			}
		}
	}

	return nil
}
