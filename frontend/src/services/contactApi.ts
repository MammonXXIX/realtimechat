export async function addContact(email: string, alias_name: string) {
  const res = await fetch("http://localhost:8081/contacts", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, alias_name }),
  });

  if (!res.ok) {
    throw new Error("Failed to add contact");
  }

  return await res.json();
}
