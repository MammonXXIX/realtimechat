// src/types/chat.ts
export interface Message {
  from: "me" | "them";
  text: string;
}

// src/types/chat.ts
export interface Contact {
  name: string;
  avatar: string;
  text: string;
  messages: { from: "me" | "them"; text: string }[];
  unread: boolean;
  roomId: string;

  // Tambahan untuk POST /chat
  addedUserId?: string; // <-- dari added_id di backend
}
