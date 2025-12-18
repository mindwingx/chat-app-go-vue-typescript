<template>
  <div class="flex flex-col flex-1 min-h-0">
    <ChatHeader :isConnected="isConnected"/>
    <div class="messages-container flex flex-col flex-1 min-h-0 p-2 gap-2 h-[calc(100%-90px)]">
      <ChatNotification :value="notif"/>
      <div class="message-area border border-gray-700 rounded bg-gray-900 h-[90%] relative overflow-hidden">
        <div
          ref="messagesContainer"
          class="absolute inset-0 overflow-y-auto p-2 flex flex-col gap-2 scrollbar-thin"
        >
          <ChatMessage
            v-for="message in messages" :key="message.id"
            :isOwn="message.isOwn"
            :isRtl="isRTL(message.text)"
          >
            <template #username>{{ message.sender }}</template>
            <template #text><span v-html="escapeHtml(message.text)"></span></template>
            <template #time>{{ message.time }}</template>
          </ChatMessage>
        </div>
      </div>
    </div>

    <div class="input-container border-t border-gray-700 p-2 relative flex-shrink-0 h-[80px] bg-gray-800">
      <div
        v-if="typingUsers.length > 0"
        class="typing-indicator absolute -top-7 left-2 right-2 text-xs text-gray-300 italic bg-gray-700 px-2 py-[2px] rounded"
      >
        {{ typingText }}
      </div>
      <ChatInput :isConnected="isConnected" @isTyping="getIsTyping" @inputVal="getInputVal" />
    </div>
  </div>
</template>

<script setup lang="ts">
import ChatHeader from './../components/chat/Header.vue'
import ChatNotification from './../components/chat/Notification.vue'
import ChatMessage from './../components/chat/Message.vue'
import ChatInput from './../components/chat/Input.vue'
import { ref, onMounted, onUnmounted, computed, watch, nextTick, watchEffect } from 'vue'
import { isRTL } from './../utils/helper'
import { useChat } from '../composable/useChat'
import { type Message } from '../types/chat'

const {
  initChat,
  isConnected,
  message,
  notification,
  typingUsers,
  sendTypingState,
  sendMessage,
  disconnectChat,
} = useChat()


const messagesContainer = ref<HTMLElement | null>(null) // Reference to scroll container
const wsConnection = ref<boolean>(false)
const users = ref<string[]>([])
const messages = ref<Message[]>([])
const notif = ref<string>("")

onMounted(() => {
  initChat()
})

watchEffect(() => {
  wsConnection.value = isConnected.value
  notif.value = notification.value
  users.value = typingUsers.value
})

watch(message, (msg) => {
    messages.value.push(msg)
    scrollToBottom()
  }, { deep: true }
)

onUnmounted(() => disconnectChat())


//

function escapeHtml(text:string) :string {
  const div = document.createElement('div')
  div.textContent = text
  return div.innerHTML.replace(/\n/g, '<br>').replace(/\\n/g, '<br>')
}

const scrollToBottom = async () => {
  await nextTick() // Wait for DOM update
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

const typingText = computed(() => {
  if (typingUsers.value.length === 1) {
    return `${typingUsers.value[0]} is typing...`
  } else if (typingUsers.value.length > 1) {
    return `${typingUsers.value.length} people are typing...`
  }
  return ''
})


function getInputVal(val: Message) {
  sendMessage(val)
  messages.value.push(val)
  scrollToBottom()
}

function getIsTyping(val: boolean) {
  sendTypingState(val)
}

</script>