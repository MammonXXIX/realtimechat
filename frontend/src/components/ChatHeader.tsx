// components/ChatHeader.tsx
import { Contact } from "@/types/chat";

export default function ChatHeader({ contact }: { contact: Contact | null }) {
  return (
    <div className="h-17 px-4 border-b bg-white flex items-center gap-3 shadow-sm">
      {contact ? (
        <>
          <img
            src={contact.avatar}
            className="w-10 h-10 rounded-full object-cover"
          />
          <div className="flex flex-col">
            <span className="font-semibold text-lg">{contact.name}</span>
          </div>
        </>
      ) : (
        <span className="text-gray-500 text-base">
          Select a contact to start chatting
        </span>
      )}
    </div>
  );
}
