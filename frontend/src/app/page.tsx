// app/page.tsx
"use client";

import { useState } from "react";
import Sidebar from "@/components/Sidebar";
import ChatHeader from "@/components/ChatHeader";
import ChatWindow from "@/components/ChatWindow";
import ChatInput from "@/components/ChatInput";
import AddContactModal from "@/components/AddContactModal";
import CreateGroupModal from "@/components/CreateGroupModal";
import { useChat } from "@/hooks/useChat";

export default function Page() {
  const {
    contacts,
    activeContact,
    selectContact,
    input,
    setInput,
    sendMessage,
    addContact,
    loadingAddContact, // ← add this
  } = useChat();

  const [openAddContact, setOpenAddContact] = useState(false);
  const [openCreateGroup, setOpenCreateGroup] = useState(false);

  return (
    <div className="flex h-screen bg-gray-100">
      <Sidebar
        contacts={contacts}
        activeContact={activeContact}
        onSelect={selectContact}
        onAddContact={() => setOpenAddContact(true)}
        onCreateGroup={() => setOpenCreateGroup(true)}
      />

      <div className="flex flex-col flex-1">
        {activeContact && <ChatHeader contact={activeContact} />}
        <ChatWindow activeContact={activeContact} />
        {activeContact && (
          <ChatInput value={input} setValue={setInput} onSend={sendMessage} />
        )}
      </div>

      {/* Popups */}
      <AddContactModal
        isOpen={openAddContact}
        onClose={() => setOpenAddContact(false)}
        onSubmit={addContact}
        loadingAddContact={loadingAddContact}
      />

      <CreateGroupModal
        isOpen={openCreateGroup}
        onClose={() => setOpenCreateGroup(false)}
        onSubmit={(name) => console.log("Create group:", name)}
      />
    </div>
  );
}
