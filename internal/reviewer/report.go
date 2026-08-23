package reviewer

import (
	"regexp"
	"strconv"
	"strings"
)

type Finding struct {
	Severity string `json:"severity"`
	Location string `json:"location,omitempty"`
	Message  string `json:"message"`
}

type ReviewReport struct {
	Format   string    `json:"format"`
	Verdict  string    `json:"verdict,omitempty"`
	Score    int       `json:"score,omitempty"`
	Findings []Finding `json:"findings"`
	Raw      string    `json:"raw"`
}

var (
	verdictPattern = regexp.MustCompile(`(?mi)(?:结论|verdict)\s*[:：]\s*\[?(PASS|WARN|WARNING|REJECT|BLOCKER|CRITICAL)\]?`)
	scorePattern   = regexp.MustCompile(`(?mi)(?:评分|score)\s*[:：]\s*(\d{1,3})`)
	findingPattern = regexp.MustCompile(`(?m)^\s*-\s*\[(BLOCKER|CRITICAL|WARNING|WARN|SUGGESTION|INFO)\]\s+(.+?)\s*$`)
)

// ParseReviewReport 将模型的 Markdown Review 转成稳定的 JSON 结构。
// Raw 始终保留完整原文，解析失败时下游仍可自行处理。
func ParseReviewReport(raw string) ReviewReport {
	report := ReviewReport{
		Format:   "git-air-review/v1",
		Findings: []Finding{},
		Raw:      strings.TrimSpace(raw),
	}

	if match := verdictPattern.FindStringSubmatch(raw); len(match) > 1 {
		report.Verdict = strings.ToUpper(match[1])
	}
	if match := scorePattern.FindStringSubmatch(raw); len(match) > 1 {
		report.Score, _ = strconv.Atoi(match[1])
		if report.Score > 100 {
			report.Score = 100
		}
	}

	for _, match := range findingPattern.FindAllStringSubmatch(raw, -1) {
		location, message := splitFindingContent(match[2])
		report.Findings = append(report.Findings, Finding{
			Severity: strings.ToUpper(match[1]),
			Location: location,
			Message:  message,
		})
	}
	return report
}

func splitFindingContent(content string) (string, string) {
	content = strings.TrimSpace(content)
	parts := strings.SplitN(content, " - ", 2)
	if len(parts) == 1 {
		return "", content
	}

	location := strings.TrimSpace(parts[0])
	message := strings.TrimSpace(parts[1])
	// 没有行号/路径特征时，不要误把普通句子当成 location。
	if !strings.Contains(location, ":") {
		return "", content
	}
	return location, message
}
