"use client";

import { createContext, useContext } from "react";
import { useChat } from "@/hooks/useChat";
import { Contact } from "@/types/chat";

interface ChatContextType {
  contacts: Contact[];
  activeContact: Contact | null;
  selectContact: (c: Contact) => Promise<void>;
  input: string;
  setInput: (value: string) => void;
  sendMessageWS: () => void;
  addContact: (email: string, aliasName: string) => Promise<void>;
  refreshContacts: () => Promise<void>;
  loadingAddContact: boolean;
}

const ChatContext = createContext<ChatContextType | null>(null);

export function ChatProvider({ children }: { children: React.ReactNode }) {
  const chat = useChat();

  console.log("🔄 ChatProvider render - activeContact:", chat.activeContact);

  return <ChatContext.Provider value={chat}>{children}</ChatContext.Provider>;
}

export function useChatContext() {
  const context = useContext(ChatContext);
  if (!context) {
    throw new Error("useChatContext must be used within ChatProvider");
  }

  console.log(
    "📞 useChatContext called - activeContact:",
    context.activeContact
  );

  return context;
}
