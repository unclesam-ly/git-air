package reviewer

import (
	"os"
	"path/filepath"
	"strings"
)

// RuleFileCandidates 支持的团队规则文件名优先级列表
var RuleFileCandidates = []string{
	".airules",
	".git-air-rules",
	".git-air.yaml",
	".git-air.md",
}

// LoadCustomRules 读取项目根目录下的团队专属规则
func LoadCustomRules(repoRoot string) string {
	if repoRoot == "" {
		return ""
	}
	rulePath := FindRuleFile(repoRoot)
	if rulePath == "" {
		return ""
	}
	data, err := os.ReadFile(rulePath)
	if err != nil || len(data) == 0 {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// FindRuleFile 探测当前仓库根目录下是否存在匹配的规则文件，返回其绝对路径
func FindRuleFile(repoRoot string) string {
	for _, filename := range RuleFileCandidates {
		fullPath := filepath.Join(repoRoot, filename)
		if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
			return fullPath
		}

	}
	return ""
}
