<template>
  <div class="flex items-start gap-2" :class="[ isOwn ? 'justify-end' : 'justify-start']">
    <div v-if="!isOwn"
         class="flex-shrink-0 w-8 h-8 rounded-full bg-slate-500 text-white flex items-center justify-center font-bold">
      {{ userChar?.charAt(0).toUpperCase() }}
    </div>

    <div class="msg p-1.5 rounded-lg max-w-[75%] min-w-[200px] break-words relative text-sm"
         :class="[ isOwn ? 'bg-indigo-500 text-white' : 'bg-blue-500 text-white']">
      
      <span v-if="!isOwn" class="text-white/60 block text-xs font-semibold mb-1">
        <slot name="username" />
      </span>

      <span :class="isRtl ? 'block text-right rtl' : 'block text-left'">
        <slot name="text" />
      </span>

      <div class="meta text-xs text-white/40 mt-1.5" :class="[!isOwn ? 'text-right' : 'text-left']">
        <slot name="time" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useSlots } from 'vue';

const slots = useSlots()
const userChar = String(slots.username?.()[0]?.children ?? "")

defineProps<{ 
    isOwn?: boolean,
    isRtl?: boolean,
}>();
</script>

<style scoped>
.rtl {
  direction: rtl;
}
</style>

