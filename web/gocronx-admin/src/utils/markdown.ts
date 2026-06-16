/**
 * 极简、零依赖、XSS 安全的 Markdown 渲染器（用于 AI 聊天面板渲染模型输出）。
 *
 * 安全策略：先对整段输入做 HTML 转义，使原文里的任何 `<...>` 变为惰性文本；
 * 之后只注入我们自己白名单内的标签（h1-h6 / strong / em / code / pre / ul / ol /
 * li / table / a / blockquote / hr / p / br）。链接仅放行 http/https/mailto。
 * 因此无需额外 sanitizer 也不会产生 XSS。
 *
 * 支持子集：标题、粗体、斜体、行内代码、围栏代码块、有序/无序列表、GFM 表格、
 * 链接、引用块、分隔线、段落。刻意不追求完整 Markdown 规范。
 */

function escapeHtml(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

// 行内渲染：先按 `code` 切分（捕获组 split 把代码段放在奇数下标），代码段不再做其它处理，
// 其余段落处理粗体/斜体/链接。这样天然隔离代码内容，无需哨兵、不会误伤正文里的数字。
function renderInline(escaped: string): string {
  return escaped
    .split(/`([^`]+)`/g)
    .map((part, idx) => {
      if (idx % 2 === 1) {
        return `<code>${part}</code>`
      }
      let s = part.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
      s = s.replace(/(^|[^*])\*([^*\n]+)\*/g, '$1<em>$2</em>')
      s = s.replace(/(^|[^_])_([^_\n]+)_/g, '$1<em>$2</em>')
      s = s.replace(/\[([^\]]+)\]\(([^)\s]+)\)/g, (_m, label: string, url: string) =>
        /^(https?:|mailto:)/i.test(url)
          ? `<a href="${url}" target="_blank" rel="noopener noreferrer">${label}</a>`
          : label
      )
      return s
    })
    .join('')
}

function isTableSeparator(line: string): boolean {
  return /^\s*\|?\s*:?-{1,}:?\s*(\|\s*:?-{1,}:?\s*)+\|?\s*$/.test(line)
}

function splitRow(line: string): string[] {
  return line
    .trim()
    .replace(/^\||\|$/g, '')
    .split('|')
    .map((c) => c.trim())
}

/** 把 Markdown 源码渲染为安全 HTML 字符串。 */
export function renderMarkdown(src: string): string {
  const escaped = escapeHtml(src.replace(/\r\n/g, '\n'))
  const lines = escaped.split('\n')
  const out: string[] = []
  let i = 0
  let para: string[] = []

  const flushPara = (): void => {
    if (para.length) {
      out.push(`<p>${para.map(renderInline).join('<br>')}</p>`)
      para = []
    }
  }

  while (i < lines.length) {
    const line = lines[i]

    // 围栏代码块 ```lang ... ```
    if (/^\s*```/.test(line)) {
      flushPara()
      const code: string[] = []
      i++
      while (i < lines.length && !/^\s*```/.test(lines[i])) {
        code.push(lines[i])
        i++
      }
      i++ // 跳过结束的 ```
      out.push(`<pre><code>${code.join('\n')}</code></pre>`)
      continue
    }

    // 表格：当前行含 | 且下一行是分隔行
    if (line.includes('|') && i + 1 < lines.length && isTableSeparator(lines[i + 1])) {
      flushPara()
      const header = splitRow(line)
      i += 2 // 跳过表头与分隔行
      const rows: string[][] = []
      while (i < lines.length && lines[i].includes('|') && lines[i].trim() !== '') {
        rows.push(splitRow(lines[i]))
        i++
      }
      const head = header.map((c) => `<th>${renderInline(c)}</th>`).join('')
      const body = rows
        .map((r) => `<tr>${r.map((c) => `<td>${renderInline(c)}</td>`).join('')}</tr>`)
        .join('')
      out.push(`<table><thead><tr>${head}</tr></thead><tbody>${body}</tbody></table>`)
      continue
    }

    // 标题 # .. ######
    const heading = /^(#{1,6})\s+(.*)$/.exec(line)
    if (heading) {
      flushPara()
      const level = heading[1].length
      out.push(`<h${level}>${renderInline(heading[2])}</h${level}>`)
      i++
      continue
    }

    // 分隔线
    if (/^\s*(-{3,}|\*{3,}|_{3,})\s*$/.test(line)) {
      flushPara()
      out.push('<hr>')
      i++
      continue
    }

    // 引用块
    if (/^\s*>\s?/.test(line)) {
      flushPara()
      const quote: string[] = []
      while (i < lines.length && /^\s*>\s?/.test(lines[i])) {
        quote.push(lines[i].replace(/^\s*>\s?/, ''))
        i++
      }
      out.push(`<blockquote>${quote.map(renderInline).join('<br>')}</blockquote>`)
      continue
    }

    // 无序列表
    if (/^\s*[-*+]\s+/.test(line)) {
      flushPara()
      const items: string[] = []
      while (i < lines.length && /^\s*[-*+]\s+/.test(lines[i])) {
        items.push(`<li>${renderInline(lines[i].replace(/^\s*[-*+]\s+/, ''))}</li>`)
        i++
      }
      out.push(`<ul>${items.join('')}</ul>`)
      continue
    }

    // 有序列表
    if (/^\s*\d+\.\s+/.test(line)) {
      flushPara()
      const items: string[] = []
      while (i < lines.length && /^\s*\d+\.\s+/.test(lines[i])) {
        items.push(`<li>${renderInline(lines[i].replace(/^\s*\d+\.\s+/, ''))}</li>`)
        i++
      }
      out.push(`<ol>${items.join('')}</ol>`)
      continue
    }

    // 空行 → 段落分隔
    if (line.trim() === '') {
      flushPara()
      i++
      continue
    }

    // 普通文本行：累积成段落
    para.push(line)
    i++
  }

  flushPara()
  return out.join('')
}
