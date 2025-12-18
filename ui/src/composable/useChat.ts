import { ref, watch, readonly } from "vue"
import { useUserStore } from "./../stores/user"
import { useWebSocket } from './../composable/useWebsocket'
import { type WsResponse } from "../lib/websocket/type"
import { type Message, type IUseChat } from "../types/chat"

function useChat (wsUrl?: string): IUseChat {
    const userStore = useUserStore()
    const url = wsUrl || `${import.meta.env.VITE_WS_URL}/ws?username=${userStore.user}`
    const { connect, send, on, disconnect, isConnected, event } = useWebSocket(url)

    let initialized = false // to avoid calling "initChat" multiple times

    const message = ref<Message>({id: "",sender: "", text: "", time: "", isOwn: false})
    const notification = ref<string>("")
    const typingUsers = ref<string[]>([])

    watch(event, (val) => {
        // handles the WS server notifications
        let time = new Date().toLocaleTimeString([], { 
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
            hour12: false,
        })

        notification.value = `[${time}] ${val}`
    })

    const initChat = () => {
        if (initialized) return
        initialized = true

        connect() // connection error is handled by "onerror" event handler

        on("open", () => {
            console.log("on open connected from useChat") // to test register multiple handlers for an event in WS
        })

        on("message", (data: WsResponse) => {
            switch (data.content.type) {
            case "message":
                message.value = {
                    id: data.content.id,
                    sender: data.username,
                    text: data.content.value,
                    time: data.time,
                    isOwn: userStore.user == data.username,
                }
                break
            case "notification":
                // handles the chat events notifications
                if (data.content.extra != null && data.content.extra.length > 0) {
                    userStore.users = userStore.users.filter(user => {
                        if (user == userStore.userName()) return true
                        return data.content.extra.includes(user)
                    })

                    data.content.extra.forEach (user => {
                        if (user == userStore.user) return

                        if (!userStore.users.includes(user)) userStore.users.push(user)
                    })
                } else {
                    userStore.users = []
                }

                notification.value = `[${data.time}] ${data.content.value}`
                break
            case "online-users":
                // initialize online users list
                if (data.content.extra != null && data.content.extra.length > 0) {
                    for (const u of data.content.extra) {
                        if(u != userStore.user && userStore.users.includes(u) == false){
                        userStore.users.push(u)
                        }
                    }
                }
                break
            case "typing":
                typingUsers.value= data.content.extra.filter(user => userStore.user != user)
                break
            }
        })
    }

    const sendTypingState = (value:boolean) => {
        if (value == true) send({ type:"typing", value: "typing" })
    }

    const sendMessage = (value: Message) => {
        send({ type:"message", value: value.text })
    }

    const disconnectChat = () => disconnect()

    return {
        initChat,
        isConnected: readonly(isConnected),
        message, 
        notification, 
        typingUsers, 
        sendTypingState,
        sendMessage, 
        disconnectChat,
     }
}

export { useChat }