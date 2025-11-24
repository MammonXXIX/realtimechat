// src/services/chatApi.ts
export async function getChatRooms(token: string) {
  return fetch("http://localhost:8081/chat/history", {
    headers: { Authorization: `Bearer ${token}` },
  }).then((r) => r.json());
}

export async function getChatMessages(token: string, chatRoomID: string) {
  return fetch(`http://localhost:8081/chat/${chatRoomID}/history`, {
    headers: { Authorization: `Bearer ${token}` },
  }).then((r) => r.json());
}

// src/services/chatApi.ts
export async function sendMessageToRoom(
  token: string,
  roomId: string,
  message: string
) {
  const res = await fetch(`http://localhost:8081/chat/${roomId}/message`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ message }),
  });

  if (!res.ok) {
    throw new Error("Failed to send message");
  }

  // 🔥 backend kemungkinan return 204 No Content, jadi jangan pakai res.json()
  try {
    return await res.json();
  } catch {
    return null; // ✅ aman, tidak crash
  }
}
