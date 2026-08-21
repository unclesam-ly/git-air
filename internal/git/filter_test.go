package git

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsIgnoredFile(t *testing.T) {
	customPatterns := []string{
		"docs/",
		"tests/mock/*.json",
		"secrets.env",
	}

	tests := []struct {
		filePath string
		expected bool
	}{
		// 默认内置忽略规则
		{"go.sum", true},
		{"package-lock.json", true},
		{"api/user.pb.go", true},
		{"web/dist/app.min.js", true},
		{"assets/logo.png", true},

		// 用户自定义 .airignore 规则
		{"docs/readme.md", true},
		{"docs/api/v1.md", true},
		{"tests/mock/user.json", true},
		{"secrets.env", true},

		// 不应被忽略的正常业务代码
		{"cmd/root.go", false},
		{"internal/service/user.go", false},
		{"pkg/utils/helper.go", false},
	}

	for _, tt := range tests {
		got := isIgnoredFile(tt.filePath, customPatterns)
		if got != tt.expected {
			t.Errorf("isIgnoredFile(%q) = %v, expected %v", tt.filePath, got, tt.expected)
		}
	}
}

func TestLoadAirIgnore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-air-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	content := `# 忽略文档
			docs/
			# 忽略临时文件
			*.tmp
			*.log
		`
	ignoreFile := filepath.Join(tempDir, ".airignore")
	if err := os.WriteFile(ignoreFile, []byte(content), 0644); err != nil {
		t.Fatalf("写入 .airignore 失败: %v", err)
	}

	patterns := LoadAirIgnore(tempDir)
	expected := []string{"docs/", "*.tmp", "*.log"}

	if len(patterns) != len(expected) {
		t.Fatalf("LoadAirIgnore 返回数量不符: got %v, expected %v", patterns, expected)
	}

	for i, p := range patterns {
		if p != expected[i] {
			t.Errorf("Pattern[%d] = %q, expected %q", i, p, expected[i])
		}
	}
}

func TestFilterDiff(t *testing.T) {
	rawDiff := `diff --git a/go.sum b/go.sum
			index 111..222 100644
			--- a/go.sum
			+++ b/go.sum
			@@ -1 +1 @@
			+some-package v1.0.0 h1:...
			diff --git a/internal/service/user.go b/internal/service/user.go
			index 333..444 100644
			--- a/internal/service/user.go
			+++ b/internal/service/user.go
			@@ -1 +1 @@
			+func HandleUser() {}
			diff --git a/docs/spec.md b/docs/spec.md
			index 555..666 100644
			--- a/docs/spec.md
			+++ b/docs/spec.md
			@@ -1 +1 @@
			+# API Spec
			`
	customPatterns := []string{"docs/"}
	filtered := FilterDiff(rawDiff, customPatterns...)

	if strings.Contains(filtered, "go.sum") {
		t.Errorf("FilterDiff 应当过滤掉 go.sum")
	}
	if strings.Contains(filtered, "docs/spec.md") {
		t.Errorf("FilterDiff 应当过滤掉 docs/spec.md")
	}
	if !strings.Contains(filtered, "internal/service/user.go") {
		t.Errorf("FilterDiff 应当保留正常的业务代码 internal/service/user.go")
	}
}
