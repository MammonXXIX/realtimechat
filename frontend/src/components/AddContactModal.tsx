// components/AddContactModal.tsx
"use client";

import { useState } from "react";

export default function AddContactModal({
  isOpen,
  onClose,
  onSubmit,
  loadingAddContact, // ← NEW
}: {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (email: string, aliasName: string) => Promise<void>;
  loadingAddContact: boolean; // ← NEW
}) {
  const [email, setEmail] = useState("");
  const [aliasName, setAliasName] = useState("");

  if (!isOpen) return null;

  const handleSubmit = async () => {
    if (!email.trim() || !aliasName.trim()) return;

    await onSubmit(email, aliasName); // tunggu proses selesai
    setEmail("");
    setAliasName("");
    onClose(); // tutup kalau sukses
  };

  return (
    <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-white w-full max-w-sm rounded-xl shadow-lg overflow-hidden animate-fadeIn">
        <div className="px-5 py-4 border-b bg-indigo-600 text-white">
          <h2 className="text-lg font-semibold">Add Contact</h2>
        </div>

        <div className="p-5 space-y-4">
          <div>
            <label className="text-sm font-medium text-gray-600">Email</label>
            <input
              type="text"
              placeholder="User Email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="mt-1 w-full border rounded-md px-3 py-2 focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>

          <div>
            <label className="text-sm font-medium text-gray-600">
              Alias Name
            </label>
            <input
              type="text"
              placeholder="Alias Name for User"
              value={aliasName}
              onChange={(e) => setAliasName(e.target.value)}
              className="mt-1 w-full border rounded-md px-3 py-2 focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
        </div>

        <div className="px-5 py-3 border-t flex justify-end gap-2">
          <button
            onClick={onClose}
            className="px-4 py-2 rounded-md hover:bg-gray-100"
          >
            Cancel
          </button>
          <button
            disabled={loadingAddContact}
            onClick={handleSubmit}
            className={`px-4 py-2 rounded-md text-white ${
              loadingAddContact
                ? "bg-indigo-400 cursor-not-allowed"
                : "bg-indigo-600 hover:bg-indigo-700"
            }`}
          >
            {loadingAddContact ? "Menambahkan..." : "Add"}
          </button>
        </div>
      </div>
    </div>
  );
}
