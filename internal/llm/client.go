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
	Provider    string  // "gemini", "deepseek", "openai", "ollama", "custom"
	APIKey      string  // 密钥 (Ollama 本地可填任意字符或留空)
	BaseURL     string  // 自定义接口地址 (非必填，会自动根据 Provider 设置默认值)
	Model       string  // 模型名称 (非必填，会自动设置默认主力模型)
	PriceInput  float64 // 用户自定义每 1M 输入 Token 美元价格 (可选，覆盖内置价格表)
	PriceOutput float64 // 用户自定义每 1M 输出 Token 美元价格 (可选，覆盖内置价格表)
}

// Client 统一大模型客户端
type Client struct {
	client      *openai.Client
	model       string
	priceInput  float64
	priceOutput float64
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
		client:      openai.NewClientWithConfig(clientConfig),
		model:       model,
		priceInput:  cfg.PriceInput,
		priceOutput: cfg.PriceOutput,
	}, nil
}

// isReasoningModel 判断是否为推理或带有 Beta 限制的模型（这类模型要求 temperature 固定为 1）
func isReasoningModel(model string) bool {
	m := strings.ToLower(model)
	return strings.HasPrefix(m, "o1") ||
		strings.HasPrefix(m, "o3") ||
		strings.HasPrefix(m, "o4") ||
		strings.Contains(m, "reasoner") ||
		strings.Contains(m, "gpt-5") ||
		strings.Contains(m, "luna") ||
		strings.Contains(m, "terra") ||
		strings.Contains(m, "sol")
}

// isTemperatureError 检测 API 报错是否由 temperature 参数限制触发
func isTemperatureError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "temperature") || strings.Contains(msg, "beta-limitations")
}

// Review 一次性完整返回代码评审结果
func (c *Client) Review(ctx context.Context, systemPrompt, diff string) (string, error) {
	temp := float32(0.2)
	if isReasoningModel(c.model) {
		temp = 1.0
	}

	req := openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: diff},
		},
		Temperature: temp,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	// 针对限制了 temperature 的模型，且首次非 1.0 时，自动调整为 1.0 重试
	if err != nil && req.Temperature != 1.0 && isTemperatureError(err) {
		req.Temperature = 1.0
		resp, err = c.client.CreateChatCompletion(ctx, req)
	}

	if err != nil {
		return "", fmt.Errorf("代码评审调用失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("模型未返回任何有效评审结果")
	}

	return resp.Choices[0].Message.Content, nil
}

// Usage 记录一次调用的 Token 消耗
type Usage struct {
	InputTokens       int
	OutputTokens      int
	Model             string
	Provider          string
	CustomPriceInput  float64 // 用户自定义输入单价（每 1M tokens 美元，>0 优先使用）
	CustomPriceOutput float64 // 用户自定义输出单价（每 1M tokens 美元，>0 优先使用）
}

