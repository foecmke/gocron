package embed

import (
	"embed"
	"io/fs"
)

// docsFiles 嵌入 VitePress 文档的 markdown 正文，供 AI 的 search_docs 工具离线检索。
// 只嵌入正文 .md（不含主题/资源），随二进制分发、无需联网。
//
//go:embed docs/zh/guide/*.md docs/en/guide/*.md docs/zh/index.md docs/en/index.md
var docsFiles embed.FS

// DocsFS 返回内嵌的文档文件系统（根为仓库根，路径形如 docs/zh/guide/xxx.md）。
func DocsFS() fs.FS {
	return docsFiles
}
