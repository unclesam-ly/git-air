package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/unclesam-ly/git-air/internal/ui"
	"github.com/unclesam-ly/git-air/internal/updater"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "查看当前 git-air 版本信息及更新检测",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("git-air %s (%s/%s, Go %s)\n",
			updater.CurrentVersion,
			runtime.GOOS,
			runtime.GOARCH,
			runtime.Version(),
		)

		checkFlag, _ := cmd.Flags().GetBool("check")
		if checkFlag {
			info := updater.CheckUpdate(true)
			if info != nil {
				ui.PrintUpdateBanner(info.CurrentVersion, info.LatestVersion, info.ReleaseURL)
			} else {
				ui.PrintSuccess("当前已是最新版本！")
			}
		}
	},
}

func init() {
	versionCmd.Flags().BoolP("check", "c", true, "向 GitHub 检查是否有可用新版本")
	rootCmd.AddCommand(versionCmd)
}
