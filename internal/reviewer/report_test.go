package reviewer

import "testing"

func TestParseReviewReport(t *testing.T) {
	raw := `#### 详细审查意见
- [BLOCKER] internal/user.go:45 - 存在 SQL 注入风险
- [WARNING] internal/user.go:82 - Redis 错误被忽略
- [WARNING] 缺少超时控制
#### 评审结论
- 结论: [REJECT]
- 评分: 58 / 100
`

	report := ParseReviewReport(raw)
	if report.Verdict != "REJECT" || report.Score != 58 {
		t.Fatalf("结论解析错误: %+v", report)
	}
	if len(report.Findings) != 3 || report.Findings[0].Location != "internal/user.go:45" {
		t.Fatalf("问题解析错误: %+v", report.Findings)
	}
	if report.Findings[2].Location != "" || report.Findings[2].Message != "缺少超时控制" {
		t.Fatalf("无位置信息的问题解析错误: %+v", report.Findings[2])
	}
}
