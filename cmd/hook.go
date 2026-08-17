package cmd

import (
	"git-air/internal/git"
	"git-air/internal/ui"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const preCommitScript = `#!/bin/sh
# git-air pre-commit hook
echo "[git-air] 正在执行 Commit 前代码自动化评审..."
git-air --cached
# 如果评审发现阻断性问题，可通过退出码阻断 commit
`

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "管理 Git Pre-commit 自动审查钩子",
}

var hookInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "在当前仓库一键安装 pre-commit 钩子",
	Run: func(cmd *cobra.Command, args []string) {
		repoRoot, err := git.GetRepoRoot(cmd.Context())
		if err != nil {
			ui.PrintError("请在 Git 仓库根目录下执行此命令")
			return
		}
		hookPath := filepath.Join(repoRoot, ".git", "hooks", "pre-commit")
		if err := os.WriteFile(hookPath, []byte(preCommitScript), 0755); err != nil {
			ui.PrintError("写入 hook 文件失败: %v", err)
			return
		}
		ui.PrintSuccess("已成功安装 pre-commit 钩子！每次执行 git commit 时将自动触发评审。")
	},
}

var hookUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "移除当前仓库的 pre-commit 钩子",
	Run: func(cmd *cobra.Command, args []string) {
		repoRoot, err := git.GetRepoRoot(cmd.Context())
		if err != nil {
			ui.PrintError("请在 Git 仓库根目录下执行此命令")
			return
		}
		hookPath := filepath.Join(repoRoot, ".git", "hooks", "pre-commit")
		_ = os.Remove(hookPath)
		ui.PrintSuccess("已成功卸载 pre-commit 钩子！")
	},
}

func init() {
	hookCmd.AddCommand(hookInstallCmd)
	hookCmd.AddCommand(hookUninstallCmd)
}
