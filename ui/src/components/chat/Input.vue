<template>
    <div class="flex items-end gap-2 h-full relative">
      <button @click="toggleEmojiPopup"
        class="emoji-toggle flex items-center justify-center bg-gray-700 text-white rounded w-10 h-10 text-lg"
      >
        😊
      </button>

      <div v-show="showEmojiPopup && isConnected"
        class="emoji-popup absolute left-0 bottom-[75%] mb-0 bg-gray-800 border border-gray-700 rounded shadow p-2 min-w-[160px] z-[1000] scrollbar-thin"
      >
        <div class="emoji-list grid grid-cols-[repeat(auto-fill,minmax(30px,1fr))] gap-1 max-h-[200px] overflow-y-auto">
          <span v-for="emoji in emojis" :key="emoji" @click="addEmoji(emoji)" 
            class="cursor-pointer text-center hover:bg-gray-700 rounded p-1">
            {{ emoji }}
          </span>
        </div>
      </div>

      <textarea :disabled="!isConnected" @keydown="handleTyping" v-model="message" @input="handleTyping" @keydown.enter.exact.prevent="sendMessage" @click="showEmojiPopup = false"
        rows="1" maxlength="1000" :placeholder="isConnected ? 'Type message...' : 'locked (disconnected)'"
        class="flex-1 border border-gray-700 rounded px-3 py-2 text-sm bg-gray-700 text-white resize-none min-h-[40px] max-h-[80px] focus:outline-none focus:border-blue-500"
      ></textarea>

      <button @click="sendMessage" :disabled="!message.trim() || !isConnected"
        class="bg-blue-500 disabled:bg-gray-500 text-white px-3 py-2 rounded h-10 min-w-[60px] text-sm"
      >Send</button>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { type Message } from '../../types/chat';

defineProps({
  isConnected: {type: Boolean}
})

const message = ref('')
const showEmojiPopup = ref(false)
const emojis = ['😊','😀', '😂', '❤️', '🔥', '👍', '🎉', '🤔', '😎', '💯', '🤯', '👀']

const toggleEmojiPopup = () => {
  showEmojiPopup.value = !showEmojiPopup.value
}

const addEmoji = (emoji: string) => {
  message.value += emoji
}

// Close emoji popup when clicking outside
const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  if (!target.closest('.emoji-toggle') && !target.closest('.emoji-popup')) {
    showEmojiPopup.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})

//

const handleTyping = () => {
  if (message.value.length > 0)
    emit("isTyping", true)
}

const emit = defineEmits(["inputVal", "isTyping"])

const sendMessage = () => {
  if (!message.value.trim()) return

  const newMessage: Message = {
    id: "", // ws handles this field
    sender: "", // ws handles this field
    text: message.value,
    time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false }),
    isOwn: true
  }

  emit("inputVal", newMessage)

  //todo: send new message. then, do the steps below:

  message.value = ''
  showEmojiPopup.value = false
}
</script>

<style scoped>
.emoji-popup {
  display: block;
}
</style>