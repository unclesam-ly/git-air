package reviewer

import (
	"strings"
	"testing"
)

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"zh", "zh"},
		{"ZH", "zh"},
		{"zh-CN", "zh"},
		{"中文", "zh"},
		{"ja", "ja"},
		{"ja-JP", "ja"},
		{"Japanese", "ja"},
		{"日本語", "ja"},
		{"ko", "ko"},
		{"ko-KR", "ko"},
		{"한국어", "ko"},
		{"en", "en"},
		{"en-US", "en"},
		{"English", "en"},
		{"fr", "en"}, // 未预置语言回退到 en
	}

	for _, tt := range tests {
		got := DetectLanguage(tt.input)
		if got != tt.expected {
			t.Errorf("DetectLanguage(%q) = %q, expected %q", tt.input, got, tt.expected)
		}
	}
}

func TestBuildCommitMsgPrompt(t *testing.T) {
	langs := []struct {
		lang     string
		expected string
	}{
		{"zh", "简体中文"},
		{"ja", "日本語"},
		{"ko", "한국어"},
		{"en", "English"},
	}

	for _, tt := range langs {
		prompt := BuildCommitMsgPrompt(tt.lang)
		if !strings.Contains(prompt, tt.expected) {
			t.Errorf("BuildCommitMsgPrompt(%q) 应包含语言名称 %q", tt.lang, tt.expected)
		}
		if !strings.Contains(prompt, "Conventional Commits") {
			t.Errorf("BuildCommitMsgPrompt(%q) 应包含 Conventional Commits 说明", tt.lang)
		}
	}
}

func TestBuildCommitMsgPromptWithCustom(t *testing.T) {
	prompt := BuildCommitMsgPromptWithCustom("en", "Use a short imperative subject and mention breaking changes explicitly.")
	if !strings.Contains(prompt, "Use a short imperative subject") {
		t.Fatal("自定义 Commit Message Prompt 未被合并")
	}
	if !strings.Contains(prompt, "Conventional Commits") {
		t.Fatal("内置 Conventional Commits 约束不应丢失")
	}
}
