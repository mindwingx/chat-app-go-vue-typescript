export interface IWebSocketClient {
    on(type: string, callback: Function): void
    send(data: any): void
    state(): number
    close(): void
}

export type WsResponse = { 
    username: string, 
    content: {
        type:string, 
        id:string, 
        value:string, 
        extra:string[],
    },
    time: string,
}