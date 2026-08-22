package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/unclesam-ly/git-air/internal/ui"
	"gopkg.in/yaml.v3"
)

// AppConfig 全局配置结构体
type AppConfig struct {
	Provider     string  `yaml:"provider"`                // gemini, deepseek, ollama, openai, custom
	APIKey       string  `yaml:"api_key"`                 // API Key
	BaseURL      string  `yaml:"base_url"`                // 自定义请求地址
	Model        string  `yaml:"model"`                   // 模型名称
	CustomPrompt string  `yaml:"custom_prompt"`           // 用户自定义的 Base System Prompt
	PriceInput   float64 `yaml:"price_input,omitempty"`   // 自定义输入单价 (每 1M tokens 美元，可选)
	PriceOutput  float64 `yaml:"price_output,omitempty"`  // 自定义输出单价 (每 1M tokens 美元，可选)
	CommitLang   string  `yaml:"commit_lang,omitempty"`   // Commit Message 生成语言 (auto, zh, en, ja, ko 等)
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".git-air")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

// LoadConfig 智能读取配置：
// 1. 先读取用户全局配置 ~/.git-air/config.yaml (保证 APIKey 等全局生效)
// 2. 若当前项目存在 .git-air.yaml 或 git-air.yaml，则进行项目级增量覆盖
// 3. 环境变量最高优先级覆盖
func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{
		Provider: "gemini",
		Model:    "gemini-3.7-flash",
	}

	// 1. 首先加载全局配置 ~/.git-air/config.yaml
	configPath, err := getConfigPath()
	if err != nil {
		// 隐式跳过时，可考虑记录日志或在 verbose 模式下输出
	} else if data, err := os.ReadFile(configPath); err == nil && len(data) > 0 {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("解析全局配置文件 [%s] 失败: %w", configPath, err)
		}
	}

	// 2. 检查当前项目是否有专属于 git-air 的配置文件 (.git-air.yaml / git-air.yaml) 进行增量覆盖
	localCandidates := []string{".git-air.yaml", ".git-air.yml", "git-air.yaml"}
	for _, name := range localCandidates {
		if data, err := os.ReadFile(name); err == nil && len(data) > 0 {
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, fmt.Errorf("解析项目配置文件 [%s] 失败: %w", name, err)
			}
			break
		}
	}

	// 3. 环境变量最高优先级覆盖
	if envKey := os.Getenv("GIT_AIR_API_KEY"); envKey != "" {
		cfg.APIKey = envKey
	}
	if envProvider := os.Getenv("GIT_AIR_PROVIDER"); envProvider != "" {
		cfg.Provider = envProvider
	}
	if envModel := os.Getenv("GIT_AIR_MODEL"); envModel != "" {
		cfg.Model = envModel
	}
	if envBaseURL := os.Getenv("GIT_AIR_BASE_URL"); envBaseURL != "" {
		cfg.BaseURL = envBaseURL
	}
	if envLang := os.Getenv("GIT_AIR_COMMIT_LANG"); envLang != "" {
		cfg.CommitLang = envLang
	}

	return cfg, nil
}

// SaveConfig 保存配置到 ~/.git-air/config.yaml
func SaveConfig(cfg *AppConfig) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0600)
}

// configCmd 配置子命令
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "查看或修改 git-air 全局配置 (~/.git-air/config.yaml)",
}
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "设置全局配置项",
	Example: `  git air config set --provider gemini --key YOUR_KEY
  git air config set --provider deepseek --key YOUR_KEY --model deepseek-chat
  git air config set --commit-lang ja
  git air config set --provider ollama --model qwen2.5-coder`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := LoadConfig()
		if err != nil || cfg == nil {
			cfg = &AppConfig{
				Provider: "gemini",
				Model:    "gemini-3.7-flash",
			}
		}
		if p, _ := cmd.Flags().GetString("provider"); p != "" {
			cfg.Provider = p
		}
		if k, _ := cmd.Flags().GetString("key"); k != "" {
			cfg.APIKey = k
		}
		if m, _ := cmd.Flags().GetString("model"); m != "" {
			cfg.Model = m
		}
		if u, _ := cmd.Flags().GetString("base-url"); u != "" {
			cfg.BaseURL = u
		}
		if pr, _ := cmd.Flags().GetString("prompt"); pr != "" {
			cfg.CustomPrompt = pr
		}
		if pi, _ := cmd.Flags().GetFloat64("price-input"); pi > 0 {
			cfg.PriceInput = pi
		}
		if po, _ := cmd.Flags().GetFloat64("price-output"); po > 0 {
			cfg.PriceOutput = po
		}
		if cl, _ := cmd.Flags().GetString("commit-lang"); cl != "" {
			cfg.CommitLang = cl
		}
		if err := SaveConfig(cfg); err != nil {
			ui.PrintError("保存配置失败: %v", err)
			return
		}
		ui.PrintSuccess("配置更新成功！")
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "查看当前所有配置",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := LoadConfig()
		if err != nil || cfg == nil {
			cfg = &AppConfig{
				Provider: "gemini",
				Model:    "gemini-3.7-flash",
			}
		}
		path, _ := getConfigPath()
		fmt.Printf("配置文件路径: %s\n\n", path)
		fmt.Printf("Provider:      %s\n", cfg.Provider)
		fmt.Printf("Model:         %s\n", cfg.Model)
		fmt.Printf("BaseURL:       %s\n", cfg.BaseURL)
		maskedKey := "未配置"
		if len(cfg.APIKey) > 8 {
			maskedKey = cfg.APIKey[:4] + "..." + cfg.APIKey[len(cfg.APIKey)-4:]
		}
		fmt.Printf("APIKey:        %s\n", maskedKey)
		hasPrompt := "否 (使用默认内置架构师 Prompt)"
		if cfg.CustomPrompt != "" {
			hasPrompt = "是 (已自定义)"
		}
		fmt.Printf("CustomPrompt:  %s\n", hasPrompt)
		commitLang := "auto (自动识别系统语言)"
		if cfg.CommitLang != "" {
			commitLang = cfg.CommitLang
		}
		fmt.Printf("CommitLang:    %s\n", commitLang)
		if cfg.PriceInput > 0 || cfg.PriceOutput > 0 {
			fmt.Printf("CustomPrice:   输入 $%.4f/1M, 输出 $%.4f/1M (已自定义覆盖)\n", cfg.PriceInput, cfg.PriceOutput)
		} else {
			fmt.Printf("CustomPrice:   默认 (跟随官方最新价格表)\n")
		}
	},
}

func init() {
	configSetCmd.Flags().StringP("provider", "p", "", "模型提供商 (gemini, claude, grok, deepseek, qwen, zhipu, moonshot, siliconflow, doubao, minimax, yi, groq, openrouter, ollama, openai)")
	configSetCmd.Flags().StringP("key", "k", "", "API 密钥")
	configSetCmd.Flags().StringP("model", "m", "", "模型名称")
	configSetCmd.Flags().StringP("base-url", "u", "", "自定义 API 地址")
	configSetCmd.Flags().String("prompt", "", "自定义 Base System Prompt")
	configSetCmd.Flags().Float64("price-input", 0, "自定义每 1M 输入 Token 美元价格")
	configSetCmd.Flags().Float64("price-output", 0, "自定义每 1M 输出 Token 美元价格")
	configSetCmd.Flags().String("commit-lang", "", "Commit Message 语言 (auto, zh, en, ja, ko 等)")

	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
}
