"use client";
import { useEffect, useState } from "react";
import { useAuth, useSession } from "@clerk/nextjs";

export default function ChatHistoryPage() {
  const { getToken } = useAuth();
  const { session } = useSession(); // Tambahkan ini
  
  const [messages, setMessages] = useState([]);
  const [sessionInfo, setSessionInfo] = useState({
    sessionId: null,
    sessionToken: null
  });

  useEffect(() => {
    const fetchHistory = async () => {
      const token = await getToken();
      
      // Simpan session info
      setSessionInfo({
        sessionId: session?.id || null,
        sessionToken: token
      });
      
      console.log("Session ID:", session?.id);
      console.log("Session Token:", token);
      
      const res = await fetch(
        "http://localhost:8081/chat/history",
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
    
    if (session) { // Pastikan session sudah ada
      fetchHistory();
    }
  }, [getToken, session]);

  return (
    <div>
      <div style={{ marginBottom: '20px', padding: '10px', background: '#f0f0f0' }}>
        <h3>Session Info:</h3>
        <p><strong>Session ID:</strong> {sessionInfo.sessionId}</p>
        <p><strong>Session Token:</strong> {sessionInfo.sessionToken?.substring(0, 50)}...</p>
      </div>
      
      <h3>Messages:</h3>
      <pre>{JSON.stringify(messages, null, 2)}</pre>
    </div>
  );
}
