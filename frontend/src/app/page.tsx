"use client";

import { useState } from "react";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";

interface Contact {
  name: string;
  color: string;
  text: string;
  messages: { from: "me" | "them"; text: string }[];
  unread?: boolean;
}

export default function ChatPage() {
  const [contacts, setContacts] = useState<Contact[]>([
    {
      name: "Alice",
      color: "ffa8e4",
      text: "Absolutely 😎",
      messages: [
        { from: "them", text: "Hey Bob! How was your day?" },
        { from: "me", text: "Pretty good! Wrapped up a project finally." },
        { from: "them", text: "Nice! Time to chill then?" },
        { from: "me", text: "Absolutely 😎" },
      ],
      unread: false,
    },
    {
      name: "Martin",
      color: "ad922e",
      text: "Deal!",
      messages: [
        { from: "them", text: "That pizza place was amazing!" },
        { from: "me", text: "Told ya! Their crust is the best in town." },
        { from: "them", text: "We should go again next week!" },
        { from: "me", text: "Deal!" },
      ],
      unread: false,
    },
    {
      name: "Charlie",
      color: "2e83ad",
      text: "Nice, you'll love the visuals!",
      messages: [
        { from: "them", text: "You watched Dune yet?" },
        { from: "me", text: "Not yet! Planning to this weekend." },
        { from: "them", text: "Nice, you'll love the visuals!" },
      ],
      unread: true, // 👈 Belum dibaca
    },
    {
      name: "David",
      color: "c2ebff",
      text: "Adding it to my list!",
      messages: [
        { from: "them", text: "This new mystery novel is wild!" },
        { from: "me", text: "What's it called?" },
        { from: "them", text: "‘The Silent Hours’. Really gripping." },
        { from: "me", text: "Adding it to my list!" },
      ],
      unread: false,
    },
  ]);

  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [activeContact, setActiveContact] = useState<Contact | null>(
    contacts[0]
  );

  const handleSelectContact = (contact: Contact) => {
    // ketika diklik, ubah unread jadi false
    setContacts((prev) =>
      prev.map((c) => (c.name === contact.name ? { ...c, unread: false } : c))
    );
    setActiveContact(contact);
    setSidebarOpen(false);
  };

  const handleAddContact = (name: string, text: string) => {
    const newContact: Contact = {
      name,
      color: Math.floor(Math.random() * 16777215).toString(16),
      text,
      messages: [
        { from: "them", text: `Hey, I'm ${name}! Nice to meet you 👋` },
        { from: "me", text: "Hey there! Welcome to the chat!" },
      ],
      unread: true,
    };
    setContacts((prev) => [...prev, newContact]);
  };

  return (
    <div className="flex h-screen overflow-hidden bg-[#edf2f7] relative">
      {/* Sidebar desktop */}
      <div className="hidden md:block w-1/4 bg-white border-r border-gray-300">
        <Sidebar
          contacts={contacts}
          activeContact={activeContact}
          onSelect={handleSelectContact}
          onClose={() => setSidebarOpen(false)}
          isMobile={false}
          onAddContact={handleAddContact}
        />
      </div>

      {/* Sidebar mobile */}
      {sidebarOpen && (
        <>
          <div
            className="fixed inset-0 backdrop-blur-md bg-white/50 z-40"
            onClick={() => setSidebarOpen(false)}
          ></div>

          <div
            className={`fixed top-0 left-0 w-3/4 sm:w-1/2 h-full bg-white border-r border-gray-300 z-50 transition-transform duration-300 ${
              sidebarOpen ? "translate-x-0" : "-translate-x-full"
            }`}
          >
            <Sidebar
              contacts={contacts}
              activeContact={activeContact}
              onSelect={handleSelectContact}
              onClose={() => setSidebarOpen(false)}
              isMobile={true}
              onAddContact={handleAddContact}
            />
          </div>
        </>
      )}

      {/* Main Chat Area */}
      <div className="flex-1 relative w-full md:w-3/4">
        <header className="bg-white p-4 text-gray-700 border-b flex items-center justify-between">
          <button
            onClick={() => setSidebarOpen(true)}
            className="block md:hidden text-gray-700 focus:outline-none mr-2"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M4 6h16M4 12h16M4 18h16"
              />
            </svg>
          </button>

          <h1 className="text-2xl font-semibold">
            {activeContact ? activeContact.name : "Select a contact"}
          </h1>
        </header>

        <div className="h-screen overflow-y-auto p-4 pb-36">
          {activeContact ? (
            activeContact.messages.map((msg, i) => (
              <div
                key={i}
                className={`flex mb-4 ${
                  msg.from === "me" ? "justify-end" : "justify-start"
                }`}
              >
                {msg.from === "them" && (
                  <div className="w-9 h-9 rounded-full flex items-center justify-center mr-2">
                    <img
                      src={`https://placehold.co/200x/${activeContact.color}/ffffff.svg?text=ʕ•́ᴥ•̀ʔ`}
                      alt="Avatar"
                      className="w-8 h-8 rounded-full"
                    />
                  </div>
                )}
                <div
                  className={`flex max-w-96 rounded-lg p-3 gap-3 ${
                    msg.from === "me"
                      ? "bg-indigo-500 text-white"
                      : "bg-white text-gray-700"
                  }`}
                >
                  <p>{msg.text}</p>
                </div>
                {msg.from === "me" && (
                  <div className="w-9 h-9 rounded-full flex items-center justify-center ml-2">
                    <img
                      src="https://placehold.co/200x/b7a8ff/ffffff.svg?text=ʕ•́ᴥ•̀ʔ"
                      alt="My Avatar"
                      className="w-8 h-8 rounded-full"
                    />
                  </div>
                )}
              </div>
            ))
          ) : (
            <p className="text-center text-gray-500 mt-10">
              Select a contact to start chatting
            </p>
          )}
        </div>

        {/* Chat Input */}
        <footer className="bg-white border-t border-gray-300 p-4 absolute bottom-0 w-full">
          <div className="flex items-center">
            <input
              type="text"
              placeholder="Type a message..."
              className="w-full p-2 rounded-md border border-gray-400 focus:outline-none focus:border-blue-500"
            />
            <button className="bg-indigo-500 text-white px-4 py-2 rounded-md ml-2">
              Send
            </button>
          </div>
        </footer>
      </div>
    </div>
  );
}

