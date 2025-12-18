import type { Ref } from "vue";

export type Message = { 
    id: string,
    sender: string, 
    text: string, 
    time: string, 
    isOwn: boolean,
}

export interface IUseChat {
  initChat: (wsUrl?: string) => void;
  isConnected: Readonly<Ref<boolean>>;
  message: Ref<Message>;
  notification: Ref<string>;
  typingUsers: Ref<string[]>;
  sendTypingState: (value: boolean) => void;
  sendMessage: (value: Message) => void;
  disconnectChat: () => void;
}