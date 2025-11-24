// src/services/contactApi.ts

export async function addContact(
  token: string,
  email: string,
  alias_name: string
) {
  const res = await fetch("http://localhost:8081/contacts", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, alias_name }),
  });

  if (!res.ok) {
    const responseText = await res.text();
    let errorMessage = "Failed to add contact";

    try {
      const errorData = JSON.parse(responseText);
      if (errorData.error && errorData.error.message) {
        errorMessage = errorData.error.message;
      } else if (errorData.message) {
        errorMessage = errorData.message;
      }
    } catch (e) {
      errorMessage = responseText || `HTTP ${res.status}`;
    }

    throw new Error(errorMessage);
  }

  // Periksa apakah ada body sebelum parse JSON
  const text = await res.text();
  if (!text) return null; // atau return {} sesuai kebutuhan
  return JSON.parse(text);
}
