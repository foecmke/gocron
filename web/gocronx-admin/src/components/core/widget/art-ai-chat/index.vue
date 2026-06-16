<!-- AI 运维助手聊天面板 -->
<template>
  <ArtIconButton icon="ri:robot-2-line" :title="t('aiChat.title')" @click="openDrawer" />

  <ElDrawer
    v-model="visible"
    :title="t('aiChat.title')"
    direction="rtl"
    size="420px"
    class="ai-chat-drawer"
    @close="cancelStream"
  >
    <template #header>
      <div class="ai-chat-header">
        <span class="ai-chat-title">{{ t('aiChat.title') }}</span>
        <ElButton link :disabled="messages.length === 0 || loading" @click="clearConversation">
          {{ t('aiChat.clear') }}
        </ElButton>
      </div>
    </template>

    <div class="ai-chat-body">
      <div ref="listRef" class="ai-chat-list">
        <p v-if="messages.length === 0" class="ai-chat-empty">{{ t('aiChat.empty') }}</p>

        <div
          v-for="(msg, index) in messages"
          :key="index"
          class="ai-chat-row"
          :class="msg.role === 'user' ? 'is-user' : 'is-assistant'"
        >
          <div class="ai-chat-bubble">
            <!-- 助手消息渲染 Markdown（renderMarkdown 已做 HTML 转义 + 白名单，XSS 安全）；用户消息保持纯文本 -->
            <div
              v-if="msg.content && msg.role === 'assistant'"
              class="ai-chat-md"
              v-html="renderMarkdown(msg.content)"
            ></div>
            <template v-else-if="msg.content">{{ msg.content }}</template>
            <span v-else-if="msg.role === 'assistant' && loading" class="ai-chat-thinking">
              {{ t('aiChat.thinking') }}
            </span>
          </div>
          <ElButton
            v-if="msg.content"
            link
            size="small"
            class="ai-chat-copy"
            @click="copyMessage(msg.content)"
          >
            {{ t('aiChat.copy') }}
          </ElButton>
          <div v-if="msg.role === 'assistant' && toolsByIndex[index]?.length" class="ai-chat-tools">
            <ElTag
              v-for="tool in toolsByIndex[index]"
              :key="tool.id"
              size="small"
              :type="
                tool.status === 'failed' ? 'danger' : tool.status === 'done' ? 'success' : 'info'
              "
            >
              {{ toolLabel(tool) }}
            </ElTag>
          </div>
        </div>
      </div>

      <div class="ai-chat-input">
        <ElInput
          v-model="draft"
          type="textarea"
          :rows="3"
          resize="none"
          :placeholder="t('aiChat.placeholder')"
          @keydown.enter="onEnter"
        />
        <ElButton type="primary" :loading="loading" :disabled="!draft.trim()" @click="send">
          {{ t('aiChat.send') }}
        </ElButton>
      </div>
    </div>
  </ElDrawer>
</template>

<script setup lang="ts">
  import { useI18n } from 'vue-i18n'
  import { ElButton, ElDrawer, ElInput, ElMessage, ElTag } from 'element-plus'
  import { streamAiChat, type AiChatMessage } from '@/api/ai'
  import { copyToClipboard } from '@/utils/clipboard'
  import { renderMarkdown } from '@/utils/markdown'

  defineOptions({ name: 'ArtAiChat' })

  const { t } = useI18n()

  // 单条消息内展示的工具调用 chip：running → done/failed。
  type ToolChip = { id: string; name: string; status: 'running' | 'done' | 'failed' }

  const visible = ref(false)
  const loading = ref(false)
  const draft = ref('')
  const messages = ref<AiChatMessage[]>([])
  // 工具列表按 messages 下标存放，仅用于展示，不回传给后端。
  const toolsByIndex = ref<Record<number, ToolChip[]>>({})
  const listRef = ref<HTMLElement>()
  // 进行中流的取消句柄，清空/关闭时用来中断，避免泄漏。
  let controller: AbortController | null = null

  const openDrawer = (): void => {
    visible.value = true
  }

  const toolLabel = (tool: ToolChip): string => {
    if (tool.status === 'failed') return t('aiChat.toolFailed', { name: tool.name })
    if (tool.status === 'done') return t('aiChat.toolDone', { name: tool.name })
    return t('aiChat.calling', { name: tool.name })
  }

  const scrollToBottom = (): void => {
    nextTick(() => {
      if (listRef.value) listRef.value.scrollTop = listRef.value.scrollHeight
    })
  }

  /**
   * 发送消息：推入用户消息 + 一条空的助手消息（流式写入目标），带全量历史
   * 调用流式接口。各事件回调直接改写最后一条助手消息的 content / 工具 chip，
   * 触发响应式更新实现逐字渲染。
   */
  const send = (): void => {
    const content = draft.value.trim()
    if (!content || loading.value) return

    messages.value = [
      ...messages.value,
      { role: 'user', content },
      { role: 'assistant', content: '' }
    ]
    const assistantIndex = messages.value.length - 1
    // 只把真正的对话历史（不含刚插入的空助手占位）传给后端。
    const history = messages.value.slice(0, assistantIndex) as AiChatMessage[]
    draft.value = ''
    loading.value = true
    scrollToBottom()

    controller = new AbortController()
    streamAiChat(
      history,
      {
        onMessage: (delta) => {
          const target = messages.value[assistantIndex]
          if (target) target.content += delta
          scrollToBottom()
        },
        onToolCall: (tcall) => {
          const list = toolsByIndex.value[assistantIndex] ?? []
          toolsByIndex.value = {
            ...toolsByIndex.value,
            [assistantIndex]: [...list, { id: tcall.id, name: tcall.name, status: 'running' }]
          }
          scrollToBottom()
        },
        onToolResult: (tres) => {
          const list = toolsByIndex.value[assistantIndex] ?? []
          toolsByIndex.value = {
            ...toolsByIndex.value,
            [assistantIndex]: list.map((chip) =>
              chip.id === tres.id ? { ...chip, status: tres.ok ? 'done' : 'failed' } : chip
            )
          }
        },
        onError: (msg) => {
          ElMessage.error(msg)
          const target = messages.value[assistantIndex]
          if (target && !target.content) target.content = t('aiChat.failed')
          loading.value = false
        },
        onDone: () => {
          loading.value = false
          controller = null
          scrollToBottom()
        }
      },
      controller.signal
    )
  }

  const cancelStream = (): void => {
    if (controller) {
      controller.abort()
      controller = null
    }
    loading.value = false
  }

  /**
   * Enter 发送，Shift+Enter 换行。
   */
  const onEnter = (e: Event | KeyboardEvent): void => {
    if ((e as KeyboardEvent).shiftKey) return
    e.preventDefault()
    send()
  }

  const copyMessage = async (content: string): Promise<void> => {
    if (await copyToClipboard(content)) {
      ElMessage.success(t('aiChat.copied'))
    } else {
      ElMessage.error(t('aiChat.copyFailed'))
    }
  }

  const clearConversation = (): void => {
    cancelStream()
    messages.value = []
    toolsByIndex.value = {}
  }

  onBeforeUnmount(cancelStream)
