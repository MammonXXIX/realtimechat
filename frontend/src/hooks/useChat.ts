// src/hooks/useChat.ts
import { useState, useEffect, useRef, useCallback } from "react";
import { useAuth } from "@clerk/nextjs";
import { Contact } from "@/types/chat";
import { getChatRooms, getChatMessages } from "@/services/chatApi";
import { safeFetchJson } from "@/utils/fetch";
import toast from "react-hot-toast";

export function useChat() {
  const { getToken, userId } = useAuth();
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const reconnectAttemptsRef = useRef(0);

  const [chatRooms, setChatRooms] = useState<any[]>([]);
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [activeContact, setActiveContact] = useState<Contact | null>(null);
  const [input, setInput] = useState("");
  const [loadingAddContact, setLoadingAddContact] = useState(false);
  const [wsConnected, setWsConnected] = useState(false);
  const hasJoinedRoomsRef = useRef(false);

  // ===================================================
  // FETCH CHAT ROOMS
  // ===================================================
  const refreshContacts = useCallback(async () => {
    const token = await getToken();
    if (!token) return;

    const data = await getChatRooms(token);
    const rooms = data.data || [];

    setChatRooms(rooms);

    const mapped = rooms.map((c: any) => ({
      name: c.chat_room_name,
      avatar:
        c.other_user?.image_url || "https://placehold.co/200x200?text=User",
      text: c.last_message?.message || "",
      messages: [],
      unread: !c.last_message?.is_read,
      roomId: c.chat_room_id,
      addedUserId: c.other_user?.id || null,
    }));

    setContacts(mapped);
  }, [getToken]);

  useEffect(() => {
    refreshContacts();
  }, [refreshContacts]);

  // ===================================================
  // INIT WEBSOCKET WITH AUTO-RECONNECT
  // ===================================================
  const connectWebSocket = useCallback(async () => {
    const token = await getToken();
    if (!token) return;

    // Clean up existing connection
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }

    const ws = new WebSocket(`ws://localhost:8081/websocket?token=${token}`);
    wsRef.current = ws;

    ws.onopen = () => {
      setWsConnected(true);
      reconnectAttemptsRef.current = 0;
      hasJoinedRoomsRef.current = false;

      // Join rooms after a small delay to ensure connection is stable
      setTimeout(() => {
        if (ws.readyState === WebSocket.OPEN && chatRooms.length > 0) {
          const roomIds = chatRooms.map((h) => h.chat_room_id);
          const payload = { type: "join_rooms", rooms: roomIds };
          ws.send(JSON.stringify(payload));
          hasJoinedRoomsRef.current = true;
        }
      }, 100);
    };

    ws.onmessage = (event) => {
      try {
        const incoming = JSON.parse(event.data);

        if (incoming.chat_room_id && incoming.message) {
          // Update contacts list
          setContacts((prev) =>
            prev.map((c) =>
              c.roomId === incoming.chat_room_id
                ? {
                    ...c,
                    text: incoming.message,
                    unread: incoming.sender_id !== userId,
                  }
                : c
            )
          );

          // Update active chat
          setActiveContact((prev) => {
            if (!prev || prev.roomId !== incoming.chat_room_id) {
              return prev;
            }

            const from: "me" | "them" =
              incoming.sender_id === userId ? "me" : "them";

            return {
              ...prev,
              messages: [
                ...prev.messages,
                { from, text: incoming.message as string },
              ],
              text: incoming.message,
            };
          });
        }
      } catch (error) {
        console.error("Error parsing WebSocket message:", error);
      }
    };

    ws.onerror = () => {
      setWsConnected(false);
    };

    ws.onclose = (event) => {
      setWsConnected(false);
      wsRef.current = null;
      hasJoinedRoomsRef.current = false;

      // Auto-reconnect if not a clean close
      if (!event.wasClean) {
        reconnectAttemptsRef.current++;
        const delay = Math.min(1000 * reconnectAttemptsRef.current, 5000);

        reconnectTimeoutRef.current = setTimeout(() => {
          connectWebSocket();
        }, delay);
      }
    };
  }, [getToken, userId, chatRooms]);

  useEffect(() => {
    connectWebSocket();

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.close();
      }
    };
  }, [connectWebSocket]);

  // ===================================================
  // JOIN ROOMS WHEN CHAT ROOMS CHANGE
  // ===================================================
  useEffect(() => {
    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
      return;
    }

    if (!chatRooms.length) {
      return;
    }

    if (hasJoinedRoomsRef.current) {
      return;
    }

    const timeoutId = setTimeout(() => {
      if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        const roomIds = chatRooms.map((h) => h.chat_room_id);
        const payload = { type: "join_rooms", rooms: roomIds };

        try {
          wsRef.current.send(JSON.stringify(payload));
          hasJoinedRoomsRef.current = true;
        } catch (error) {
          console.error("Failed to send join_rooms:", error);
        }
      }
    }, 100);

    return () => clearTimeout(timeoutId);
  }, [chatRooms, wsConnected]);

  // ===================================================
  // SELECT CONTACT
  // ===================================================
  const selectContact = async (contact: Contact) => {
    const token = await getToken();
    if (!token) return;

    try {
      const history = await getChatMessages(token, contact.roomId);
      const mappedMessages = (history?.data || []).map((msg: any) => {
        const from: "me" | "them" = msg.sender_id === userId ? "me" : "them";
        return { from, text: msg.message as string };
      });

      const updated: Contact = {
        ...contact,
        messages: mappedMessages,
        unread: false,
      };

      setContacts((prev) =>
        prev.map((c) => (c.roomId === contact.roomId ? updated : c))
      );

      setActiveContact(updated);
    } catch (error) {
      console.error("Error selecting contact:", error);
    }
  };

  // ===================================================
  // SEND MESSAGE
  // ===================================================
  const sendMessageWS = () => {
    if (!input.trim()) {
      return;
    }

    if (!activeContact) {
      return;
    }

    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
      connectWebSocket();
      alert("Connection lost. Reconnecting... Please try again in a moment.");
      return;
    }

    const messageText = input.trim();
    const payload = {
      type: "send_message",
      room_id: activeContact.roomId,
      message: messageText,
    };

    try {
      wsRef.current.send(JSON.stringify(payload));
      setInput("");
    } catch (error) {
      console.error("Failed to send message:", error);
      alert("Failed to send message. Please try again.");
    }
  };

  // ===================================================
  // ADD CONTACT
  // ===================================================
  const addContact = async (email: string, aliasName: string) => {
    setLoadingAddContact(true);
    try {
      const token = await getToken();
      if (!token) {
        toast.error("Token tidak tersedia. Silakan login ulang.");
        return;
      }

      // 1️⃣ Add contact
      const postResponse = await fetch("http://localhost:8081/contacts", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, alias_name: aliasName }),
      });

      if (!postResponse.ok) {
        toast.error("Gagal menambahkan kontak");
        return;
      }

      // 2️⃣ Get all contacts
      const contactsData = await safeFetchJson(
        "http://localhost:8081/contacts",
        { headers: { Authorization: `Bearer ${token}` } },
        { data: [] } // default kalau kosong
      );

      // 3️⃣ Find the newly added contact
      const matchedContact = contactsData?.data?.find(
        (c: any) => c.added_user?.email === email
      );

      if (!matchedContact) {
        toast.error("Kontak tidak ditemukan setelah ditambahkan");
        return;
      }

      const addedUserId = matchedContact.added_user.id;

      // 4️⃣ Create chat room
      const chatResponse = await fetch("http://localhost:8081/chat", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ added_user_id: addedUserId }),
      });

      if (!chatResponse.ok) {
        toast.error("Gagal membuat chat room untuk kontak baru");
        return;
      }

      // 5️⃣ Refresh contacts
      await refreshContacts();
      toast.success("Kontak berhasil ditambahkan!");
    } catch (error: any) {
      console.error(error);
      toast.error(error?.message || "Terjadi kesalahan");
    } finally {
      setLoadingAddContact(false);
    }
  };

  return {
    contacts,
    activeContact,
    selectContact,
    input,
    setInput,
    sendMessageWS,
    addContact,
    refreshContacts,
    loadingAddContact,
    wsConnected,
  };
}
