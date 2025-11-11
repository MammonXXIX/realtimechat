import { useState, useEffect } from "react";
import { useAuth } from "@clerk/nextjs";

export function useChatHistory() {
  const { getToken } = useAuth();
  const [histories, setHistories] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchHistory = async () => {
    setLoading(true);
    const token = await getToken();
    if (!token) return;

    const res = await fetch("http://localhost:8081/chat/history", {
      headers: { Authorization: `Bearer ${token}` },
    });

    const data = await res.json();
    setHistories(data.data || []);
    setLoading(false);
  };

  useEffect(() => {
    fetchHistory();
  }, []);

  return { histories, loading, fetchHistory };
}
