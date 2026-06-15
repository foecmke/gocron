/**
 * 复制文本到剪贴板。
 *
 * navigator.clipboard 仅在安全上下文(HTTPS 或 localhost)可用；
 * 通过局域网 IP 明文 HTTP 访问时它为 undefined，会导致直接调用抛错，
 * 因此这里在不可用或失败时自动降级到 document.execCommand。
 *
 * @returns 是否复制成功
 */
export async function copyToClipboard(text: string): Promise<boolean> {
  if (navigator.clipboard?.writeText) {
    try {
      await navigator.clipboard.writeText(text)
      return true
    } catch {
      // 继续走 execCommand 降级方案
    }
  }

  const el = document.createElement('textarea')
  el.value = text
  el.style.position = 'fixed'
  el.style.opacity = '0'
  document.body.appendChild(el)
  el.select()
  let ok = false
  try {
    ok = document.execCommand('copy')
  } catch {
    ok = false
  }
  document.body.removeChild(el)
  return ok
}
