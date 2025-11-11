// components/ChatInput.tsx
export default function ChatInput({
  value,
  setValue,
  onSend,
}: {
  value: string;
  setValue: (v: string) => void;
  onSend: () => void;
}) {
  return (
    <div className="p-4 border-t bg-white flex gap-3">
      <input
        className="flex-1 border rounded-md px-3 py-2"
        placeholder="Type a message..."
        value={value}
        onChange={(e) => setValue(e.target.value)}
      />
      <button
        onClick={onSend}
        className="bg-indigo-600 text-white px-4 py-2 rounded-md"
      >
        Send
      </button>
    </div>
  );
}
