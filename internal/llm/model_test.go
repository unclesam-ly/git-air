package llm

import (
	"fmt"
	"testing"
)

func TestIsReasoningModel(t *testing.T) {
	tests := []struct {
		model    string
		expected bool
	}{
		{"o1", true},
		{"o1-mini", true},
		{"o1-preview", true},
		{"o3-mini", true},
		{"gpt-5.6-luna", true},
		{"gpt-5.6-sol", true},
		{"gpt-5.6-terra", true},
		{"deepseek-reasoner", true},
		{"gemini-3.7-flash", false},
		{"claude-3.7-sonnet", false},
		{"gpt-4o", false},
	}

	for _, tt := range tests {
		got := isReasoningModel(tt.model)
		if got != tt.expected {
			t.Errorf("isReasoningModel(%q) = %v, expected %v", tt.model, got, tt.expected)
		}
	}
}

func TestIsTemperatureError(t *testing.T) {
	tests := []struct {
		errStr   string
		expected bool
	}{
		{"this model has beta-limitations, temperature, top_p and n are fixed at 1", true},
		{"temperature must be 1.0", true},
		{"invalid api key", false},
		{"rate limit exceeded", false},
	}

	for _, tt := range tests {
		got := isTemperatureError(fmt.Errorf("%s", tt.errStr))
		if got != tt.expected {
			t.Errorf("isTemperatureError(%q) = %v, expected %v", tt.errStr, got, tt.expected)
		}
	}
}
