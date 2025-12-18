import { type IWebSocketClient } from "./type"

export class WebSocketClient implements IWebSocketClient {
    private ws: WebSocket
    private callbacks = new Map<string, Function[]>() // it registers multiple event/handler per map item

    constructor(url:string) {
        this.ws = new WebSocket(url)
        this.ws.onopen = () => this.callbacks.get("open")?.forEach(cb => cb()) 
        this.ws.onclose = (err) => this.callbacks.get("close")?.forEach(cb => cb(err))
        this.ws.onerror = (err) => this.callbacks.get("error")?.forEach(cb => cb(err))
        this.ws.onmessage = (event) => {
            const data = JSON.parse(event.data)
            this.callbacks.get("message")?.forEach(cb => cb(data))
        }
    }

    on(type: string, callback: Function): void {
        if (!this.callbacks.has(type)) {
            this.callbacks.set(type, [])
        }
        this.callbacks.get(type)!.push(callback)
    }

    send(data: any): void {
        if (this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify( data ))
        }
    }

    state(): number {
        return this.ws.readyState
    }

    close(): void {
        this.ws.close()
    }
}
