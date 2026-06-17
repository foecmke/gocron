package docsearch

import (
	"strings"
	"testing"
)

func TestSearch_FindsRelevantEnglish(t *testing.T) {
	got := Search("MCP", 4)
	if len(got) == 0 {
		t.Fatal("expected results for 'MCP'")
	}
	found := false
	for _, c := range got {
		if strings.Contains(strings.ToLower(c.Content+c.Heading), "mcp") {
			found = true
		}
	}
	if !found {
		t.Fatalf("no returned chunk mentions MCP: %+v", got)
	}
}

func TestSearch_FindsRelevantChinese(t *testing.T) {
	// 中文按二元字组匹配
	got := Search("数据库", 4)
	if len(got) == 0 {
		t.Fatal("expected results for '数据库'")
	}
}

func TestSearch_NoMatch(t *testing.T) {
	if got := Search("zzzznotinanydocxyzzy", 4); len(got) != 0 {
		t.Fatalf("expected no results, got %d", len(got))
	}
}

func TestSearch_EmptyQuery(t *testing.T) {
	if got := Search("   ", 4); len(got) != 0 {
		t.Fatalf("expected no results for blank query, got %d", len(got))
	}
}

func TestSplitChunks_StripsFrontmatter(t *testing.T) {
	md := "---\nlayout: home\n---\n# Title\n\nbody text here\n\n## Section\nmore\n"
	chunks := splitChunks("x.md", md)
	for _, c := range chunks {
		if strings.Contains(c.Content, "layout: home") {
			t.Fatalf("frontmatter leaked into chunk: %q", c.Content)
		}
	}
	if len(chunks) == 0 {
		t.Fatal("expected chunks")
	}
}

func TestTokenize(t *testing.T) {
	terms := tokenize("how to 配置数据库")
	has := func(x string) bool {
		for _, t := range terms {
			if t == x {
				return true
			}
		}
		return false
	}
	if !has("how") || !has("配置") || !has("数据") {
		t.Fatalf("unexpected tokens: %v", terms)
	}
}
