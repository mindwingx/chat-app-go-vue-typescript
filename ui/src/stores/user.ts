import {defineStore} from "pinia";

const chatAppSession = "chat-app-user"

const useUserStore = defineStore('user', {
    state: () => ({
        active: false,
        user:  null as string | null,
        users:[] as string[]
    }),
    
    actions: {
        setUser(data: string) {
            this.user = data
            this.active = true
            this.users?.push(this.userName())
        },

        userName() :string {
            return `${this.user} (You)`
        },

        logout() {
            this.user = null
        }
    },

    persist: {
        storage: sessionStorage,
        key: chatAppSession
    }
})

export {chatAppSession, useUserStore}