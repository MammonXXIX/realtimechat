// src/components/ChatList.tsx
import { Contact } from "@/types/chat";

interface Props {
  contacts: Contact[];
  activeContact: Contact | null;
  selectContact: (contact: Contact) => void;
}

export default function ChatList({
  contacts,
  activeContact,
  selectContact,
}: Props) {
  return (
    <div className="w-1/3 border-r">
      {contacts.map((c) => (
        <div
          key={c.roomId}
          onClick={() => selectContact(c)}
          className={`p-4 cursor-pointer hover:bg-gray-100 ${
            activeContact?.roomId === c.roomId ? "bg-gray-100" : ""
          }`}
        >
          <div className="flex items-center gap-3">
            <img src={c.avatar} className="w-12 h-12 rounded-full" />
            <div>
              <p className="font-semibold">{c.name}</p>
              <p className="text-sm text-gray-600">{c.text}</p>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
