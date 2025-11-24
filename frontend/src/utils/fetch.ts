export async function safeFetchJson(
  url: string,
  options?: RequestInit,
  defaultValue: any = null
) {
  const res = await fetch(url, options);

  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }

  const text = await res.text();
  if (!text) return defaultValue;
  try {
    return JSON.parse(text);
  } catch (error) {
    console.warn("Failed to parse JSON, returning default value:", error);
    return defaultValue;
  }
}
