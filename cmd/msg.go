package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/unclesam-ly/git-air/internal/git"
	"github.com/unclesam-ly/git-air/internal/llm"
	"github.com/unclesam-ly/git-air/internal/reviewer"
	"github.com/unclesam-ly/git-air/internal/ui"
)

var msgCmd = &cobra.Command{
	Use:     "msg [git diff 参数，如 HEAD~1, main..feature, 文件名等]",
	Aliases: []string{"commit-msg", "commit"},
	Short:   "根据代码变更智能生成规范的 Conventional Commits 提交信息",
	Example: `  git air msg               # 根据当前暂存区代码自动生成
  git air msg -l ja          # 指定生成日语提交说明
  git air msg -l en          # 指定生成英文提交说明
  git air msg -c             # 生成后直接执行 git commit`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		// 1. 检查并获取 Diff (默认优先获取暂存区)
		diff, err := git.GetDiff(ctx, args...)
		if err != nil {
			ui.PrintError("%v", err)
			return
		}

		// 2. 加载配置
		cfg, err := LoadConfig()
		if err != nil {
			ui.PrintError("加载配置失败: %v", err)
			return
		}

		if flagKey, _ := cmd.Flags().GetString("key"); flagKey != "" {
			cfg.APIKey = flagKey
		}
		if flagModel, _ := cmd.Flags().GetString("model"); flagModel != "" {
			cfg.Model = flagModel
		}
		if flagProvider, _ := cmd.Flags().GetString("provider"); flagProvider != "" {
			cfg.Provider = flagProvider
		}

		// 3. 初始化大模型客户端
		llmClient, err := llm.NewClient(llm.Config{
			Provider:    cfg.Provider,
			APIKey:      cfg.APIKey,
			BaseURL:     cfg.BaseURL,
			Model:       cfg.Model,
			PriceInput:  cfg.PriceInput,
			PriceOutput: cfg.PriceOutput,
		})
		if err != nil {
			ui.PrintError("初始化模型失败: %v", err)
			return
		}

		// 4. 确定目标语言 (命令行 Flag > 配置文件 > 自动探测系统语言)
		lang, _ := cmd.Flags().GetString("lang")
		if lang == "" {
			lang = cfg.CommitLang
		}

		repoRoot, _ := git.GetRepoRoot(ctx)
		rev := reviewer.NewReviewer(llmClient, repoRoot, "")

		targetLang := reviewer.DetectLanguage(lang)
		langDesc := "English"
		switch targetLang {
		case "zh":
			langDesc = "简体中文"
		case "ja":
			langDesc = "日本語"
		case "ko":
			langDesc = "한국어"
		}

		fmt.Printf("\033[1;36m[git-air]\033[0m 正在生成 Conventional Commit Message... \033[0;90m(Language: %s / Engine: %s)\033[0m\n", langDesc, cfg.Model)
		fmt.Println("\033[0;90m" + strings.Repeat("─", 65) + "\033[0m")

		commitMsg, err := rev.GenerateCommitMsg(ctx, diff, lang)
		if err != nil {
			ui.PrintError("生成 Commit Message 失败: %v", err)
			return
		}

		// 5. 打印生成的提交信息
		fmt.Printf("\033[1;32m%s\033[0m\n", commitMsg)
		fmt.Println("\033[0;90m" + strings.Repeat("─", 65) + "\033[0m")

		// 6. 是否自动执行提交
		autoCommit, _ := cmd.Flags().GetBool("commit")
		if autoCommit {
			executeGitCommit(commitMsg)
			return
		}

		// 交互式询问是否直接提交
		fmt.Print("? 是否直接以此信息执行 git commit？[Y/n]: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input == "" || input == "y" || input == "yes" {
			executeGitCommit(commitMsg)
		} else {
			fmt.Println("\033[0;90m已跳过提交。你可以复制上方内容手动提交。\033[0m")
		}
	},
}

func executeGitCommit(msg string) {
	cmd := exec.Command("git", "commit", "-m", msg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		ui.PrintError("执行 git commit 失败: %v", err)
		return
	}
	ui.PrintSuccess("代码提交成功！🎉")
}

func init() {
	msgCmd.Flags().StringP("lang", "l", "", "生成语言 (auto, zh, en, ja, ko)")
	msgCmd.Flags().BoolP("commit", "c", false, "生成后无需确认直接执行 git commit")
	msgCmd.Flags().StringP("key", "k", "", "临时指定 API Key")
	msgCmd.Flags().StringP("model", "m", "", "临时指定模型名称")
	msgCmd.Flags().StringP("provider", "p", "", "临时指定 Provider")

	rootCmd.AddCommand(msgCmd)
}