/* --- Sidebar --- */
interface SidebarProps {
  contacts: Contact[];
  activeContact: Contact | null;
  onSelect: (c: Contact) => void;
  onClose: () => void;
  isMobile: boolean;
  onAddContact: (name: string, text: string) => void;
}

function Sidebar({
  contacts,
  activeContact,
  onSelect,
  onClose,
  isMobile,
  onAddContact,
}: SidebarProps) {
  const [name, setName] = useState("");
  const [text, setText] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    onAddContact(name, text || "Hey there! I'm using Chat Web 😄");
    setName("");
    setText("");
  };

  return (
    <>
      <header className="p-4 border-b border-gray-300 flex justify-between items-center bg-indigo-600 text-white">
        <h1 className="text-2xl font-semibold">Chat Web</h1>
        <div className="flex items-center gap-2 relative">
          {/* Tombol Add Contact */}
          <Dialog>
            <DialogTrigger asChild>
              <button
                className="focus:outline-none bg-indigo-700 hover:bg-indigo-500 rounded-full p-1.5"
                aria-label="Add Contact"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-5 w-5 text-white"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  strokeWidth={2}
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M12 4v16m8-8H4"
                  />
                </svg>
              </button>
            </DialogTrigger>

            <DialogContent className="sm:max-w-[425px]">
              <DialogHeader>
                <DialogTitle>Add Contact</DialogTitle>
                <DialogDescription>
                  Enter alias name and email for the new contact.
                </DialogDescription>
              </DialogHeader>
              <form onSubmit={handleSubmit} className="grid gap-4">
                <div className="grid gap-3">
                  <Label htmlFor="name">Alias name</Label>
                  <Input
                    id="name"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder="Ziofinga Kurniawan"
                  />
                </div>
                <div className="grid gap-3">
                  <Label htmlFor="text">Email</Label>
                  <Input
                    id="text"
                    value={text}
                    onChange={(e) => setText(e.target.value)}
                    placeholder="ziofinga.kurniawan@gmail.com"
                  />
                </div>
                <DialogFooter>
                  <DialogClose asChild>
                    <Button variant="outline">Cancel</Button>
                  </DialogClose>
                  <Button type="submit">Add</Button>
                </DialogFooter>
              </form>
            </DialogContent>
          </Dialog>

          {isMobile && (
            <button
              onClick={onClose}
              className="focus:outline-none bg-indigo-700 hover:bg-indigo-500 rounded-full p-1.5"
              aria-label="Close sidebar"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="h-4 w-4 text-white"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  fillRule="evenodd"
                  d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                  clipRule="evenodd"
                />
              </svg>
            </button>
          )}
        </div>
      </header>

      {/* daftar kontak */}
      <div className="overflow-y-auto h-screen p-3 mb-9 pb-20">
        {contacts.map((c) => (
          <div
            key={c.name}
            onClick={() => onSelect(c)}
            className={`flex items-center mb-4 cursor-pointer hover:bg-gray-100 p-2 rounded-md ${
              activeContact?.name === c.name ? "bg-indigo-100" : ""
            }`}
          >
            <div className="w-12 h-12 bg-gray-300 rounded-full mr-3 relative">
              <img
                src={`https://placehold.co/200x/${c.color}/ffffff.svg?text=ʕ•́ᴥ•̀ʔ`}
                alt={c.name}
                className="w-12 h-12 rounded-full"
              />
              {c.unread && (
                <span className="absolute top-0 right-0 bg-blue-500 w-3 h-3 rounded-full border-2 border-white"></span>
              )}
            </div>
            <div className="flex-1">
              <h2
                className={`text-lg ${
                  c.unread ? "font-semibold text-gray-900" : "font-medium"
                }`}
              >
                {c.name}
              </h2>
              <p
                className={`text-sm truncate ${
                  c.unread ? "font-semibold text-gray-800" : "text-gray-600"
                }`}
              >
                {c.text}
              </p>
            </div>
          </div>
        ))}
      </div>
    </>
  );
}
