<!-- AI 运维助手聊天面板 -->
<template>
  <ArtIconButton icon="ri:robot-2-line" :title="t('aiChat.title')" @click="openDrawer" />

  <ElDrawer
    v-model="visible"
    :title="t('aiChat.title')"
    direction="rtl"
    size="420px"
    class="ai-chat-drawer"
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
          <div class="ai-chat-bubble">{{ msg.content }}</div>
          <div v-if="msg.role === 'assistant' && toolsByIndex[index]?.length" class="ai-chat-tools">
            <span class="ai-chat-tools-label">{{ t('aiChat.toolsUsed') }}</span>
            <ElTag v-for="tool in toolsByIndex[index]" :key="tool" size="small" type="info">
              {{ tool }}
            </ElTag>
          </div>
        </div>

        <div v-if="loading" class="ai-chat-row is-assistant">
          <div class="ai-chat-bubble ai-chat-thinking">{{ t('aiChat.thinking') }}</div>
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
  import { ElButton, ElDrawer, ElInput, ElTag } from 'element-plus'
  import { sendAiChat, type AiChatMessage } from '@/api/ai'

  defineOptions({ name: 'ArtAiChat' })

  const { t } = useI18n()

  const visible = ref(false)
  const loading = ref(false)
  const draft = ref('')
  const messages = ref<AiChatMessage[]>([])
  // 工具列表按 messages 下标存放，仅用于展示，不回传给后端。
  const toolsByIndex = ref<Record<number, string[]>>({})
  const listRef = ref<HTMLElement>()

  const openDrawer = (): void => {
    visible.value = true
  }

  const scrollToBottom = (): void => {
    nextTick(() => {
      if (listRef.value) listRef.value.scrollTop = listRef.value.scrollHeight
    })
  }

  /**
   * 发送消息：推入用户消息后带全量历史调用接口，成功后追加助手回复。
   * 失败时只复位 loading，错误提示交由 http 拦截器统一弹出。
   */
  const send = async (): Promise<void> => {
    const content = draft.value.trim()
    if (!content || loading.value) return

    messages.value = [...messages.value, { role: 'user', content }]
    draft.value = ''
    loading.value = true
    scrollToBottom()

    try {
      const res = await sendAiChat(messages.value)
      const index = messages.value.length
      messages.value = [...messages.value, { role: 'assistant', content: res.reply }]
      if (res.tools_used?.length) {
        toolsByIndex.value = { ...toolsByIndex.value, [index]: res.tools_used }
      }
    } catch {
      // 错误已由 http 拦截器提示，这里仅复位状态。
    } finally {
      loading.value = false
      scrollToBottom()
    }
  }

  /**
   * Enter 发送，Shift+Enter 换行。
   */
  const onEnter = (e: Event | KeyboardEvent): void => {
    if ((e as KeyboardEvent).shiftKey) return
    e.preventDefault()
    send()
  }

  const clearConversation = (): void => {
    messages.value = []
    toolsByIndex.value = {}
  }
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