</script>

<style scoped>
  @reference '@styles/core/tailwind.css';

  .ai-chat-header {
    @apply flex items-center justify-between w-full pr-6;
  }

  .ai-chat-title {
    @apply text-base font-medium;
  }

  .ai-chat-body {
    @apply flex flex-col h-full;
  }

  .ai-chat-list {
    @apply flex-1 overflow-y-auto pr-1;
  }

  .ai-chat-empty {
    @apply mt-10 text-sm text-center text-g-500 select-none;
  }

  .ai-chat-row {
    @apply flex flex-col mb-3;
  }

  .ai-chat-row.is-user {
    @apply items-end;
  }

  .ai-chat-row.is-assistant {
    @apply items-start;
  }

  .ai-chat-bubble {
    @apply max-w-[85%] px-3 py-2 text-sm rounded-lg;

    word-break: break-word;
    white-space: pre-wrap;
  }

  .is-user .ai-chat-bubble {
    color: #fff;
    background-color: var(--main-color);
  }

  .is-assistant .ai-chat-bubble {
    background-color: var(--art-gray-200);
  }

  .ai-chat-thinking {
    @apply text-g-500;
  }

  /* Markdown 渲染内容（v-html，不受 scoped 影响，需用 :deep） */
  .ai-chat-md {
    white-space: normal;
  }

  .ai-chat-md :deep(p) {
    margin: 0 0 8px;
  }

  .ai-chat-md :deep(p:last-child) {
    margin-bottom: 0;
  }

  .ai-chat-md :deep(:is(h1, h2, h3, h4, h5, h6)) {
    margin: 10px 0 6px;
    font-size: 14px;
    font-weight: 600;
  }

  .ai-chat-md :deep(:is(ul, ol)) {
    padding-left: 20px;
    margin: 4px 0;
  }

  .ai-chat-md :deep(li) {
    margin: 2px 0;
  }

  .ai-chat-md :deep(code) {
    padding: 1px 4px;
    font-family: monospace;
    font-size: 12px;
    background: var(--art-gray-300);
    border-radius: 3px;
  }

  .ai-chat-md :deep(pre) {
    padding: 8px 10px;
    margin: 6px 0;
    overflow-x: auto;
    background: var(--art-gray-300);
    border-radius: 6px;
  }

  .ai-chat-md :deep(pre code) {
    padding: 0;
    background: none;
  }

  .ai-chat-md :deep(table) {
    margin: 6px 0;
    font-size: 12px;
    border-collapse: collapse;
  }

  .ai-chat-md :deep(:is(th, td)) {
    padding: 4px 8px;
    border: 1px solid var(--art-gray-400);
  }

  .ai-chat-md :deep(a) {
    color: var(--main-color);
    text-decoration: underline;
  }

  .ai-chat-md :deep(blockquote) {
    padding-left: 10px;
    margin: 6px 0;
    color: var(--art-gray-600);
    border-left: 3px solid var(--art-gray-400);
  }

  .ai-chat-md :deep(hr) {
    margin: 8px 0;
    border: none;
    border-top: 1px solid var(--art-gray-400);
  }

  .ai-chat-copy {
    @apply h-auto p-0 mt-0.5 text-xs text-g-500;
  }

  .ai-chat-tools {
    @apply flex flex-wrap items-center gap-1 mt-1;
  }

  .ai-chat-tools-label {
    @apply text-xs text-g-500;
  }

  .ai-chat-input {
    @apply flex flex-col gap-2 pt-3 mt-2 border-t border-g-300/80;
  }
</style>
