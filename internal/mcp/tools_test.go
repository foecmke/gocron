package mcp

import (
	"testing"
	"time"
)

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
