package security

import "regexp"

type redactRule struct {
	pattern     *regexp.Regexp
	replacement string
}

var redactRules = []redactRule{
	{
		// PEM 私钥必须优先处理，避免后续规则只替换其中一部分。
		pattern:     regexp.MustCompile(`(?s)-----BEGIN [A-Z ]*PRIVATE KEY-----.*?-----END [A-Z ]*PRIVATE KEY-----`),
		replacement: "[REDACTED_SECRET]",
	},
	{
		// 常见 Token 前缀：OpenAI/GitHub/AWS/Slack 等。
		pattern:     regexp.MustCompile(`(?i)\b(?:sk-[A-Za-z0-9]{16,}|gh[pousr]_[A-Za-z0-9_]{16,}|AKIA[0-9A-Z]{16}|xox[baprs]-[A-Za-z0-9-]{10,})\b`),
		replacement: "[REDACTED_SECRET]",
	},
	{
		// Authorization: Bearer ...
		pattern:     regexp.MustCompile(`(?i)(Authorization\s*:\s*Bearer\s+)[A-Za-z0-9._~+/=-]+`),
		replacement: "${1}[REDACTED_BEARER_TOKEN]",
	},
	{
		// JWT。
		pattern:     regexp.MustCompile(`\b[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\b`),
		replacement: "[REDACTED_SECRET]",
	},
	{
		// 数据库连接串。
		pattern:     regexp.MustCompile(`(?i)\b(?:postgres(?:ql)?|mysql|mongodb(?:\+srv)?):\/\/[^\s"'<>]+`),
		replacement: "[REDACTED_SECRET]",
	},
	{
		// 常见密钥字段赋值。值至少 8 个字符，减少普通变量误报。
		pattern:     regexp.MustCompile(`(?i)(\b(?:api[_-]?key|access[_-]?key|secret(?:[_-]?key)?|password|passwd|token)\b\s*[:=]\s*["']?)[^\s"' ,;]{8,}`),
		replacement: "${1}[REDACTED_SECRET]",
	},
}

// RedactSecrets 在 Diff 发送给模型前替换常见敏感信息。
// 脱敏只作用于发送给模型的字符串，不修改工作区文件或 Git 内容。
func RedactSecrets(input string) string {
	output := input
	for _, rule := range redactRules {
		output = rule.pattern.ReplaceAllString(output, rule.replacement)
	}
	return output
}
