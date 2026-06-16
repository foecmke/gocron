package mcp

import (
	"testing"
	"time"

	"github.com/gocronx-team/gocron/internal/models"
)

func TestDiagnoseTaskLog_GuardPaths(t *testing.T) {
	defer setupTestDb(t)()

	// 日志不存在 → 报错（不触达 LLM）
	if _, err := diagnoseTaskLog(diagnoseTaskLogInput{Id: 9999}); err == nil {
		t.Fatal("expected error for missing log")
	}

	// 日志存在但无执行输出 → 报错（不触达 LLM）
	log := &models.TaskLog{Id: 1, Name: "a", Result: "   "}
	if err := models.Db.Create(log).Error; err != nil {
		t.Fatalf("seed log: %v", err)
	}
	if _, err := diagnoseTaskLog(diagnoseTaskLogInput{Id: 1}); err == nil {
		t.Fatal("expected error when log has no result")
	}
}

func TestParseTimeArg(t *testing.T) {
	cases := []struct {
		name string
		in   string
		ok   bool
		want time.Time
	}{
		{"empty", "", false, time.Time{}},
		{"blank", "   ", false, time.Time{}},
		{"date only", "2026-06-15", true, time.Date(2026, 6, 15, 0, 0, 0, 0, time.Local)},
		{"datetime", "2026-06-15 20:30:00", true, time.Date(2026, 6, 15, 20, 30, 0, 0, time.Local)},
		{"rfc3339", "2026-06-15T20:30:00Z", true, time.Date(2026, 6, 15, 20, 30, 0, 0, time.UTC)},
		{"garbage", "not a time", false, time.Time{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := parseTimeArg(c.in)
			if ok != c.ok {
				t.Fatalf("ok = %v, want %v", ok, c.ok)
			}
			if ok && !got.Equal(c.want) {
				t.Fatalf("time = %v, want %v", got, c.want)
			}
		})
	}
}
