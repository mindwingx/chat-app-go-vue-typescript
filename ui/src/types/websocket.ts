export interface IUseWebSocket {
    connect: () => void;
    send: (data: any) => void;
    on: (type: string, callback: Function) => void;
    disconnect: () => void;
    isConnected: { readonly value: boolean };
    event: { readonly value: string };
}