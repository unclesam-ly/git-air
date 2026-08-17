package ui

import (
	"fmt"
	"os"
	"strings"
)

// ANSI 终端颜色常量
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[1;31m"
	colorGreen  = "\033[1;32m"
	colorYellow = "\033[1;33m"
	colorBlue   = "\033[1;34m"
	colorCyan   = "\033[1;36m"
	colorGray   = "\033[0;90m"
	colorBold   = "\033[1m"
)

// PrintBanner 打印开场横幅与当前评审状态
func PrintBanner(provider, model string) {
	fmt.Printf("%s[git-air]%s 代码评审中... %s(Engine: %s / %s)%s\n", colorCyan, colorReset, colorGray, provider, model, colorReset)
	fmt.Println(colorGray + strings.Repeat("─", 65) + colorReset)
}

// PrintFooter 打印结束分界线
func PrintFooter() {
	fmt.Println()
	fmt.Println(colorGray + strings.Repeat("─", 65) + colorReset)
}

// PrintError 打印错误信息
func PrintError(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stderr, "%s[ERROR]%s %s\n", colorRed, colorReset, msg)
}

// PrintSuccess 打印成功信息
func PrintSuccess(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s[SUCCESS]%s %s\n", colorGreen, colorReset, msg)
}

// StreamPrinter 带有标签高亮的流式输出打印器
type StreamPrinter struct{}

func NewStreamPrinter() *StreamPrinter {
	return &StreamPrinter{}
}

// PrintChunk 逐字流式输出，遇到关键标签进行颜色增强
func (p *StreamPrinter) PrintChunk(chunk string) {
	// 对关键标签进行颜色替换
	highlighted := chunk
	highlighted = strings.ReplaceAll(highlighted, "[BLOCKER]", colorRed+"[BLOCKER]"+colorReset)
	highlighted = strings.ReplaceAll(highlighted, "[CRITICAL]", colorRed+"[CRITICAL]"+colorReset)
	highlighted = strings.ReplaceAll(highlighted, "[WARNING]", colorYellow+"[WARNING]"+colorReset)
	highlighted = strings.ReplaceAll(highlighted, "[WARN]", colorYellow+"[WARN]"+colorReset)
	highlighted = strings.ReplaceAll(highlighted, "[SUGGESTION]", colorBlue+"[SUGGESTION]"+colorReset)
	highlighted = strings.ReplaceAll(highlighted, "[PASS]", colorGreen+"[PASS]"+colorReset)
	highlighted = strings.ReplaceAll(highlighted, "[LGTM]", colorGreen+"[LGTM]"+colorReset)
	highlighted = strings.ReplaceAll(highlighted, "[REJECT]", colorRed+"[REJECT]"+colorReset)

	fmt.Print(highlighted)
}
