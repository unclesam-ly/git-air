package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/unclesam-ly/git-air/internal/ui"
	"github.com/unclesam-ly/git-air/internal/updater"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade"},
	Short:   "一键升级 git-air 到最新版本",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("🔍 正在检查并拉取最新版本 (当前: %s)...\n", updater.CurrentVersion)

		updateCmd := exec.Command("go", "install", "github.com/unclesam-ly/git-air@latest")
		updateCmd.Stdout = os.Stdout
		updateCmd.Stderr = os.Stderr

		if err := updateCmd.Run(); err != nil {
			ui.PrintError("升级失败: %v (提示: 请确保已安装 Go 并配置环境变量)", err)
			return
		}

		ui.PrintSuccess("git-air 已成功升级至最新版本！🎉")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
