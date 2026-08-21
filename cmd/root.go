package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/unclesam-ly/git-air/internal/git"
	"github.com/unclesam-ly/git-air/internal/llm"
	"github.com/unclesam-ly/git-air/internal/reviewer"
	"github.com/unclesam-ly/git-air/internal/ui"
)

var rootCmd = &cobra.Command{
	Use:   "git-air [git diff 参数，如 HEAD~1, main..feature, 文件名等]",
	Short: "git-air: 极轻量、极速的 Git 原生 AI 代码评审工具",
	Long: `git-air 是一个纯 Go 编写的 Git AI 评审利器。
支持一键评审暂存区、指定 Commit、分支差异或单文件，无缝兼容 Gemini、DeepSeek、Ollama 和 OpenAI。`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		// 1. 检查并获取 Diff
		if flagCached, _ := cmd.Flags().GetBool("cached"); flagCached {
			args = append([]string{"--cached"}, args...)
		} else if flagStaged, _ := cmd.Flags().GetBool("staged"); flagStaged {
			args = append([]string{"--cached"}, args...)
		}

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

		// 允许命令行临时覆盖模型/Key
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
			ui.PrintError("初始化模型失败: %v (提示: 请先执行 'git air config set --key YOUR_KEY')", err)
			return
		}

		// 4. 获取仓库根目录 (用于加载 .airules)
		repoRoot, _ := git.GetRepoRoot(ctx)

		// 5. 初始化评审引擎
		customPrompt, _ := cmd.Flags().GetString("prompt")
		if customPrompt == "" {
			customPrompt = cfg.CustomPrompt
		}
		rev := reviewer.NewReviewer(llmClient, repoRoot, customPrompt)

		// 6. 打印开场横幅并开始流式评审
		ui.PrintBanner(cfg.Provider, cfg.Model)
		printer := ui.NewStreamPrinter()
		var reviewOutput strings.Builder
		_, usage, err := rev.Execute(ctx, diff, true, func(chunk string) {
			reviewOutput.WriteString(chunk)
			printer.PrintChunk(chunk)
		})
		if err != nil {
			ui.PrintError("\n评审过程中发生错误: %v", err)
			os.Exit(1)
		}

		// 7. 打印 Token 消耗统计
		if usage != nil {
			ui.PrintUsage(usage.InputTokens, usage.OutputTokens, usage.EstimateCost())
		}

		ui.PrintFooter()

		// 8. 门禁阻断检测 (Pre-commit 拦截逻辑)
		outText := reviewOutput.String()
		hasBlocker := strings.Contains(outText, "[BLOCKER]") || strings.Contains(outText, "[CRITICAL]")
		hasReject := strings.Contains(outText, "[REJECT]")
		hasWarning := strings.Contains(outText, "[WARNING]") || strings.Contains(outText, "[WARN]")

		strict, _ := cmd.Flags().GetBool("strict")
		noBlock, _ := cmd.Flags().GetBool("no-block")

		if !noBlock {
			if hasBlocker || hasReject {
				ui.PrintError("[git-air] 审查未通过：检测到 [BLOCKER] 严重缺陷或 [REJECT] 结论，已拦截提交！")
				ui.PrintError("修复上述问题后再提交；如需临时跳过可执行 'git commit --no-verify' 或加上 '--no-block'")
				os.Exit(1)
			} else if strict && hasWarning {
				ui.PrintError("[git-air] 严格模式（--strict）生效：检测到 [WARNING] 警告，已拦截提交！")
				ui.PrintError("请优化上述警告或去掉 --strict 标志。")
				os.Exit(1)
			}
		}
	},
}

func Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(hookCmd)

	rootCmd.Flags().StringP("key", "k", "", "临时指定 API Key")
	rootCmd.Flags().StringP("model", "m", "", "临时指定模型名称")
	rootCmd.Flags().StringP("provider", "p", "", "临时指定 Provider (gemini, deepseek, ollama)")
	rootCmd.Flags().String("prompt", "", "临时指定 System Prompt")
	rootCmd.Flags().Bool("cached", false, "显式审查暂存区（git add）代码")
	rootCmd.Flags().Bool("staged", false, "显式审查暂存区（git add）代码")
	rootCmd.Flags().Bool("strict", false, "严格模式：遇到 WARNING 也会阻断退出")
	rootCmd.Flags().Bool("no-block", false, "非阻断模式：即使发现 BLOCKER 也不以非零退出码阻断")
}
