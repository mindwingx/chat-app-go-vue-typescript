<template>
  <div class="flex-shrink-0 header bg-gray-800 text-white p-3 px-4 flex items-center justify-center relative">
    <button
        @click="toggleUsersPopup"
        class="users-toggle bg-gray-700 text-white text-[11px] px-2 py-1 rounded cursor-pointer"
        :class="{ 'hidden': !userStore.active }"
    >Online Users</button>

    <div id="users" class="absolute left-3 top-full -mt-2 bg-gray-600 border border-gray-700 rounded-md shadow p-2 min-w-[140px] z-[1000] text-white text-xs scrollbar-thin"
      v-show="showUsersPopup">
      <div class="flex flex-col gap-1 max-h-[180px] overflow-y-auto">
        <div v-for="(user, index) in userStore.users" :key="index" class="px-2 py-1 rounded">
          <span class="status-dot" :class="[isConnected ? 'status-online' : 'status-offline']"></span>{{ user }}
        </div>
      </div>
    </div>

    <div class="flex-1 flex justify-center">
      <h2 v-if="!userStore.active" class="text-[1.55rem] font-semibold mt-4">
        {{ headerTitle }}
      </h2>
      <div v-else>
        <Logo :size="32"/>
      </div>
    </div>
    
    <div v-if="userStore.active" class="user-info text-xs opacity-90 mx-2">
      [ {{ userStore.user }} ]
    </div>

    <div class="right-controls flex items-center">
      <button
          v-if="userStore.active"
          @click.prevent="exit"
          class="exit-btn bg-gray-700 text-white text-xs px-2 py-1 rounded cursor-pointer"
      >
        Exit
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import Logo from './../icons/Logo.vue'
import { useRouter } from 'vue-router'
import { useUserStore } from "./../../stores/user"
import { chatAppSession } from './../../stores/user'

let headerTitle = import.meta.env.VITE_APP_NAME
const router = useRouter()
const userStore = useUserStore()

const showUsersPopup = ref(false)

defineProps<{
  isConnected: Boolean,
}>()

const toggleUsersPopup = () => {
  showUsersPopup.value = !showUsersPopup.value
}

// Close emoji popup when clicking outside
const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  if (!target.closest('.users-toggle') && !target.closest('#users')) {
    showUsersPopup.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})

const exit = () => {
  userStore.$reset()
  sessionStorage.removeItem(chatAppSession)
  router.push('/')
}
</script>

<style scoped>
.status-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
  margin-right: 6px;
  vertical-align: middle;
}


.status-online {
  background: #84cc16;
}

.status-offline {
  background: #e9682c;
}
</style>