// reviewPriceTable 各模型每 1M token 的美元单价（输入/输出）
// 数据来源：各厂商官方定价页面，更新于 2026-08
// 如实际价格有变动，可在 config.yaml 中通过 price_input/price_output 字段自定义覆盖
var reviewPriceTable = map[string][2]float64{

	// ── Google Gemini ─────────────────────────────────────────────
	// 3.x 系列（当前主力）
	"gemini-3.7-flash":      {0.75, 3.75}, // 推介价至 2026-12-31，之后 1.50/7.50
	"gemini-3.6-flash":      {1.50, 7.50},
	"gemini-3.5-flash":      {1.50, 9.00},
	"gemini-3.5-flash-lite": {0.30, 2.50},
	"gemini-3.1-flash-lite": {0.25, 1.50},
	"gemini-3.1-pro":        {2.00, 12.00},
	// 2.5 系列（2026-10 退役）
	"gemini-2.5-pro":        {1.25, 10.00},
	"gemini-2.5-flash":      {0.30, 2.50},
	"gemini-2.5-flash-lite": {0.10, 0.40},
	// 旧款（兼容）
	"gemini-2.0-flash": {0.10, 0.40},

	// ── DeepSeek ─────────────────────────────────────────────────
	// 注意：DeepSeek 已启用峰谷定价，以下为非高峰标准价
	"deepseek-v4-flash": {0.22, 0.66},
	"deepseek-chat":     {0.22, 0.66}, // V4-Flash 别名
	"deepseek-reasoner": {0.55, 2.19},

	// ── Anthropic Claude ──────────────────────────────────────────
	"claude-opus-5":              {5.00, 25.00},
	"claude-sonnet-5":            {2.00, 10.00}, // 推介价至 2026-08-31
	"anthropic/claude-sonnet-5":  {2.00, 10.00},
	"claude-haiku-4-5":           {1.00, 5.00},
	"anthropic/claude-haiku-4-5": {1.00, 5.00},
	// 旧款（兼容）
	"anthropic/claude-3.7-sonnet": {3.00, 15.00},
	"anthropic/claude-3.5-sonnet": {3.00, 15.00},
	"anthropic/claude-3.5-haiku":  {0.80, 4.00},

	// ── xAI Grok ─────────────────────────────────────────────────
	"grok-4.6":      {2.00, 6.00}, // >200k tokens 翻倍
	"grok-4.3":      {1.25, 2.50},
	"grok-2-latest": {2.00, 6.00}, // 旧款映射

	// ── OpenAI ───────────────────────────────────────────────────
	"gpt-5.6-sol":   {5.00, 30.00},
	"gpt-5.6-terra": {2.00, 12.00},
	"gpt-5.6-luna":  {0.20, 1.20},
	"gpt-4o":        {2.50, 10.00}, // 旧款保留
	"gpt-4o-mini":   {0.15, 0.60},
	"o3-mini":       {1.10, 4.40},

	// ── SiliconFlow（聚合平台，按 DeepSeek V4 主力估算）────────────
	"deepseek-ai/DeepSeek-V3": {0.22, 0.66},
	"deepseek-ai/DeepSeek-R1": {0.55, 2.19},

	// ── Qwen / 阿里通义千问 ───────────────────────────────────────
	"qwen-3.8-max":   {2.00, 6.00},
	"qwen-max":       {2.00, 6.00},
	"qwen-plus":      {0.40, 1.20},
	"qwen-3.7-flash": {0.03, 0.13},

	// ── Kimi / Moonshot AI ────────────────────────────────────────
	"kimi-k3":          {3.00, 15.00},
	"moonshot-v1-8k":   {0.95, 4.00},
	"moonshot-v1-32k":  {1.90, 8.00},
	"moonshot-v1-128k": {3.50, 15.00},

	// ── 智谱 GLM ─────────────────────────────────────────────────
	"glm-5.3":    {1.40, 4.40},
	"glm-5.2":    {1.40, 4.40},
	"glm-4-plus": {1.40, 4.40},

	// ── Doubao / 字节跳动 ─────────────────────────────────────────
	"doubao-seed-2.1-pro": {0.88, 4.42},
	"doubao-pro-32k":      {0.88, 4.42},

	// ── Groq（超快推理，按量计费）────────────────────────────────
	"llama-3.3-70b-versatile": {0.59, 0.79},
	"llama-3.1-8b-instant":    {0.05, 0.08},

	// ── 本地模型（完全免费）──────────────────────────────────────
	"ollama":      {0, 0},
	"local-model": {0, 0},
}

// EstimateCost 估算费用（美元），返回 -1 表示无定价数据
func (u *Usage) EstimateCost() float64 {
	// 1. 优先使用用户在 config.yaml 中配置的自定义单价（每 1M tokens 美元）
	if u.CustomPriceInput > 0 || u.CustomPriceOutput > 0 {
		inputCost := float64(u.InputTokens) / 1_000_000 * u.CustomPriceInput
		outputCost := float64(u.OutputTokens) / 1_000_000 * u.CustomPriceOutput
		return inputCost + outputCost
	}

	// 2. 回退使用内置官方模型价格表
	prices, ok := reviewPriceTable[u.Model]
	if !ok {
		return -1
	}
	// Ollama / 本地模型 price = 0，直接返回 0
	inputCost := float64(u.InputTokens) / 1_000_000 * prices[0]
	outputCost := float64(u.OutputTokens) / 1_000_000 * prices[1]
	return inputCost + outputCost
}

// ReviewStream 流式逐字返回评审结果（打字机效果），同时返回 Token 使用量
func (c *Client) ReviewStream(ctx context.Context, systemPrompt, diff string, onChunk func(chunk string)) (*Usage, error) {
	temp := float32(0.2)
	if isReasoningModel(c.model) {
		temp = 1.0
	}

	req := openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: diff},
		},
		Temperature: temp,
		Stream:      true,
		StreamOptions: &openai.StreamOptions{
			IncludeUsage: true,
		},
	}

	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	// 针对限制了 temperature 的模型，且首次非 1.0 时，自动调整为 1.0 重试
	if err != nil && req.Temperature != 1.0 && isTemperatureError(err) {
		req.Temperature = 1.0
		stream, err = c.client.CreateChatCompletionStream(ctx, req)
	}

	if err != nil {
		return nil, fmt.Errorf("创建流式请求失败: %w", err)
	}
	if stream != nil {
		defer stream.Close()
	}

	usage := &Usage{
		Model:             c.model,
		CustomPriceInput:  c.priceInput,
		CustomPriceOutput: c.priceOutput,
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("流式读取异常: %w", err)
		}

		// 捕获最后一个 chunk 中携带的 Usage 信息
		if response.Usage != nil && response.Usage.TotalTokens > 0 {
			usage.InputTokens = response.Usage.PromptTokens
			usage.OutputTokens = response.Usage.CompletionTokens
		}

		if len(response.Choices) > 0 {
			chunk := response.Choices[0].Delta.Content
			if chunk != "" {
				onChunk(chunk)
			}
		}
	}

	return usage, nil
}
