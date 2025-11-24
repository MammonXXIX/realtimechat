// components/ChatWindow.tsx
import { Contact } from "@/types/chat";

export default function ChatWindow({
  activeContact,
}: {
  activeContact: Contact | null;
}) {
  if (!activeContact) {
    return (
      <div className="flex-1 flex items-center justify-center text-gray-500">
        Select a contact to start chatting
      </div>
    );
  }

  return (
    <div className="flex-1 p-4 overflow-y-auto bg-gray-50">
      {activeContact.messages.map((m, idx) => (
        <div
          key={idx}
          className={`mb-2 flex ${
            m.from === "me" ? "justify-end" : "justify-start"
          }`}
        >
          <div
            className={`px-3 py-2 rounded-lg max-w-[70%] ${
              m.from === "me"
                ? "bg-indigo-600 text-white"
                : "bg-gray-200 text-gray-800"
            }`}
          >
            {m.text}
          </div>
        </div>
      ))}
    </div>
  );
}
