// src/hooks/useChat.ts
import { useState, useEffect } from "react";
import { useAuth } from "@clerk/nextjs";
import { Contact } from "@/types/chat";
import {
  getChatRooms,
  getChatMessages,
  sendMessageToRoom,
} from "@/services/chatApi";

export function useChat() {
  const { getToken, userId } = useAuth();
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [activeContact, setActiveContact] = useState<Contact | null>(null);
  const [input, setInput] = useState("");
  const [loadingAddContact, setLoadingAddContact] = useState(false);

  // Refresh contacts from API
  const refreshContacts = async () => {
    const token = await getToken();
    if (!token) return;

    const data = await getChatRooms(token);

    const mapped = data.data.map((c: any) => ({
      name: c.chat_room_name,
      avatar:
        c.other_user?.image_url || "https://placehold.co/200x200?text=User",
      text: c.last_message?.message || "", // <-- updated for new structure
      messages: [], // messages will be loaded when selecting contact
      unread: !c.last_message?.is_read, // unread status from last_message
      roomId: c.chat_room_id,
      addedUserId: c.other_user?.id || null,
    }));

    setContacts(mapped);
  };

  useEffect(() => {
    refreshContacts();
  }, [getToken]);

  // Add a new contact
  const addContact = async (email: string, aliasName: string) => {
    setLoadingAddContact(true);
    try {
      const token = await getToken();
      if (!token) return;

      // 1️⃣ Add contact
      const postResponse = await fetch("http://localhost:8081/contacts", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, alias_name: aliasName }),
      });

      if (!postResponse.ok) throw new Error("Gagal menambahkan kontak");

      // 2️⃣ Get all contacts
      const getContactsResponse = await fetch(
        "http://localhost:8081/contacts",
        {
          method: "GET",
          headers: { Authorization: `Bearer ${token}` },
        }
      );

      const contactsData = await getContactsResponse.json();

      // 3️⃣ Find the newly added contact (based on email)
      const matchedContact = contactsData?.data?.find(
        (c: any) => c.added_user?.email === email
      );

      if (!matchedContact)
        throw new Error("Kontak tidak ditemukan setelah ditambahkan");

      const addedUserId = matchedContact.added_user.id;

      // 4️⃣ Create chat room for the new contact
      await fetch("http://localhost:8081/chat", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ added_user_id: addedUserId }),
      });

      // 5️⃣ Refresh contacts
      await refreshContacts();
    } finally {
      setLoadingAddContact(false);
    }
  };

  // Select a contact (load chat messages)
  const selectContact = async (contact: Contact) => {
    const token = await getToken();
    if (!token) return;

    const history = await getChatMessages(token, contact.roomId);

    // ✅ Safely map messages, default to empty array if undefined
    const mappedMessages = (history?.data || []).map((msg: any) => ({
      from: msg.sender_id === userId ? "me" : "them",
      text: msg.message,
    }));

    const updated = { ...contact, messages: mappedMessages, unread: false };

    setContacts((prev) =>
      prev.map((c) => (c.roomId === contact.roomId ? updated : c))
    );
    setActiveContact(updated);
  };

  // Send message
  const sendMessage = async () => {
    if (!activeContact || !input.trim()) return;
    const token = await getToken();
    if (!token) return;

    // Send the message
    await sendMessageToRoom(token, activeContact.roomId, input);

    // Update activeContact messages and last text
    setActiveContact((prev) => {
      if (!prev) return prev;

      const updated: Contact = {
        ...prev,
        messages: [
          ...prev.messages,
          { from: "me" as "me", text: input }, // ← cast as "me"
        ],
        text: input, // update last message
        unread: false,
      };

      // Update contacts array
      setContacts((prevContacts) =>
        prevContacts.map((c) =>
          c.roomId === updated.roomId ? { ...c, text: input } : c
        )
      );

      return updated;
    });

    setInput("");
  };

  return {
    contacts,
    activeContact,
    selectContact,
    input,
    setInput,
    sendMessage,
    addContact,
    refreshContacts,
    loadingAddContact,
  };
}
