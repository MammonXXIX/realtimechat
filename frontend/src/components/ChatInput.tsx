import { useChatContext } from "@/context/ChatContext";

export default function ChatInput() {
  const { input, setInput, sendMessageWS, activeContact } = useChatContext();

  console.log("🔍 ChatInput - activeContact:", activeContact);
  console.log("🔍 ChatInput - input:", input);

  const disabled = !input.trim();

  const handleSend = () => {
    console.log("🔘 Send button clicked!");
    console.log("Disabled state:", disabled);
    console.log("Input value:", input);

    if (!disabled) {
      console.log("✅ Calling sendMessageWS...");
      sendMessageWS();
    } else {
      console.log("❌ Button is disabled, not sending");
    }
  };

  return (
    <div className="p-4 border-t bg-white flex gap-3">
      <input
        className="flex-1 border rounded-md px-3 py-2"
        placeholder="Type a message..."
        value={input}
        onChange={(e) => {
          console.log("Input changed:", e.target.value);
          setInput(e.target.value);
        }}
        onKeyDown={(e) => {
          if (e.key === "Enter") {
            console.log("Enter key pressed");
            handleSend();
          }
        }}
        onClick={() => console.log("Input clicked!")}
      />

      <button
        onClick={handleSend}
        disabled={disabled}
        className={`px-4 py-2 rounded-md text-white 
          ${disabled ? "bg-gray-400" : "bg-indigo-600 active:scale-95"}
        `}
      >
        Send
      </button>
    </div>
  );
}
