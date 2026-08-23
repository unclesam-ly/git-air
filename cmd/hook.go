package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/unclesam-ly/git-air/internal/git"
	"github.com/unclesam-ly/git-air/internal/ui"
)

const (
	preCommitMarker         = "# git-air-managed-pre-commit-v1"
	preCommitBackupSuffix   = ".git-air-backup"
	preCommitBackupFileName = "pre-commit.git-air-backup"
)

func preCommitScript(bin ...string) string {
	binPath := "git-air"
	if len(bin) > 0 && bin[0] != "" {
		binPath = bin[0]
	}
	return fmt.Sprintf(`#!/bin/sh
%s
echo "[git-air] 正在执行 Commit 前代码自动化评审..."
"%s" --staged
git_air_status=$?
if [ "$git_air_status" -ne 0 ]; then
  exit "$git_air_status"
fi

hook_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
old_hook="$hook_dir/%s"
if [ -x "$old_hook" ]; then
	exec "$old_hook" "$@"
fi
exit 0
`, preCommitMarker, binPath, preCommitBackupFileName)
}

func isGitAirHook(data []byte) bool {
	content := string(data)
	return strings.Contains(content, preCommitMarker) ||
		strings.Contains(content, "# git-air pre-commit hook") ||
		strings.Contains(content, "git-air --staged")
}

func installPreCommitHook(hookPath string) error {
	backupPath := hookPath + preCommitBackupSuffix
	hookData, hookReadErr := os.ReadFile(hookPath)
	hookExists := hookReadErr == nil
	if hookReadErr != nil && !errors.Is(hookReadErr, os.ErrNotExist) {
		return hookReadErr
	}
	// 幂等安装：已经是 git-air 包装 Hook 时直接返回，绝不把自己再次备份。
	if hookExists && isGitAirHook(hookData) {
		return nil
	}

	if _, err := os.Stat(backupPath); err == nil {
		return fmt.Errorf("检测到已有 git-air 备份文件 %s，请先执行 hook uninstall", backupPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if hookExists {
		// 保留原 Hook：先移动为备份，再安装包装 Hook。写入失败时尝试恢复。
		if err := os.Rename(hookPath, backupPath); err != nil {
			return fmt.Errorf("备份已有 pre-commit Hook 失败: %w", err)
		}
	}

	binPath := "git-air"
	if exe, err := os.Executable(); err == nil && exe != "" {
		binPath = exe
	}

	if err := os.WriteFile(hookPath, []byte(preCommitScript(binPath)), 0755); err != nil {
		if hookExists {
			_ = os.Remove(hookPath)
			_ = os.Rename(backupPath, hookPath)
		}
		return fmt.Errorf("写入 git-air pre-commit Hook 失败: %w", err)
	}
	return nil
}

func uninstallPreCommitHook(hookPath string) error {
	data, err := os.ReadFile(hookPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("当前仓库没有 git-air pre-commit Hook")
		}
		return err
	}
	if !isGitAirHook(data) {
		return fmt.Errorf("当前 pre-commit Hook 不是由 git-air 管理，已保留原文件")
	}

	backupPath := hookPath + preCommitBackupSuffix
	if _, err := os.Stat(backupPath); err == nil {
		return restoreOriginalHook(hookPath, backupPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return os.Remove(hookPath)
}

// restoreOriginalHook 恢复备份 Hook。
// POSIX 系统优先使用 Rename 的原子替换；Windows 目标文件存在时 Rename 可能失败，
// 此时先把包装 Hook 移到临时路径，恢复失败则尝试回滚，避免先删除导致数据丢失。
func restoreOriginalHook(hookPath, backupPath string) error {
	if err := os.Rename(backupPath, hookPath); err == nil {
		return nil
	}

	tempPath := hookPath + ".git-air-restore-temp"
	if _, err := os.Stat(tempPath); err == nil {
		return fmt.Errorf("恢复原 pre-commit Hook 失败：临时文件已存在 %s", tempPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.Rename(hookPath, tempPath); err != nil {
		return fmt.Errorf("暂存 git-air Hook 失败: %w", err)
	}
	if err := os.Rename(backupPath, hookPath); err != nil {
		_ = os.Rename(tempPath, hookPath)
		return fmt.Errorf("恢复原 pre-commit Hook 失败: %w", err)
	}
	if err := os.Remove(tempPath); err != nil {
		return fmt.Errorf("清理临时 Hook 文件失败: %w", err)
	}
	return nil
}

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "管理 Git Pre-commit 自动审查钩子",
}

var hookInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "在当前仓库安装 pre-commit 钩子并保留已有 Hook",
	Run: func(cmd *cobra.Command, args []string) {
		repoRoot, err := git.GetRepoRoot(cmd.Context())
		if err != nil {
			ui.PrintError("请在 Git 仓库根目录下执行此命令")
			return
		}
		hooksDir := filepath.Join(repoRoot, ".git", "hooks")
		if err := os.MkdirAll(hooksDir, 0755); err != nil {
			ui.PrintError("创建 hooks 目录失败: %v", err)
			return
		}
		hookPath := filepath.Join(hooksDir, "pre-commit")
		if err := installPreCommitHook(hookPath); err != nil {
			ui.PrintError("安装 Hook 失败: %v", err)
			return
		}
		ui.PrintSuccess("已安装 git-air Hook；已有 Hook 已保存在 pre-commit.git-air-backup，提交通过后会继续执行。")
	},
}

var hookUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "移除 git-air Hook 并恢复已有 Hook",
	Run: func(cmd *cobra.Command, args []string) {
		repoRoot, err := git.GetRepoRoot(cmd.Context())
		if err != nil {
			ui.PrintError("请在 Git 仓库根目录下执行此命令")
			return
		}
		hookPath := filepath.Join(repoRoot, ".git", "hooks", "pre-commit")
		if err := uninstallPreCommitHook(hookPath); err != nil {
			ui.PrintError("卸载 Hook 失败: %v", err)
			return
		}
		ui.PrintSuccess("已卸载 git-air Hook；如原来存在 Hook，已恢复原文件。")
	},
}

func init() {
	hookCmd.AddCommand(hookInstallCmd)
	hookCmd.AddCommand(hookUninstallCmd)
}
