<template>
  <div class="header bg-gray-800 text-white p-3 px-4 flex items-center justify-center relative">
    <div class="flex-1 flex justify-center">
      <h2 class="text-[1.55rem] font-semibold mt-4">
        {{ headerTitle }}
      </h2>
    </div>
    <div class="right-controls flex items-center">
    </div>
  </div>
  <div class="p-6 flex flex-col gap-3 items-center justify-center flex-1">
    <div>
      <Logo :size="200" class="mb-20"/>
    </div>

    <input v-model="userName" maxlength="20" autofocus @keyup.enter="join"
        class="login-input w-full max-w-[240px] border border-gray-700 rounded-md px-3 py-2 text-sm bg-gray-700 text-white focus:outline-none focus:border-blue-500"
        placeholder="Enter your name"
    >
    <button @click="join" :disabled="!available"
        class="login-btn w-full max-w-[240px] bg-blue-500 text-white px-3 py-2 rounded-md cursor-pointer text-sm disabled:opacity-50 disabled:cursor-not-allowed"
    >Join</button>
  </div>
</template>

<script setup lang="ts">
import Logo from './../components/icons/Logo.vue'
import {ref, onMounted, onUnmounted} from 'vue'
import {useRouter} from "vue-router";
import {handshake} from "./../api/wsHandshankeApi"
import {useUserStore} from "./../stores/user"

const headerTitle = import.meta.env.VITE_APP_NAME

const router = useRouter()
const userStore = useUserStore()

const userName = ref('')
const available = ref(false)

// const randomStr = (len: number = 8): string => {
//   return Math.random().toString(36).substring(2, 2 + len)
// }

const join = () => {
  if (userName.value.trim()) {
    userStore.setUser(userName.value)
    router.push("/chat")
  }
}


onMounted(async () => {
  const check = async () => {
    available.value = await handshake.check()
  }

  check()

  const interval = setInterval(check, 2000)

  onUnmounted(async () => {
    clearInterval(interval)
  })
})

</script>