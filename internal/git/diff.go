package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var (
	ErrNotGitRepo = errors.New("当前目录不是 Git 仓库（未找到 .git）")
	ErrEmptyDiff  = errors.New("没有检测到任何需要评审的代码变更（Diff 为空）")
)

// IsGitRepo 检查当前工作目录是否在一个 Git 仓库内
func IsGitRepo(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}

// GetRepoRoot 获取当前 Git 仓库的根目录绝对路径
func GetRepoRoot(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", ErrNotGitRepo
	}

	return strings.TrimSpace(string(out)), nil
}

const maxDiffSize = 300 * 1024 // 限制最大 300KB Diff，避免打爆上下文与 OOM

// GetDiff 提取并过滤 Git 代码差异
// args 可传：
// - 无参数：默认先查暂存区 (git diff --cached)，若暂存区为空则自动查工作区 (git diff)
// - 自定义参数：如 "HEAD~1"、"main..feature"、"internal/service/chat.go" 等
func GetDiff(ctx context.Context, args ...string) (string, error) {
	if !IsGitRepo(ctx) {
		return "", ErrNotGitRepo
	}

	var cmdArgs []string
	if len(args) == 0 {
		// 1. 默认先提取暂存区代码 (准备 commit 的代码)
		stagedDiff, err := execGitDiff(ctx, "--cached")
		if err != nil {
			return "", err
		}
		if strings.TrimSpace(stagedDiff) != "" {
			return truncateDiffIfNeeded(FilterDiff(stagedDiff)), nil
		}
		// 2. 如果暂存区为空，回退提取工作区未暂存的修改
		cmdArgs = []string{}
	} else {
		cmdArgs = args
	}

	rawDiff, err := execGitDiff(ctx, cmdArgs...)
	if err != nil {
		return "", err
	}

	filteredDiff := FilterDiff(rawDiff)
	if strings.TrimSpace(filteredDiff) == "" {
		return "", ErrEmptyDiff
	}

	return truncateDiffIfNeeded(filteredDiff), nil
}

func truncateDiffIfNeeded(diff string) string {
	runes := []rune(diff)
	if len(runes) > maxDiffSize {
		runes = runes[:maxDiffSize]
	}
	return string(runes) + "\n\n[WARNING: 代码变更量过大，已自动截断前 300KB 内容进行评审]"
}

// execGitDiff 执行底层的 git diff 命令
func execGitDiff(ctx context.Context, args ...string) (string, error) {
	fullArgs := append([]string{"diff"}, args...)
	cmd := exec.CommandContext(ctx, "git", fullArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return "", errors.New("系统未检测到 git 命令，请确保已安装 Git 并配置到系统 PATH 中")
		}
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg != "" {
			return "", fmt.Errorf("git diff 执行失败: %s (%w)", errMsg, err)
		}
		return "", fmt.Errorf("git diff 执行失败: %w", err)
	}
	return stdout.String(), nil
}
