// Assigns a warm tint to a room. Rooms only have a free-text name (no type
// field in the data model), so we match common keywords first and fall back
// to cycling through the palette by ID — keeps every room panel visually
// distinct without needing backend changes.

const KEYWORD_TINTS = [
  { pattern: /cocina/i, tint: "sand" },
  { pattern: /ba(ñ|n)o/i, tint: "slate" },
  { pattern: /(dormitorio|habitaci[oó]n|cuarto|pieza)/i, tint: "blush" },
  { pattern: /(living|sala|estar)/i, tint: "sage" },
  { pattern: /comedor/i, tint: "clay" },
  { pattern: /(garage|garaje|bodega)/i, tint: "slate" },
  { pattern: /(patio|jard[ií]n|terraza)/i, tint: "sage" },
  { pattern: /(oficina|estudio)/i, tint: "linen" },
];

const FALLBACK_TINTS = ["clay", "sage", "sand", "slate", "blush", "linen"];

export function getRoomTint(room) {
  const name = room?.Name ?? "";
  const match = KEYWORD_TINTS.find((entry) => entry.pattern.test(name));
  if (match) return match.tint;
  const seed = room?.ID ?? 0;
  return FALLBACK_TINTS[seed % FALLBACK_TINTS.length];
}
