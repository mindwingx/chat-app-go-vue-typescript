<template>
    <div class="notifications-area border border-gray-700 rounded bg-gray-900 h-[10%] overflow-hidden relative">
        <div ref="notificationContainer" class="absolute inset-0 overflow-y-auto px-2 py-1 flex flex-col gap-1 text-[11px] scrollbar-thin">
          <div v-for="(notification, index) in notifications" :key="index" class="text-gray-300">
            {{ notification }}
          </div>
        </div>
      </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'

const notifications = ref<string[]>([])
const notificationContainer = ref<HTMLElement | null>(null) // Reference to scroll container


const props = defineProps<{
  value?: string | null
}>()

watch(() => props.value, (newVal) => {
  if (newVal) {
    notifications.value.push(newVal)
    scrollToBottom()
  }
})

const scrollToBottom = async () => {
  await nextTick() // Wait for DOM update
  if (notificationContainer.value) {
    notificationContainer.value.scrollTop = notificationContainer.value.scrollHeight
  }
}

</script>