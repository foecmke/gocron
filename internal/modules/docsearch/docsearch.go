// Package docsearch 对内嵌的 gocron 文档做离线关键词检索，供 AI 的 search_docs 工具使用。
// 纯本地、无依赖、随二进制分发：按标题把文档切块，按关键词重叠度打分返回最相关的若干段。
// 中英混合:英文按单词、中文按相邻二元字组(bigram)匹配。
package docsearch

import (
	"io/fs"
	"sort"
	"strings"
	"sync"
	"unicode"

	gocronembed "github.com/gocronx-team/gocron"
)

const (
	maxChunkRunes = 1200 // 单块返回上限，控制 token
	defaultTopN   = 4
)

// Chunk 是一段可检索的文档片段。
type Chunk struct {
	Source  string `json:"source"`  // 文件路径，如 docs/zh/guide/configuration.md
	Heading string `json:"heading"` // 所属标题
	Content string `json:"content"` // 正文（已截断）
}

var (
	once   sync.Once
	chunks []Chunk
)

// load 解析并切块所有内嵌文档（只做一次）。
func load() {
	docs := gocronembed.DocsFS()
	_ = fs.WalkDir(docs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		data, rerr := fs.ReadFile(docs, path)
		if rerr != nil {
			return nil
		}
		chunks = append(chunks, splitChunks(path, string(data))...)
		return nil
	})
}

// splitChunks 去掉 frontmatter，按 markdown 标题把一篇文档切成若干块。
func splitChunks(path, content string) []Chunk {
	content = stripFrontmatter(content)
	lines := strings.Split(content, "\n")
	var out []Chunk
	heading := ""
	var body []string
	flush := func() {
		text := strings.TrimSpace(strings.Join(body, "\n"))
		if text == "" && heading == "" {
			return
		}
		if r := []rune(text); len(r) > maxChunkRunes {
			text = string(r[:maxChunkRunes])
		}
		out = append(out, Chunk{Source: path, Heading: heading, Content: text})
	}
	for _, line := range lines {
		if h := strings.TrimLeft(line, "#"); h != line && strings.HasPrefix(strings.TrimSpace(line), "#") {
			flush()
			heading = strings.TrimSpace(h)
			body = nil
			continue
		}
		body = append(body, line)
	}
	flush()
	return out
}

func stripFrontmatter(s string) string {
	s = strings.TrimLeft(s, "\ufeff \n\r\t")
	if strings.HasPrefix(s, "---") {
		if end := strings.Index(s[3:], "\n---"); end >= 0 {
			rest := s[3+end+4:]
			if i := strings.IndexByte(rest, '\n'); i >= 0 {
				return rest[i+1:]
			}
			return ""
		}
	}
	return s
}

func isCJK(r rune) bool {
	return unicode.Is(unicode.Han, r)
}

// tokenize 把查询拆成检索词:英文按单词(≥2 字符)、中文按相邻二元字组。
func tokenize(s string) []string {
	s = strings.ToLower(s)
	seen := map[string]bool{}
	var terms []string
	add := func(t string) {
		if t != "" && !seen[t] {
			seen[t] = true
			terms = append(terms, t)
		}
	}
	// ASCII 单词
	for _, w := range strings.FieldsFunc(s, func(r rune) bool {
		return r > unicode.MaxASCII || (!unicode.IsLetter(r) && !unicode.IsDigit(r))
	}) {
		if len(w) >= 2 {
			add(w)
		}
	}
	// 中文二元字组
	runes := []rune(s)
	for i := 0; i+1 < len(runes); i++ {
		if isCJK(runes[i]) && isCJK(runes[i+1]) {
			add(string(runes[i : i+2]))
		}
	}
	return terms
}

// Search 返回与 query 最相关的至多 topN 段文档。topN<=0 时取默认值。
func Search(query string, topN int) []Chunk {
	once.Do(load)
	if topN <= 0 {
		topN = defaultTopN
	}
	terms := tokenize(query)
	if len(terms) == 0 {
		return nil
	}
	type scored struct {
		c     Chunk
		score int
	}
	var ranked []scored
	for _, c := range chunks {
		body := strings.ToLower(c.Content)
		head := strings.ToLower(c.Heading)
		score := 0
		for _, t := range terms {
			score += strings.Count(body, t)
			score += 3 * strings.Count(head, t) // 命中标题权重更高
		}
		if score > 0 {
			ranked = append(ranked, scored{c, score})
		}
	}
	sort.SliceStable(ranked, func(i, j int) bool { return ranked[i].score > ranked[j].score })
	if len(ranked) > topN {
		ranked = ranked[:topN]
	}
	out := make([]Chunk, 0, len(ranked))
	for _, r := range ranked {
		out = append(out, r.c)
	}
	return out
}
