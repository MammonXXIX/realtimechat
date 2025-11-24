import { Contact } from "@/types/chat";
import { Plus, Users, Search, LogOut, MoreVertical } from "lucide-react";
import { useAuth } from "@clerk/nextjs";
import { useState, useEffect } from "react";

interface Props {
  contacts: Contact[];
  activeContact: Contact | null;
  onSelect: (c: Contact) => void;
  onAddContact?: () => void;

  onViewProfile?: (c: Contact) => void;
  onRemoveContact?: (c: Contact) => void;
}

export default function Sidebar({
  contacts,
  activeContact,
  onSelect,
  onAddContact,
  onViewProfile,
  onRemoveContact,
}: Props) {
  const { signOut } = useAuth();
  const [tab, setTab] = useState<"contacts" | "chats">("chats");
  const [selectedMenuContact, setSelectedMenuContact] = useState<string | null>(
    null
  );

  // Close popup if clicked outside
  useEffect(() => {
    const closeMenu = () => setSelectedMenuContact(null);
    window.addEventListener("click", closeMenu);
    return () => window.removeEventListener("click", closeMenu);
  }, []);

  return (
    <div className="w-1/4 min-w-[260px] bg-white border-r flex flex-col shadow-sm">
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-4 border-b">
        <h2 className="font-bold text-xl text-gray-800">Messages</h2>

        <div className="flex gap-2">
          <button
            onClick={() => signOut({ redirectUrl: "/" })}
            className="p-2 rounded-full bg-red-100 hover:bg-red-200 text-red-600 shadow-sm transition"
            title="Log Out"
          >
            <LogOut size={18} />
          </button>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex border-b">
        <button
          className={`flex-1 py-3 text-center font-semibold transition ${
            tab === "contacts"
              ? "text-indigo-600 border-b-2 border-indigo-600"
              : "text-gray-500 hover:text-gray-700"
          }`}
          onClick={() => setTab("contacts")}
        >
          Contacts
        </button>

        <button
          className={`flex-1 py-3 text-center font-semibold transition ${
            tab === "chats"
              ? "text-indigo-600 border-b-2 border-indigo-600"
              : "text-gray-500 hover:text-gray-700"
          }`}
          onClick={() => setTab("chats")}
        >
          Chats
        </button>
      </div>

      {/* Search */}
      <div className="px-4 py-3">
        <div className="flex items-center gap-2 px-3 py-2 bg-gray-100 rounded-lg">
          <Search size={18} className="text-gray-500" />
          <input
            placeholder={`Search ${tab}...`}
            className="bg-transparent outline-none text-sm w-full"
          />
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto py-3 px-3 space-y-2">
        {/* CONTACTS TAB */}
        {tab === "contacts" && (
          <>
            {contacts.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-8 text-center text-gray-500">
                <p className="text-sm">No contacts yet</p>
                <p className="text-xs mt-1">Add a contact to get started</p>
              </div>
            ) : (
              contacts.map((c) => (
                <div
                  key={c.roomId}
                  onClick={() => onSelect(c)}
                  className="relative flex items-center justify-between gap-3 p-3 rounded-xl bg-white border shadow-sm hover:shadow-md transition cursor-pointer"
                >
                  <div className="flex items-center gap-3">
                    <img
                      src={c.avatar}
                      alt={c.name}
                      className="w-11 h-11 rounded-full object-cover border border-gray-200"
                    />
                    <span className="font-medium text-gray-800">{c.name}</span>
                  </div>
                </div>
              ))
            )}

            <button
              onClick={onAddContact}
              className="flex items-center gap-2 justify-center py-2 px-3 rounded-lg bg-indigo-600 text-white hover:bg-indigo-700 transition shadow-md w-full mt-3"
            >
              <Plus size={18} /> Add Contact
            </button>
          </>
        )}

        {/* CHATS TAB */}
        {tab === "chats" && (
          <>
            {contacts.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-8 text-center text-gray-500">
                <p className="text-sm">No chats yet</p>
                <p className="text-xs mt-1">Add a contact to start chatting</p>
              </div>
            ) : (
              contacts.map((c) => (
                <div
                  key={c.roomId}
                  onClick={() => onSelect(c)}
                  className={`flex items-center gap-3 p-3 rounded-xl cursor-pointer transition relative ${
                    activeContact?.roomId === c.roomId
                      ? "bg-indigo-50 shadow-md border border-indigo-200"
                      : "hover:bg-gray-50 border border-transparent"
                  }`}
                >
                  {/* Avatar with status */}
                  <div className="relative">
                    <img
                      src={c.avatar}
                      alt={c.name}
                      className="w-12 h-12 rounded-full object-cover border border-gray-200"
                    />
                    {c.unread && (
                      <span className="absolute -top-1 -right-1 w-3 h-3 bg-indigo-600 rounded-full border-2 border-white"></span>
                    )}
                  </div>

                  {/* Name and last message */}
                  <div className="flex flex-col w-full overflow-hidden">
                    <span
                      className={`truncate text-sm ${
                        c.unread ? "font-bold text-gray-900" : "text-gray-700"
                      }`}
                    >
                      {c.name}
                    </span>
                    <span
                      className={`truncate text-xs ${
                        c.unread ? "font-medium text-gray-900" : "text-gray-500"
                      }`}
                    >
                      {c.text || "No messages yet"}
                    </span>
                  </div>

                  {/* Unread badge */}
                  {c.unread && (
                    <span className="ml-auto inline-flex items-center justify-center min-w-[20px] h-5 px-1.5 text-xs font-semibold text-white bg-indigo-600 rounded-full">
                      •
                    </span>
                  )}
                </div>
              ))
            )}
          </>
        )}
      </div>
    </div>
  );
}
