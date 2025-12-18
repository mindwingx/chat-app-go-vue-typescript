import { onUnmounted, ref } from "vue";
import { type IWebSocketClient } from "./../lib/websocket/type"
import { WebSocketClient } from "./../lib/websocket/client"
import { type IUseWebSocket } from "../types/websocket";

function useWebSocket(url:string): IUseWebSocket {
    const client = ref<IWebSocketClient>()

    const isConnected = ref(false)
    const event = ref<string>("")
    const handlers = new Map<string, Function>()

    let   reconnecting = false
    const retries = ref<number>(0)
    const maxRetries:number = 5

    const connect = () => {
        if (retries.value >= maxRetries) {
            event.value = "max retry attempts reached"
            return
        }

        client.value = new WebSocketClient(url)

        // re-attach all handlers
        handlers.forEach((callback, type) => {
            client.value!.on(type, callback)
        })
        
        client.value.on("open", () => {
            isConnected.value = true
            retries.value = 0 // reset retry count
            event.value = "connected to server"
        })

        client.value.on("close", (_err:any) => reconnect())

        client.value.on("error", (_err:any) => reconnect())
    }

    const reconnect = () => {
        if (reconnecting || retries.value >= maxRetries) return

        reconnecting = true

        isConnected.value = false
        retries.value++
        event.value = `reconnecting attempt ${retries.value}...`

        setTimeout(() => {
            reconnecting = false
            connect()
        }, 2000 * retries.value)
    }

    const send = (data: any) => {
        client.value?.send(data)
    }

    const on = (type: string, callback: Function) => {
        handlers.set(type, callback) // store permanently
        client.value?.on(type, callback) // attach immediately if connected
    }

    const disconnect = () => {
        retries.value = maxRetries; // stop any further retries
        client.value?.close()
        
        const interval = setInterval(() => {
            if (client.value?.state() === WebSocket.CLOSED) {
                isConnected.value = false
                event.value = "disconnected"
                clearInterval(interval)
            }
        }, 200)
    }

    onUnmounted(() => {
        disconnect()
    })

    return { connect, send, on, disconnect, isConnected, event }
}

export {useWebSocket}