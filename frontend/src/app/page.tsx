
"use client";
import { useEffect, useState } from "react";
import { useAuth } from "@clerk/nextjs";

export default function ChatHistoryPage() {
  const { getToken } = useAuth();
  const [messages, setMessages] = useState([]);

  useEffect(() => {
    const fetchHistory = async () => {
      const token = await getToken()

      const res = await fetch(
        "http://localhost:8081/chat/d89cab7c-f9ed-4f5a-a0b5-6aa5b18b001d/history",
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      const data = await res.json();
      console.log("Response:", data);
      setMessages(data.data);
    };

    fetchHistory();
  }, [getToken]);

  return <pre>{JSON.stringify(messages, null, 2)}</pre>;
}
