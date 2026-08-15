// Thin fetch wrapper around the Go backend.
//
// Two of the backend's list handlers (GET /rooms, GET /plugs) call
// json.NewEncoder(w).Encode(...) and then write a trailing plain-text
// string ("Get rooms" / "Get plugs") straight after it, so the raw
// response body is JSON immediately followed by garbage text. The Go
// encoder always emits compact (single-line) JSON terminated by "\n",
// so we can recover the real payload by only parsing the first line.
async function parseBody(text) {
  if (!text) return null;
  const firstLine = text.split("\n")[0];
  try {
    return JSON.parse(firstLine);
  } catch {
    return null;
  }
}

export async function request(path, options = {}) {
  const res = await fetch(path, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
  });

  const text = await res.text();

  if (!res.ok) {
    throw new Error(text || `Request failed with status ${res.status}`);
  }

  return parseBody(text);
}
